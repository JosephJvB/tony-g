package main

import (
	"fmt"
	"os"
	"time"
	"tony-gony/internal/googlesheets"
	"tony-gony/internal/scraping"
)

func main() {
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

	// then lookup those tracks in spotify
	// then get my spotify playlists
	// then create one if not exists
	// then add found tracks to spotify playlist if not already in there
	// then mark tracks as found / not
	// then update scraped tracks google sheet
}
