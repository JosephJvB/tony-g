package youtube

import (
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/h2non/gock"
	"github.com/joho/godotenv"
)

func TestYoutube(t *testing.T) {
	t.Run("it can create a new youtube client", func(t *testing.T) {
		t.Skip("its messing with the youtube api key")
		testApiKey := "_test_youtubeApiKey"

		yt := NewClient(testApiKey)

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

		apiKey := os.Getenv("YOUTUBE_API_KEY")
		yt := NewClient(apiKey)

		items := yt.LoadPlaylistItems()

		if len(items) == 0 {
			t.Errorf("Failed to load playlist items")
		}

		b, err := json.MarshalIndent(items, "", "	")
		if err != nil {
			panic(err)
		}

		err = os.WriteFile("../../data/youtube-videos.json", b, 0666)
		if err != nil {
			panic(err)
		}
	})

	// TODO: mock Youtube HTTP response https://pkg.go.dev/net/http/httptest
	t.Run("makes correctly formatted API call", func(t *testing.T) {
		defer gock.Off()
		// gock.Observe(gock.DumpRequest)

		testApiKey := "_test_youtubeApiKey"
		testPageToken := "_test_pageToken"

		gock.New("https://www.googleapis.com").
			Get("/youtube/v3/playlistItems").
			MatchParam("maxResults", "50").
			MatchParam("playlistId", "PLP4CSgl7K7or84AAhr7zlLNpghEnKWu2c").
			MatchParam("part", "snippet,status").
			MatchParam("key", testApiKey).
			Reply(200).
			JSON(map[string]any{
				"nextPageToken": testPageToken,
				"items": []PlaylistItem{
					{
						Id: "_test_id1",
					},
				},
			})
		// 2nd
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
						Id: "_test_id2",
					},
				},
			})

		yt := NewClient(testApiKey)

		items := yt.LoadPlaylistItems()

		if len(items) != 2 {
			t.Errorf("Expected to load two playlist items received %d", len(items))
		}

		if items[0].Id != "_test_id1" {
			t.Errorf("Expected test playlist item 1 to have Id _test_id1. Received %s", items[0].Id)
		}
		if items[1].Id != "_test_id2" {
			t.Errorf("Expected test playlist item 2 to have Id _test_id2. Received %s", items[1].Id)
		}
	})
}
