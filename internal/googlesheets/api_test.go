package googlesheets

import (
	"encoding/json"
	"os"
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
	t.Run("Can load videos from google sheets", func(t *testing.T) {
		t.Skip("skip test calling real google sheets api")

		err := godotenv.Load("../../.env")
		if err != nil {
			panic(err)
		}

		// .env file does not handle private keys gracefully
		// probably would be better saved to a file than in .env. Oh well.
		invalidKey := os.Getenv("GOOGLE_SHEETS_PRIVATE_KEY")
		fixedKey := strings.ReplaceAll(invalidKey, "__n__", "\n")

		os.Setenv("GOOGLE_SHEETS_PRIVATE_KEY", fixedKey)

		gs := NewClient(Secrets{
			Email:      os.Getenv("GOOGLE_SHEETS_EMAIL"),
			PrivateKey: fixedKey,
		})

		gs.LoadScrapedTracks()

		if len(gs.ScrapedTracks) == 0 {
			t.Errorf("Expected parsed videos to be loaded")
		}
		if len(gs.ScrapedTracksMap) == 0 {
			t.Errorf("Expected parsed videos map to be loaded")
		}

		b, err := json.MarshalIndent(gs.ScrapedTracks, "", "	")
		if err != nil {
			panic(err)
		}

		err = os.WriteFile("../../data/scraped-tracks.json", b, 0666)
		if err != nil {
			panic(err)
		}
	})
}
