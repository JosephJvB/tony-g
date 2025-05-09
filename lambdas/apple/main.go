package scrapeapple

import (
	"fmt"
	"strconv"
	"time"
	"tony-g/internal/apple"
	"tony-g/internal/googlesheets"
	"tony-g/internal/spotify"
	"tony-g/internal/ssm"

	"github.com/aws/aws-lambda-go/lambda"
)

// AppleTrackListItem
// -> ScrapedTrack - can I skip this?
// -> AppleTrackRow

// if this lambda was regularly invoked
// I would initialize AWS clients here
// but it's not so I wont!
// var ()
// func init() {}

type Evt struct {
	Year int `json:"year"`
}

func handleLambdaEvent(evt Evt) {
	now := time.Now()
	timestamp := now.Format(time.RFC3339)
	if evt.Year == 0 {
		evt.Year = now.Year()
	}

	ac := apple.NewClient()
	scrapedTracks := ac.GetTracksForYear(evt.Year)

	fmt.Printf(
		"scraped %d tracks apple music playlist:%d\n",
		len(scrapedTracks),
		evt.Year,
	)

	if len(scrapedTracks) == 0 {
		return
	}

	paramClient := ssm.NewClient()
	paramClient.LoadParameterValues()

	gs := googlesheets.NewClient(googlesheets.Secrets{
		Email:      paramClient.GoogleClientEmail.Value,
		PrivateKey: paramClient.GooglePrivateKey.Value,
	})
	prevTracks := gs.GetAppleTracks()

	fmt.Printf(
		"loaded %d tracks from google sheets\n",
		len(prevTracks),
	)

	m := map[string]bool{}
	for _, t := range prevTracks {
		m[t.GetAppleTrackId()] = true
	}

	// don't lookup tracks if they're already in Google Sheets
	toLookup := []apple.ScrapedTrack{}
	for _, t := range scrapedTracks {
		i := t.GetAppleTrackId()

		if !m[i] {
			toLookup = append(toLookup, t)
		}
	}

	fmt.Printf("you gotta find %d tracks\n", len(toLookup))

	if len(toLookup) == 0 {
		return
	}

	spc := spotify.NewClient(spotify.Secrets{
		ClientId:     paramClient.SpotifyClientId.Value,
		ClientSecret: paramClient.SpotifyClientSecret.Value,
		RefreshToken: paramClient.SpotifyRefreshToken.Value,
	})

	nextRows := []googlesheets.AppleTrackRow{}
	foundTracks := []spotify.SpotifyTrack{}
	for i, t := range toLookup {
		fmt.Printf("finding track %d/%d\r", i+1, len(toLookup))
		results := spc.FindTrack(spotify.FindTrackInput{
			Title:  t.Title,
			Artist: t.Artist,
		})

		if len(results) > 0 {
			foundTracks = append(foundTracks, results[0])
		}

		nextRows = append(nextRows, googlesheets.AppleTrackRow{
			Title:   t.Title,
			Artist:  t.Artist,
			Album:   t.Album,
			Year:    t.Year,
			Found:   len(results) > 0,
			AddedAt: timestamp,
		})
	}

	fmt.Printf("\nfound %d/%d tracks\n", len(foundTracks), len(toLookup))

	tonyPlaylistName := spotify.TonyPlaylistPrefix + strconv.Itoa(evt.Year)
	fmt.Printf("finding playlist %s\n", tonyPlaylistName)

	myPlaylists := spc.GetMyPlaylists()
	fmt.Printf("loaded %d playlists\n", len(myPlaylists))

	// issue: same as previous service, sometimes this code is not finding playlist by name
	// 1. did I not structure the name correctly?
	// 2. more likely: I didn't load the right playlist from Spotify
	// if problem persists, I'll make a new Sheet storing year -> playlistId mapping
	//
	// choosing this as my pattern for handling struct not found in list
	// copying `value, ok := dict["key"] access`
	tonyPlaylist, ok := spotify.SpotifyPlaylist{}, false
	for _, p := range myPlaylists {
		if p.Name == tonyPlaylistName {
			tonyPlaylist = p
			ok = true
		}
	}

	fmt.Printf("spotify playlist for %d exists: %t\n", evt.Year, ok)

	// keyed by spotify track id
	currentTrackMap := map[string]bool{}
	if !ok {
		tonyPlaylist = spc.CreatePlaylist(tonyPlaylistName)
	} else {
		currentTracks := spc.GetPlaylistItems(tonyPlaylist.Id)
		for _, t := range currentTracks {
			currentTrackMap[t.Track.Id] = true
		}
	}

	fmt.Printf("%d tracks already in %d playlist\n", len(currentTrackMap), evt.Year)

	toAdd := []string{}
	for _, t := range foundTracks {
		if !currentTrackMap[t.Id] {
			toAdd = append(toAdd, t.Uri)
		}
	}

	fmt.Printf("adding %d tracks to %d playlist\n", len(toAdd), evt.Year)

	spc.AddPlaylistItems(tonyPlaylist.Id, toAdd)

	fmt.Printf("adding %d rows to scraped google sheet\n", len(nextRows))

	gs.AddAppleTracks(nextRows)
}

func main() {
	lambda.Start(handleLambdaEvent)
}
