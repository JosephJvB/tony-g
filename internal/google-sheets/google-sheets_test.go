package googlesheets

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/joho/godotenv"
)

func TestGoogleSheets(t *testing.T) {
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

		gs := NewClient()

		gs.LoadSheetData()

		if len(gs.parsedVideos) == 0 {
			t.Errorf("Expected parsed videos to be loaded")
		}
		if len(gs.parsedVideosMap) == 0 {
			t.Errorf("Expected parsed videos map to be loaded")
		}

		b, err := json.MarshalIndent(gs.parsedVideos, "", "	")
		if err != nil {
			panic(err)
		}

		err = os.WriteFile("../../videos.json", b, 0666)
		if err != nil {
			panic(err)
		}
	})
}
