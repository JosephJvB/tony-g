package scrapeyoutube

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"tony-g/internal/gemini"
	"tony-g/internal/googlesearch"
	"tony-g/internal/googlesheets"
	"tony-g/internal/spotify"
	"tony-g/internal/ssm"
	"tony-g/internal/youtube"

	"github.com/aws/aws-lambda-go/lambda"
)

// TODO: only parse videos in input event
type Evt struct {
	VideoIds []string `json:"year"`
}

func handleLambdaEvent(evt Evt) {
	now := time.Now()
	timestamp := now.Format(time.RFC3339)

	paramClient := ssm.NewClient()
	paramClient.LoadParameterValues()

	gs := googlesheets.NewClient(googlesheets.Secrets{
		Email:      paramClient.GoogleClientEmail.Value,
		PrivateKey: paramClient.GooglePrivateKey.Value,
	})

	prevVideos := gs.GetYoutubeVideos()
	fmt.Printf("Loaded %d scraped youtube videos from google sheets\n", len(prevVideos))

	prevVideoMap := map[string]bool{}
	for _, v := range prevVideos {
		prevVideoMap[v.Id] = true
	}

	yt := youtube.NewClient(paramClient.YoutubeApiKey.Value)
	allVideos := yt.LoadPlaylistItems()
	// TODO one at a time
	// allVideos = []youtube.PlaylistItem{
	// 	allVideos[0],
	// }
	fmt.Printf("Loaded %d youtube videos\n", len(allVideos))
	if len(allVideos) == 0 {
		return
	}

	nextVideos := []youtube.PlaylistItem{}
	for _, v := range allVideos {
		if v.Status.PrivacyStatus == "private" {
			continue
		}
		if v.Snippet.ChannelId != v.Snippet.VideoOwnerChannelId {
			continue
		}

		if !prevVideoMap[v.Snippet.ResourceId.VideoId] {
			nextVideos = append(nextVideos, v)
		}
	}
	fmt.Printf("%d Youtube Videos to pull tracks from\n", len(nextVideos))
	if len(nextVideos) == 0 {
		return
	}

	gem := gemini.NewClient(
		paramClient.GeminiApiKey.Value,
	)

	nextTrackRows := []googlesheets.YoutubeTrackRow{}
	nextVideoRows := []googlesheets.YoutubeVideoRow{}
	for i, v := range nextVideos {
		fmt.Printf("Getting tracks from description %d/%d\r", i+1, len(nextVideos))
		nextTracks := gem.ParseYoutubeDescription(v.Snippet.Description)

		nv := googlesheets.YoutubeVideoRow{
			Id:          v.Snippet.ResourceId.VideoId,
			Title:       v.Snippet.Title,
			PublishedAt: v.Snippet.PublishedAt,
			TotalTracks: len(nextTracks),
			FoundTracks: 0,
			AddedAt:     timestamp,
		}
		nextVideoRows = append(nextVideoRows, nv)

		for _, t := range nextTracks {
			r := googlesheets.YoutubeTrackRow{
				Title:            t.Title,
				Artist:           t.Artist,
				SpotifyUrl:       "",
				Link:             t.Url,
				VideoId:          v.Snippet.ResourceId.VideoId,
				VideoPublishDate: v.Snippet.PublishedAt,
				AddedAt:          timestamp,
			}

			nextTrackRows = append(nextTrackRows, r)
		}
	}
	fmt.Printf("Gemini found %d tracks in %d video descriptions\n", len(nextTrackRows), len(nextVideos))
	if len(nextTrackRows) == 0 {
		return
	}

	spc := spotify.NewClient(spotify.Secrets{
		ClientId:     paramClient.SpotifyClientId.Value,
		ClientSecret: paramClient.SpotifyClientSecret.Value,
		RefreshToken: paramClient.SpotifyRefreshToken.Value,
	})
	gcs := googlesearch.NewClient(googlesearch.Config{
		ApiKey: paramClient.GoogleSearchApiKey.Value,
		Cx:     paramClient.GoogleSearchCx.Value,
	})

	toAddByYear := map[int][]string{}
	foundMap := map[string]int{}
	for i, t := range nextTrackRows {
		fmt.Printf("finding track %d/%d\r", i+1, len(nextTrackRows))

		videoDate, err := time.Parse(time.RFC3339, t.VideoPublishDate)
		year := -1
		if err == nil {
			year = videoDate.Year()
		}

		res := spc.FindTrack(spotify.FindTrackInput{
			Title:  t.Title,
			Artist: t.Artist,
		})
		if len(res) > 0 {
			nextTrackRows[i].SpotifyUrl = res[0].Href
			toAddByYear[year] = append(toAddByYear[year], res[0].Uri)
			foundMap[t.VideoId]++
			continue
		}

		href, ok := gcs.FindSpotifyTrackHref(googlesearch.FindTrackInput{
			Title:  t.Title,
			Artist: t.Artist,
		})
		if !ok {
			continue
		}

		uri, ok := spotify.LinkToTrackUri(href)
		if !ok {
			continue
		}

		nextTrackRows[i].SpotifyUrl = href
		toAddByYear[year] = append(toAddByYear[year], uri)
		foundMap[t.VideoId]++
		continue
	}

	myPlaylists := spc.GetMyPlaylists()
	byYear := map[int]spotify.SpotifyPlaylist{}
	for _, p := range myPlaylists {
		if strings.HasPrefix(p.Name, spotify.YoutubePlaylistPrefix) {
			yearStr := strings.TrimPrefix(p.Name, spotify.YoutubePlaylistPrefix)
			year, err := strconv.Atoi(yearStr)
			if err == nil {
				byYear[year] = p
			}
		}
	}

	for year, uris := range toAddByYear {
		playlistName := spotify.YoutubePlaylistPrefix + strconv.Itoa(year)
		fmt.Printf("finding playlist %s\n", playlistName)

		fmt.Printf("loaded %d playlists\n", len(myPlaylists))

		// issue: same as previous service, sometimes this code is not finding playlist by name
		// 1. did I not structure the name correctly?
		// 2. more likely: I didn't load the right playlist from Spotify
		// if problem persists, I'll make a new Sheet storing year -> playlistId mapping
		currentTrackMap := map[string]bool{}
		playlist, ok := byYear[year]
		fmt.Printf("spotify playlist for %d exists: %t\n", year, ok)

		if !ok {
			playlist = spc.CreatePlaylist(playlistName)
		} else {
			currentTracks := spc.GetPlaylistItems(playlist.Id)
			for _, t := range currentTracks {
				currentTrackMap[t.Track.Id] = true
			}
		}

		add := []string{}
		for _, uri := range uris {
			if !currentTrackMap[uri] {
				add = append(add, uri)
			}
		}

		spc.AddPlaylistItems(playlist.Id, add)
	}

	gs.AddYoutubeTracks(nextTrackRows)
	for i, v := range nextVideoRows {
		nextVideoRows[i].FoundTracks = foundMap[v.Id]
	}
	gs.AddYoutubeVideos(nextVideoRows)
}

func main() {
	lambda.Start(handleLambdaEvent)
}
