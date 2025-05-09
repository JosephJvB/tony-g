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

		q := "spotify track clipping blood on the fang"
		result, ok := client.FindSpotifyTrackUrl(q)
		if !ok {
			t.Error("Failed to find spotify track clipping blood on the fang")
		}

		if result != "https://open.spotify.com/track/0jv5VgdENAPV7lHtBlsaXE" {
			t.Errorf("Expected \"https://open.spotify.com/track/0jv5VgdENAPV7lHtBlsaXE\" received \"%s\"", result)
		}
	})
}
