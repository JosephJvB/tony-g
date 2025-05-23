package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"tony-g/internal/googlesearch"
	"tony-g/internal/googlesheets"
	"tony-g/internal/spotify"

	"github.com/joho/godotenv"
)

// add Source and FoundTrackInfo to YoutubeTracks in Google Sheet
// Purpose: at a glance, review what was searched for and what was found
// Background: noticed one case where Google Search returned the wrong track (THE FINALS by Joey Bada$$)
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file:\n%s", err.Error())
	}

	sc := spotify.NewClient(spotify.Secrets{
		ClientId:     os.Getenv("SPOTIFY_CLIENT_ID"),
		ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
		RefreshToken: os.Getenv("SPOTIFY_REFRESH_TOKEN"),
	})
	cs := googlesearch.NewClient(googlesearch.Secrets{
		ApiKey: os.Getenv("GOOGLE_SEARCH_API_KEY"),
		Cx:     os.Getenv("GOOGLE_SEARCH_CX"),
	})

	invalidKey := os.Getenv("GOOGLE_SHEETS_PRIVATE_KEY")
	fixedKey := strings.ReplaceAll(invalidKey, "__n__", "\n")
	gs := googlesheets.NewClient(googlesheets.Secrets{
		Email:      os.Getenv("GOOGLE_SHEETS_EMAIL"),
		PrivateKey: fixedKey,
	})

	testTracks := gs.GetYoutubeTracks()

	googleSearchCount := 0
	for i, track := range testTracks {
		if googleSearchCount >= 100 {
			fmt.Printf("Reached search limit of %d, early exit and update sheet\n", googleSearchCount)
			fmt.Println("Pls resume in 24 hours")
			break
		}

		fmt.Printf("Searching for track %d/%d \r", i+1, len(testTracks))

		if track.Source != "" && track.FoundTrackInfo != "" {
			continue
		}
		// ie: original search never found the track anyway :/
		if track.SpotifyUrl == "" {
			continue
		}

		res := sc.FindTrack(spotify.FindTrackInput{
			Title:  track.Title,
			Artist: track.Artist,
		})
		if len(res) > 0 {
			artistsStr := ""
			for i, a := range res[0].Artists {
				if i > 0 {
					artistsStr += ", "
				}
				artistsStr += a.Name
			}

			testTracks[i].Source = "Spotify Search"
			testTracks[i].FoundTrackInfo = res[0].Name + " By " + artistsStr
			continue
		}

		googleSearchCount++
		fmt.Printf("Google Search #%d \"%s by %s\"\n", googleSearchCount, track.Title, track.Artist)
		res2 := cs.FindSpotifyTrack(googlesearch.FindTrackInput{
			Title:  track.Title,
			Artist: track.Artist,
		})
		if len(res2) > 0 {
			testTracks[i].Source = "Google Search"
			testTracks[i].FoundTrackInfo = res2[0].Title
		} else {
			fmt.Printf("Failed to find track \"%s by %s\"\n", track.Title, track.Artist)
		}

	}

	gs.UpdateYoutubeTracksSourceInfo(testTracks)
}
