package spotify

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"tony-gony/internal/scraping"

	"github.com/joho/godotenv"
)

// TODO:
// can i make all these requests with just the basic access token?

func TestSpotify(t *testing.T) {
	t.Run("can get basic token", func(t *testing.T) {
		t.Skip("skip test calling live spotify api")

		err := godotenv.Load("../../.env")
		if err != nil {
			panic(err)
		}

		s := NewClient(Secrets{
			ClientId:     os.Getenv("SPOTIFY_CLIENT_ID"),
			ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
			RefreshToken: os.Getenv("SPOTIFY_REFRESH_TOKEN"),
		})

		if s.clientId == "" {
			t.Errorf("failed to load clientId from .env")
		}
		if s.clientSecret == "" {
			t.Errorf("failed to load clientSecret from .env")
		}

		s.loadBasicToken()

		if s.basicToken == "" {
			t.Errorf("failed to load basic token")
		}
	})

	t.Run("can get access token", func(t *testing.T) {
		t.Skip("skip test calling live spotify api")

		err := godotenv.Load("../../.env")
		if err != nil {
			panic(err)
		}

		s := NewClient(Secrets{
			ClientId:     os.Getenv("SPOTIFY_CLIENT_ID"),
			ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
			RefreshToken: os.Getenv("SPOTIFY_REFRESH_TOKEN"),
		})

		if s.refreshToken == "" {
			t.Errorf("failed to load clientId from .env")
		}

		s.loadAccessToken()

		if s.accessToken == "" {
			t.Errorf("failed to load access token")
		}
	})

	t.Run("can load tony playlists", func(t *testing.T) {
		t.Skip("skip test calling live spotify api")

		err := godotenv.Load("../../.env")
		if err != nil {
			panic(err)
		}

		s := NewClient(Secrets{
			ClientId:     os.Getenv("SPOTIFY_CLIENT_ID"),
			ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
			RefreshToken: os.Getenv("SPOTIFY_REFRESH_TOKEN"),
		})

		p := s.GetMyPlaylists()

		if len(p) == 0 {
			t.Error("Failed to load tony playlists")
		}

		b, err := json.MarshalIndent(p, "", "	")
		if err != nil {
			panic(err)
		}

		err = os.WriteFile("../../data/playlists.json", b, 0666)
		if err != nil {
			panic(err)
		}
	})

	t.Run("can load playlist items", func(t *testing.T) {
		t.Skip("skip test calling live spotify api")

		err := godotenv.Load("../../.env")
		if err != nil {
			panic(err)
		}

		s := NewClient(Secrets{
			ClientId:     os.Getenv("SPOTIFY_CLIENT_ID"),
			ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
			RefreshToken: os.Getenv("SPOTIFY_REFRESH_TOKEN"),
		})

		items := s.GetPlaylistItems("1Jqly9vntOGxebDIgHSdWt")

		if len(items) == 0 {
			t.Error("Failed to load playlist items")
		}

		b, err := json.MarshalIndent(items, "", "	")
		if err != nil {
			panic(err)
		}

		err = os.WriteFile("../../data/playlist-items.json", b, 0666)
		if err != nil {
			panic(err)
		}
	})

	t.Run("can create a playlist", func(t *testing.T) {
		t.Skip("skip test calling live spotify api")

		err := godotenv.Load("../../.env")
		if err != nil {
			panic(err)
		}

		s := NewClient(Secrets{
			ClientId:     os.Getenv("SPOTIFY_CLIENT_ID"),
			ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
			RefreshToken: os.Getenv("SPOTIFY_REFRESH_TOKEN"),
		})

		playlist := s.CreatePlaylist("jvb testo 123")
		if playlist.Name == "" {
			t.Error("Failed to create playlist")
		}
	})

	t.Run("can add items to a playlist", func(t *testing.T) {
		t.Skip("skip test calling live spotify api")

		playlistId := "2bNqL6lXSzg7rZlTxDl117"
		trackUris := []string{
			"spotify:track:0AO3ejChi1gRBWvUDMH2kg",
			"spotify:track:7LXN0LffItjMb9bq61htdB",
		}

		err := godotenv.Load("../../.env")
		if err != nil {
			panic(err)
		}

		s := NewClient(Secrets{
			ClientId:     os.Getenv("SPOTIFY_CLIENT_ID"),
			ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
			RefreshToken: os.Getenv("SPOTIFY_REFRESH_TOKEN"),
		})

		s.AddPlaylistItems(playlistId, trackUris)
	})

	t.Run("can find My Golden Years", func(t *testing.T) {
		t.Skip("skip test calling live spotify api")
		err := godotenv.Load("../../.env")
		if err != nil {
			panic(err)
		}

		s := NewClient(Secrets{
			ClientId:     os.Getenv("SPOTIFY_CLIENT_ID"),
			ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
			RefreshToken: os.Getenv("SPOTIFY_REFRESH_TOKEN"),
		})

		found := s.FindTrack(scraping.ScrapedTrack{
			Title:  "My Golden Years",
			Artist: "The Lemon Twigs",
			Year:   2025,
		})

		if len(found) == 0 {
			t.Error("Failed to find track: My Golden Years")
		}

		b, err := json.MarshalIndent(found, "", "	")
		if err != nil {
			panic(err)
		}

		err = os.WriteFile("../../data/found-track.json", b, 0666)
		if err != nil {
			panic(err)
		}
	})

	t.Run("can find Somethin'", func(t *testing.T) {
		t.Skip("skip test calling live spotify api")
		err := godotenv.Load("../../.env")
		if err != nil {
			panic(err)
		}

		s := NewClient(Secrets{
			ClientId:     os.Getenv("SPOTIFY_CLIENT_ID"),
			ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
			RefreshToken: os.Getenv("SPOTIFY_REFRESH_TOKEN"),
		})

		found := s.FindTrack(scraping.ScrapedTrack{
			Title:  "Somethin' (feat. Sexyy Red)",
			Artist: "Nardo Wick",
			Year:   2025,
		})

		if len(found) == 0 {
			t.Error("Failed to find track: Somethin'")
		}

		b, err := json.MarshalIndent(found, "", "	")
		if err != nil {
			panic(err)
		}

		err = os.WriteFile("../../data/found-track.json", b, 0666)
		if err != nil {
			panic(err)
		}
	})

	t.Run("can find FUNKFEST'", func(t *testing.T) {
		t.Skip("skip test calling live spotify api")
		err := godotenv.Load("../../.env")
		if err != nil {
			panic(err)
		}

		s := NewClient(Secrets{
			ClientId:     os.Getenv("SPOTIFY_CLIENT_ID"),
			ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
			RefreshToken: os.Getenv("SPOTIFY_REFRESH_TOKEN"),
		})

		found := s.FindTrack(scraping.ScrapedTrack{
			// i think i gotta trim "(feat. ...)" and "[feat. ...]"
			Title: "FUNKFEST (feat. TJOnline)", // fails
			// Title:  "FUNKFEST)", // works
			Artist: "grouptherapy., Jadagrace & SWIM",
			Year:   2023,
		})

		if len(found) == 0 {
			t.Error("Failed to find track: FUNKFEST'")
		}

		b, err := json.MarshalIndent(found, "", "	")
		if err != nil {
			panic(err)
		}

		err = os.WriteFile("../../data/found-track.json", b, 0666)
		if err != nil {
			panic(err)
		}
	})

	t.Run("can find Times Is Rough'", func(t *testing.T) {
		t.Skip("skip test calling live spotify api")
		err := godotenv.Load("../../.env")
		if err != nil {
			panic(err)
		}

		s := NewClient(Secrets{
			ClientId:     os.Getenv("SPOTIFY_CLIENT_ID"),
			ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
			RefreshToken: os.Getenv("SPOTIFY_REFRESH_TOKEN"),
		})

		found := s.FindTrack(scraping.ScrapedTrack{
			// i think i gotta trim "(feat. ...)" and "[feat. ...]"
			Title: "Times Is Rough (feat. Heem B$F & Rick Hyde)", // fails
			// Title:  "Times Is Rough", // works
			Artist: "Black Soprano Family, Benny the Butcher & DJ Premier",
			Year:   2022,
		})

		if len(found) == 0 {
			t.Error("Failed to find track: FUNKFEST'")
		}

		b, err := json.MarshalIndent(found, "", "	")
		if err != nil {
			panic(err)
		}

		err = os.WriteFile("../../data/found-track.json", b, 0666)
		if err != nil {
			panic(err)
		}
	})

	t.Run("Can find 2025 playlist", func(t *testing.T) {
		t.Skip("skip test calling live spotify api")

		err := godotenv.Load("../../.env")
		if err != nil {
			panic(err)
		}

		s := NewClient(Secrets{
			ClientId:     os.Getenv("SPOTIFY_CLIENT_ID"),
			ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
			RefreshToken: os.Getenv("SPOTIFY_REFRESH_TOKEN"),
		})

		p := s.GetMyPlaylists()

		if len(p) == 0 {
			t.Error("Failed to load tony playlists")
		}

		fmt.Printf("Found %d playlists\n", len(p))

		playlistName := TonyPlaylistPrefix + "2025"
		playlist, ok := SpotifyPlaylist{}, false
		for _, p := range p {
			if p.Name == playlistName {
				playlist = p
				ok = true
			}
		}

		if !ok {
			t.Error("Failed to find playlist: 2025")
		}

		b, err := json.MarshalIndent(playlist, "", "	")
		if err != nil {
			panic(err)
		}

		err = os.WriteFile("../../data/2025-playlist.json", b, 0666)
		if err != nil {
			panic(err)
		}
	})
}
