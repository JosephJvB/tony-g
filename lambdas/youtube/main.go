package scrapeyoutube

import (
	"fmt"
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
	fmt.Printf("Loaded %d youtube videos\n", len(allVideos))

	nextVideos := []youtube.PlaylistItem{}
	for _, v := range allVideos {
		if v.Status.PrivacyStatus == "private" {
			continue
		}
		if v.Snippet.ChannelId != v.Snippet.VideoOwnerChannelId {
			continue
		}

		if !prevVideoMap[v.Id] {
			nextVideos = append(nextVideos, v)
		}
	}
	fmt.Printf("%d Youtube Videos to pull tracks from\n", len(nextVideos))

	gem := gemini.NewClient(
		paramClient.GeminiApiKey.Value,
	)

	nextTrackRows := []googlesheets.YoutubeTrackRow{}
	nextVideoRows := []googlesheets.YoutubeVideoRow{}
	for i, v := range nextVideos {
		fmt.Printf("Getting tracks from description %d/%d\r", i+1, len(nextVideos))
		nextTracks := gem.ParseYoutubeDescription(v.Snippet.Description)

		nv := googlesheets.YoutubeVideoRow{
			Id:          v.Id,
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
				Found:            false,
				Link:             t.Link,
				VideoId:          v.Id,
				VideoPublishDate: v.Snippet.PublishedAt,
				AddedAt:          timestamp,
			}

			nextTrackRows = append(nextTrackRows, r)
		}
	}
	fmt.Printf("Gemini found %d tracks in %d video descriptions\n", len(nextTrackRows), len(nextVideos))

	spc := spotify.NewClient(spotify.Secrets{
		ClientId:     paramClient.SpotifyClientId.Value,
		ClientSecret: paramClient.SpotifyClientSecret.Value,
		RefreshToken: paramClient.SpotifyRefreshToken.Value,
	})
	gcs := googlesearch.NewClient(googlesearch.Config{
		ApiKey: paramClient.GoogleSearchApiKey.Value,
		Cx:     paramClient.GoogleSearchCx.Value,
	})

	toAdd := map[int][]string{}
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
			nextTrackRows[i].Found = true
			toAdd[year] = append(toAdd[year], res[0].Uri)
			foundMap[t.VideoId]++
			continue
		}

		uri, ok := gcs.FindSpotifyTrackUri(googlesearch.FindTrackInput{
			Title:  t.Title,
			Artist: t.Artist,
		})
		if ok {
			nextTrackRows[i].Found = true
			toAdd[year] = append(toAdd[year], uri)
			foundMap[t.VideoId]++
			continue
		}
	}

	// TODO:
	// for year, uris := range toAdd {
	// 	// get or create spotify playlist
	// 	// get spotify playlist items
	// 	// add those uris that aren't already in playlist
	// }

	gs.AddYoutubeTracks(nextTrackRows)
	gs.AddYoutubeVideos(nextVideoRows)
}

func main() {
	lambda.Start(handleLambdaEvent)
}
