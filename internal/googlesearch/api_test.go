package googlesearch

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestGoogleSearchApi(t *testing.T) {
	t.Run("search for clipping blood on the fang", func(t *testing.T) {
		t.Skip("Skip calling real CustomSearch API")

		err := godotenv.Load("../../.env")
		if err != nil {
			t.Errorf("Error loading .env file")
		}

		cfg := Config{
			ApiKey: os.Getenv("GOOGLE_SEARCH_API_KEY"),
			Cx:     os.Getenv("GOOGLE_SEARCH_CX"),
		}
		client := NewClient(cfg)

		result, ok := client.FindSpotifyTrackUri(FindTrackInput{
			Title:  "Blood on the Fang",
			Artist: "clipping.",
		})
		if !ok {
			t.Error("Failed to find spotify track clipping blood on the fang")
		}

		if result != "spotify:track:0jv5VgdENAPV7lHtBlsaXE" {
			t.Errorf("Expected \"spotify:track:0jv5VgdENAPV7lHtBlsaXE\" received \"%s\"", result)
		}
	})

	t.Run("turns spotify link to uri", func(t *testing.T) {
		input := "https://open.spotify.com/track/0jv5VgdENAPV7lHtBlsaXE"

		result, ok := linkToTrackUri(input)

		if !ok {
			t.Error("linkToTrackUri failed")
		}

		if result != "spotify:track:0jv5VgdENAPV7lHtBlsaXE" {
			t.Errorf("Failed to get uri from link. Received %s", result)
		}
	})
}
