package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"time"
	"tony-gony/internal/googlesheets"
	"tony-gony/internal/scraping"
	"tony-gony/internal/spotify"
)

func main() {
	timestamp := time.Now().Format(time.RFC3339)

	sc := scraping.NewClient()

	thisYear := time.Now().Year()
	sc.LoadTracksForYear(thisYear)

	gs := googlesheets.NewClient(googlesheets.Secrets{
		// TODO: from parameter store
		Email:      os.Getenv("GOOGLE_SHEETS_EMAIL"),
		PrivateKey: os.Getenv("GOOGLE_SHEETS_PRIVATE_KEY"),
	})
	gs.LoadScrapedTracks()

	toLookup := []scraping.ScrapedTrack{}
	for _, t := range sc.TracksByYear[thisYear] {
		if !gs.ScrapedTracksMap[t.Id] {
			toLookup = append(toLookup, t)
		}
	}

	fmt.Printf("you gotta find %d tracks\n", len(toLookup))

	spc := spotify.NewClient(spotify.Secrets{
		// TODO: from parameter store
		ClientId:     os.Getenv("SPOTIFY_CLIENT_ID"),
		ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
		RefreshToken: os.Getenv("SPOTIFY_REFRESH_TOKEN"),
	})

	nextRows := []googlesheets.ScrapedTrackRow{}
	foundTracks := []spotify.SpotifyTrack{}
	for _, t := range toLookup {
		results := spc.FindTrack(t)
		if len(results) > 0 {
			foundTracks = append(foundTracks, results[0])
		}

		nextRows = append(nextRows, googlesheets.ScrapedTrackRow{
			Title:   t.Title,
			Artist:  t.Artist,
			Album:   t.Album,
			Year:    t.Year,
			Found:   len(results) > 0,
			AddedAt: timestamp,
		})
	}

	tonyPlaylistName := spotify.TonyPlaylistPrefix + strconv.Itoa(thisYear)
	tonyPlaylist := spc.GetPlaylistByName(tonyPlaylistName)

	currentTrackMap := map[string]bool{}
	if tonyPlaylist.Id == "" {
		// tonyPlaylist = spc.CreatePlaylist(tonyPlaylistName)
	} else {
		currentTracks := spc.GetPlaylistItems(tonyPlaylist.Id)
		for _, t := range currentTracks {
			currentTrackMap[t.Track.Id] = true
		}
	}

	toAdd := []spotify.SpotifyTrack{}
	for _, t := range foundTracks {
		if !currentTrackMap[t.Id] {
			toAdd = append(toAdd, t)
		}
	}

	spc.AddPlaylistItems(tonyPlaylist.Id, toAdd)

	// add next rows to front
	slices.Reverse(nextRows)
	nextRows = append(nextRows, gs.ScrapedTracks...)

	gs.SetScrapedTracks(nextRows)
}
