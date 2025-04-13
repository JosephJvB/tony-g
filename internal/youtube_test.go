package youtube

import (
	"log"
	"testing"

	"github.com/joho/godotenv"
)

func TestYoutube(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	t.Run("it can create a new youtube client", func(t *testing.T) {
		yt := NewYoutubeClient()

		if yt.playlistId == "" {
			t.Errorf("playlistId not set on Youtube Client")
		}

		if yt.apiKey == "" {
			t.Errorf("apiKey not set on Youtube Client")
		}
	})

	t.Run("can load all youtube items", func(t *testing.T) {
		t.Skip("skip test calling YoutubeAPI")

		yt := NewYoutubeClient()

		yt.playlistItems = []PlaylistItem{}
		yt.LoadPlaylistItems("")

		if len(yt.playlistItems) == 0 {
			t.Errorf("Failed to load playlist items")
		}
	})
}
