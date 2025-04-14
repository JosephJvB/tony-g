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

		s := NewClient()

		if s.clientId == "" {
			t.Errorf("failed to load clientId from .env")
		}
		if s.clientSecret == "" {
			t.Errorf("failed to load clientSecret from .env")
		}

		s.LoadBasicToken()

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

		s := NewClient()

		if s.refreshToken == "" {
			t.Errorf("failed to load clientId from .env")
		}

		s.LoadAccessToken()

		if s.accessToken == "" {
			t.Errorf("failed to load access token")
		}
	})

	t.Run("Can load tony playlists", func(t *testing.T) {
		t.Skip("skip test calling live spotify api")

		err := godotenv.Load("../../.env")
		if err != nil {
			panic(err)
		}

		s := NewClient()

		s.LoadTonyPlaylists()

		if len(s.tonyPlaylists) == 0 {
			t.Error("Failed to load tony playlists")
		}

		b, err := json.MarshalIndent(s.tonyPlaylists, "", "	")
		if err != nil {
			panic(err)
		}

		err = os.WriteFile("../../playlists.json", b, 0666)
		if err != nil {
			panic(err)
		}
	})

	t.Run("Can load playlist items", func(t *testing.T) {
		t.Skip("skip test calling live spotify api")

		err := godotenv.Load("../../.env")
		if err != nil {
			panic(err)
		}

		s := NewClient()

		items := s.GetPlaylistItems("1Jqly9vntOGxebDIgHSdWt")

		if len(items) == 0 {
			t.Error("Failed to load playlist items")
		}

		b, err := json.MarshalIndent(items, "", "	")
		if err != nil {
			panic(err)
		}

		err = os.WriteFile("../../playlist-items.json", b, 0666)
		if err != nil {
			panic(err)
		}
	})
}
