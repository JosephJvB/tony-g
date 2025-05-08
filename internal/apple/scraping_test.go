package apple

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestAppleScraping(t *testing.T) {
	t.Run("Can get embedded apple music playlist url for 2025", func(t *testing.T) {
		t.Skip("making live api call")
		playlistUrl := scrapeApplePlaylistUrlFromTony(2025)

		fmt.Printf("url: \"%s\"\n", playlistUrl)

		if playlistUrl == "" {
			t.Errorf("Failed to get playlist url for 2025")
		}
	})

	t.Run("Can get tracklist from apple music playlist url", func(t *testing.T) {
		t.Skip("making live api call")
		playlistUrl := "https://music.apple.com/us/playlist/my-fav-singles-of-2025/pl.u-ayeZTygbKDy"

		trackList := scrapeTrackListFromApple(playlistUrl)

		b, err := json.MarshalIndent(trackList, "", "	")
		if err != nil {
			panic(err)
		}

		err = os.WriteFile("../../data/scraped-tracks.json", b, 0666)
		if err != nil {
			panic(err)
		}

		if len(trackList) == 0 {
			t.Errorf("Failed to get tracklist from url")
		}
	})
}
