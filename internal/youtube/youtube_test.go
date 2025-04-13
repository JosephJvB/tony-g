package youtube

import (
	"log"
	"os"
	"testing"

	"github.com/h2non/gock"
	"github.com/joho/godotenv"
)

func TestYoutube(t *testing.T) {
	t.Run("it can create a new youtube client", func(t *testing.T) {
		testApiKey := "_test_youtubeApiKey"
		os.Setenv("YOUTUBE_API_KEY", testApiKey)

		yt := NewClient()

		if yt.playlistId == "" {
			t.Errorf("playlistId not set on Youtube Client")
		}

		if yt.apiKey == "" {
			t.Errorf("apiKey not set on Youtube Client")
		}
	})

	t.Run("can load all youtube items", func(t *testing.T) {
		t.Skip("skip test calling YoutubeAPI")

		// Load actual Youtube API Key
		err := godotenv.Load("../../.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		yt := NewClient()

		yt.playlistItems = []PlaylistItem{}
		yt.LoadPlaylistItems("")

		if len(yt.playlistItems) == 0 {
			t.Errorf("Failed to load playlist items")
		}
	})

	// TODO: mock Youtube HTTP response https://pkg.go.dev/net/http/httptest
	t.Run("makes correctly formatted API call", func(t *testing.T) {
		defer gock.Off()
		// gock.Observe(gock.DumpRequest)

		testApiKey := "_test_youtubeApiKey"
		os.Setenv("YOUTUBE_API_KEY", testApiKey)

		testPageToken := "_test_pageToken"

		gock.New("https://www.googleapis.com").
			Get("/youtube/v3/playlistItems").
			MatchParam("maxResults", "50").
			MatchParam("playlistId", "PLP4CSgl7K7or84AAhr7zlLNpghEnKWu2c").
			MatchParam("part", "snippet,status").
			MatchParam("key", testApiKey).
			MatchParam("pageToken", testPageToken).
			Reply(200).
			JSON(map[string]any{
				"nextPageToken": "",
				"items": []PlaylistItem{
					{
						Id: "_test_id",
					},
				},
			})

		yt := NewClient()

		yt.LoadPlaylistItems(testPageToken)

		if len(yt.playlistItems) != 1 {
			t.Errorf("Expected to load one playlist item received %d", len(yt.playlistItems))
		}

		testItemId := yt.playlistItems[0].Id
		if testItemId != "_test_id" {
			t.Errorf("Expected test playlist item to have Id _test_id. Received %s", testItemId)
		}
	})
}
