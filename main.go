package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"tony-gony/internal/googlesheets"
	"tony-gony/internal/scraping"
	"tony-gony/internal/spotify"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	invalidKey := os.Getenv("GOOGLE_SHEETS_PRIVATE_KEY")
	fixedKey := strings.ReplaceAll(invalidKey, "__n__", "\n")

	os.Setenv("GOOGLE_SHEETS_PRIVATE_KEY", fixedKey)
}

func main() {
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
		// TODO: from parameter store
		Email:      os.Getenv("GOOGLE_SHEETS_EMAIL"),
		PrivateKey: os.Getenv("GOOGLE_SHEETS_PRIVATE_KEY"),
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
		// TODO: from parameter store
		ClientId:     os.Getenv("SPOTIFY_CLIENT_ID"),
		ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
		RefreshToken: os.Getenv("SPOTIFY_REFRESH_TOKEN"),
	})

	nextRows := []googlesheets.ScrapedTrackRow{}
	foundTracks := []spotify.SpotifyTrack{}
	for i, t := range toLookup {
		fmt.Printf("finding track %d/%d\r", i+1, len(toLookup))
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
