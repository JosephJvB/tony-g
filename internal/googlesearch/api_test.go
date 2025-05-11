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

		result, ok := client.FindSpotifyTrackHref(FindTrackInput{
			Title:  "Blood on the Fang",
			Artist: "clipping.",
		})
		if !ok {
			t.Error("Failed to find spotify track clipping blood on the fang")
		}

		if result != "https://open.spotify.com/track/0jv5VgdENAPV7lHtBlsaXE" {
			t.Errorf("Expected \"https://open.spotify.com/track/0jv5VgdENAPV7lHtBlsaXE\" received \"%s\"", result)
		}
	})

	t.Run("loadRequestCount handles file not exist", func(t *testing.T) {
		os.Remove("../../data/locking_test.txt")

		c := loadRequestCount("../../data/locking_test.txt")

		if c != 0 {
			t.Errorf("Expected 0, got %d", c)
		}
	})

	t.Run("loadRequestCount handles invalid file contents", func(t *testing.T) {
		os.WriteFile("../../data/locking_test.txt", []byte("invalid"), 0666)

		c := loadRequestCount("../../data/locking_test.txt")

		if c != 0 {
			t.Errorf("Expected 0, got %d", c)
		}
	})

	t.Run("loadRequestCount loads prev count", func(t *testing.T) {
		os.WriteFile("../../data/locking_test.txt", []byte("299"), 0666)

		c := loadRequestCount("../../data/locking_test.txt")

		if c != 299 {
			t.Errorf("Expected 299, got %d", c)
		}
	})
}
