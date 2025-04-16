package spotify

import (
	"encoding/json"
	"os"
	"testing"

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
}
