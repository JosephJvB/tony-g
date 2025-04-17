package main

import (
	"fmt"
	"strconv"
	"time"
	"tony-gony/internal/aws/ssm"
	"tony-gony/internal/googlesheets"
	"tony-gony/internal/scraping"
	"tony-gony/internal/spotify"
	"tony-gony/internal/util"
)

func main() {
	paramClient := ssm.NewClient()
	paramClient.LoadParameterValues()

	timestamp := time.Now().Format(time.RFC3339)
	thisYear := time.Now().Year()

	sc := scraping.NewClient()
	scrapedTracks := sc.GetTracksForYear(thisYear)

	fmt.Printf(
		"scraped %d tracks apple music playlist:%d\n",
		len(scrapedTracks),
		thisYear,
	)

	gs := googlesheets.NewClient(googlesheets.Secrets{
		Email:      paramClient.GoogleClientEmail.Value,
		PrivateKey: paramClient.GooglePrivateKey.Value,
	})
	gs.LoadScrapedTracks()

	fmt.Printf(
		"loaded %d tracks from google sheets\n",
		len(gs.ScrapedTracksMap),
	)

	// don't lookup tracks if they're already in Google Sheets
	toLookup := []scraping.ScrapedTrack{}
	for _, t := range scrapedTracks {
		// keyed by custom id. See `util.go`
		if !gs.ScrapedTracksMap[t.Id] {
			toLookup = append(toLookup, t)
		}
	}

	fmt.Printf("you gotta find %d tracks\n", len(toLookup))

	if len(toLookup) == 0 {
		return
	}

	spc := spotify.NewClient(spotify.Secrets{
		ClientId:     paramClient.SpotifyClientId.Value,
		ClientSecret: paramClient.GoogleClientEmail.Value,
		RefreshToken: paramClient.SpotifyRefreshToken.Value,
	})

	nextRows := []googlesheets.ScrapedTrackRow{}
	foundTracks := []spotify.SpotifyTrack{}
	for i, t := range toLookup {
		fmt.Printf("finding track %d/%d\r", i+1, len(toLookup))
		results := spc.FindTrack(t)

		// on first failure - try normalize track title
		if len(results) == 0 {
			withoutFeatureStr := util.RemoveFeatureString(t.Title)
			if withoutFeatureStr != t.Title {
				t2 := t
				t2.Title = withoutFeatureStr
				results = spc.FindTrack(t2)
			}
		}

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

	fmt.Printf("\nfound %d/%d tracks\n", len(foundTracks), len(toLookup))

	tonyPlaylistName := spotify.TonyPlaylistPrefix + strconv.Itoa(thisYear)
	myPlaylists := spc.GetMyPlaylists()
	// choosing this as my pattern for handling struct not found in list
	// copying `value, ok := dict["key"] access`
	tonyPlaylist, ok := spotify.SpotifyPlaylist{}, false
	for _, p := range myPlaylists {
		if p.Name == tonyPlaylistName {
			tonyPlaylist = p
			ok = true
		}
	}

	fmt.Printf("spotify playlist for %d exists: %t\n", thisYear, ok)

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

	fmt.Printf("%d tracks already in %d playlist\n", len(currentTrackMap), thisYear)

	toAdd := []string{}
	for _, t := range foundTracks {
		if !currentTrackMap[t.Id] {
			toAdd = append(toAdd, t.Uri)
		}
	}

	fmt.Printf("adding %d tracks to %d playlist\n", len(toAdd), thisYear)

	spc.AddPlaylistItems(tonyPlaylist.Id, toAdd)

	fmt.Printf("adding %d rows to scraped google sheet\n", len(nextRows))

	gs.AddNextRows(nextRows)
}
