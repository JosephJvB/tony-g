package googlesheets

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/joho/godotenv"
)

func TestGoogleSheets(t *testing.T) {

	t.Run("append to undefined map key", func(t *testing.T) {
		m := map[int][]string{}

		// m[20] is not set - will append throw?
		m[20] = append(m[20], "123")

		// t.Logf("%v", m)

		if len(m[20]) != 1 {
			t.Error("something went wrong")
		}
	})
	t.Run("can load videos from google sheets", func(t *testing.T) {
		t.Skip("skip test calling real google sheets api")

		err := godotenv.Load("../../.env")
		if err != nil {
			panic(err)
		}

		// .env file does not handle private keys gracefully
		// probably would be better saved to a file than in .env. Oh well.
		invalidKey := os.Getenv("GOOGLE_SHEETS_PRIVATE_KEY")
		fixedKey := strings.ReplaceAll(invalidKey, "__n__", "\n")

		gs := NewClient(Secrets{
			Email:      os.Getenv("GOOGLE_SHEETS_EMAIL"),
			PrivateKey: fixedKey,
		})

		tracks := gs.GetAppleTracks()

		if len(tracks) == 0 {
			t.Errorf("Expected parsed videos to be loaded")
		}

		b, err := json.MarshalIndent(tracks, "", "	")
		if err != nil {
			panic(err)
		}

		err = os.WriteFile("../../data/scraped-tracks.json", b, 0666)
		if err != nil {
			panic(err)
		}
	})

	t.Run("can append tracks to google sheets", func(t *testing.T) {
		t.Skip("skip test calling real google sheets api")

		err := godotenv.Load("../../.env")
		if err != nil {
			panic(err)
		}

		// .env file does not handle private keys gracefully
		// probably would be better saved to a file than in .env. Oh well.
		invalidKey := os.Getenv("GOOGLE_SHEETS_PRIVATE_KEY")
		fixedKey := strings.ReplaceAll(invalidKey, "__n__", "\n")

		gs := NewClient(Secrets{
			Email:      os.Getenv("GOOGLE_SHEETS_EMAIL"),
			PrivateKey: fixedKey,
		})

		toAdd := []AppleTrackRow{
			{
				Title:      "song 9",
				Artist:     "artist 9",
				Album:      "album 9",
				SpotifyUrl: "https://open.spotify.com/track/123",
				Year:       2023,
				AddedAt:    "2024-04-16T00:00:00.000Z",
			},
			{
				Title:      "song 2",
				Artist:     "artist 2",
				Album:      "album 2",
				SpotifyUrl: "https://open.spotify.com/track/123",
				Year:       2025,
				AddedAt:    "2025-04-16T00:00:00.000Z",
			},
		}

		gs.AddAppleTracks(toAdd)
	})

	t.Run("can update 4 source and info columns", func(t *testing.T) {
		t.Skip("skip test calling real google sheets api")

		err := godotenv.Load("../../.env")
		if err != nil {
			panic(err)
		}

		// .env file does not handle private keys gracefully
		// probably would be better saved to a file than in .env. Oh well.
		invalidKey := os.Getenv("GOOGLE_SHEETS_PRIVATE_KEY")
		fixedKey := strings.ReplaceAll(invalidKey, "__n__", "\n")

		gs := NewClient(Secrets{
			Email:      os.Getenv("GOOGLE_SHEETS_EMAIL"),
			PrivateKey: fixedKey,
		})

		values := make([][]interface{}, 4)
		v1 := make([]interface{}, 2)
		v1[0] = "Spotify"
		v1[1] = "Spotify Track information!"
		values[0] = v1

		v2 := make([]interface{}, 2)
		v2[0] = ""
		v2[1] = ""
		values[1] = v2

		v3 := make([]interface{}, 2)
		v3[0] = "GoogleSearch"
		v3[1] = "Found this one from Google"
		values[2] = v3

		v4 := make([]interface{}, 2)
		v4[0] = "GoogleSearch"
		v4[1] = "Found this one from Google"
		values[3] = v4

		gs.updateValues(TestYoutubeTrackSheet, "C2:D", values)
	})

	t.Run("can update 4 source and info columns with TEST method", func(t *testing.T) {
		t.Skip("skip test calling real google sheets api")

		err := godotenv.Load("../../.env")
		if err != nil {
			panic(err)
		}

		// .env file does not handle private keys gracefully
		// probably would be better saved to a file than in .env. Oh well.
		invalidKey := os.Getenv("GOOGLE_SHEETS_PRIVATE_KEY")
		fixedKey := strings.ReplaceAll(invalidKey, "__n__", "\n")

		gs := NewClient(Secrets{
			Email:      os.Getenv("GOOGLE_SHEETS_EMAIL"),
			PrivateKey: fixedKey,
		})

		rows := []YoutubeTrackRow{}

		for i := range 4 {
			rows = append(rows, YoutubeTrackRow{
				Source:         "Spotify-" + strconv.Itoa(i),
				FoundTrackInfo: "Info-" + strconv.Itoa(i),
			})
		}

		gs.UpdateTESTSourceInfo(rows)
	})

	t.Run("can update all source and info columns with TEST method", func(t *testing.T) {
		t.Skip("skip test calling real google sheets api")

		err := godotenv.Load("../../.env")
		if err != nil {
			panic(err)
		}

		// .env file does not handle private keys gracefully
		// probably would be better saved to a file than in .env. Oh well.
		invalidKey := os.Getenv("GOOGLE_SHEETS_PRIVATE_KEY")
		fixedKey := strings.ReplaceAll(invalidKey, "__n__", "\n")

		gs := NewClient(Secrets{
			Email:      os.Getenv("GOOGLE_SHEETS_EMAIL"),
			PrivateKey: fixedKey,
		})

		rows := gs.getRows(TestYoutubeTrackSheet)
		nextRows := []YoutubeTrackRow{}
		for i := range rows {
			r := YoutubeTrackRow{
				Source:         "Spotify-" + strconv.Itoa(i),
				FoundTrackInfo: "Info-" + strconv.Itoa(i),
			}
			nextRows = append(nextRows, r)
		}

		gs.UpdateTESTSourceInfo(nextRows)
	})

	// in case I wanna not update all cells... but maybe I can still update all cells
	t.Run("can update 4 source and info columns: dynamic notation", func(t *testing.T) {
		t.Skip("skip test calling real google sheets api")

		err := godotenv.Load("../../.env")
		if err != nil {
			panic(err)
		}

		// .env file does not handle private keys gracefully
		// probably would be better saved to a file than in .env. Oh well.
		invalidKey := os.Getenv("GOOGLE_SHEETS_PRIVATE_KEY")
		fixedKey := strings.ReplaceAll(invalidKey, "__n__", "\n")

		gs := NewClient(Secrets{
			Email:      os.Getenv("GOOGLE_SHEETS_EMAIL"),
			PrivateKey: fixedKey,
		})

		values := make([][]interface{}, 4)
		v1 := make([]interface{}, 2)
		v1[0] = "Spotify"
		v1[1] = "Spotify Track information!"
		values[0] = v1

		v2 := make([]interface{}, 2)
		v2[0] = ""
		v2[1] = ""
		values[1] = v2

		v3 := make([]interface{}, 2)
		v3[0] = "GoogleSearch"
		v3[1] = "Found this one from Google"
		values[2] = v3

		v4 := make([]interface{}, 2)
		v4[0] = "GoogleSearch"
		v4[1] = "Found this one from Google"
		values[3] = v4

		notation := fmt.Sprintf("C%d:D", 2+3)

		gs.updateValues(TestYoutubeTrackSheet, notation, values)
	})

	t.Run("TEST has almost all rows with source and info", func(t *testing.T) {
		t.Skip("skip test calling real google sheets api. Also the data will be different")

		err := godotenv.Load("../../.env")
		if err != nil {
			panic(err)
		}

		// .env file does not handle private keys gracefully
		// probably would be better saved to a file than in .env. Oh well.
		invalidKey := os.Getenv("GOOGLE_SHEETS_PRIVATE_KEY")
		fixedKey := strings.ReplaceAll(invalidKey, "__n__", "\n")

		gs := NewClient(Secrets{
			Email:      os.Getenv("GOOGLE_SHEETS_EMAIL"),
			PrivateKey: fixedKey,
		})

		tracks := gs.GetTESTTracks()

		if len(tracks) == 0 {
			t.Error("Failed to load TEST tracks")
		}

		if tracks[0].Title != "Young" {
			t.Errorf("Expected first track to be Young, got %s", tracks[0].Title)
		}
		if tracks[0].Artist != "Little Simz" {
			t.Errorf("Expected first track to be Little Simz, got %s", tracks[0].Artist)
		}
		if tracks[0].Source != "Spotify-0" {
			t.Errorf("Expected first track to be Spotify-0, got %s", tracks[0].Source)
		}
		if tracks[0].FoundTrackInfo != "Info-0" {
			t.Errorf("Expected first track to be Info-0, got %s", tracks[0].FoundTrackInfo)
		}

		missingSourceInfo := 0
		for _, t := range tracks {
			if t.Source == "" || t.FoundTrackInfo == "" {
				missingSourceInfo++
			}
		}

		if missingSourceInfo != 1 {
			t.Errorf("Expected only 1 test track with missing source and info, got %d", missingSourceInfo)
		}
	})
}
