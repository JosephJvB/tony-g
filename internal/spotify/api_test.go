package spotify

import (
	"encoding/json"
	"fmt"
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

		found := s.FindTrack(FindTrackInput{
			Title:  "My Golden Years",
			Artist: "The Lemon Twigs",
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

		found := s.FindTrack(FindTrackInput{
			Title:  "Somethin' (feat. Sexyy Red)",
			Artist: "Nardo Wick",
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

		found := s.FindTrack(FindTrackInput{
			// i think i gotta trim "(feat. ...)" and "[feat. ...]"
			Title: "FUNKFEST (feat. TJOnline)", // fails
			// Title:  "FUNKFEST)", // works
			Artist: "grouptherapy., Jadagrace & SWIM",
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

		found := s.FindTrack(FindTrackInput{
			// i think i gotta trim "(feat. ...)" and "[feat. ...]"
			Title: "Times Is Rough (feat. Heem B$F & Rick Hyde)", // fails
			// Title:  "Times Is Rough", // works
			Artist: "Black Soprano Family, Benny the Butcher & DJ Premier",
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

		playlistName := ApplePlaylistPrefix + "2025"
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

	t.Run("can find Blood of the Fang with typo", func(t *testing.T) {
		t.Skip("skip tonys typo")
		err := godotenv.Load("../../.env")
		if err != nil {
			panic(err)
		}

		s := NewClient(Secrets{
			ClientId:     os.Getenv("SPOTIFY_CLIENT_ID"),
			ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
			RefreshToken: os.Getenv("SPOTIFY_REFRESH_TOKEN"),
		})

		found := s.FindTrack(FindTrackInput{
			Title:  "Blood on the Fang", // fails
			Artist: "clipping.",
		})

		if len(found) == 0 {
			t.Error("Failed to find track: 'Blood on the Fang'")
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

	t.Run("can find (What's So Funny Bout) Peace, Love and Understanding", func(t *testing.T) {
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

		found := s.FindTrack(FindTrackInput{
			Title:  "(What's So Funny Bout) Peace, Love and Understanding",
			Artist: "Cheekface",
		})

		if len(found) == 0 {
			t.Error("Failed to find track: (What's So Funny Bout) Peace, Love and Understanding")
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

	t.Run("can find Peace, Love and Understanding", func(t *testing.T) {
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

		found := s.FindTrack(FindTrackInput{
			Title:  "Peace, Love and Understanding",
			Artist: "Cheekface",
		})

		if len(found) == 0 {
			t.Error("Failed to find track: Peace, Love and Understanding")
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

	t.Run("can find Too Fast (Pull Over)", func(t *testing.T) {
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

		found := s.FindTrack(FindTrackInput{
			Title:  "Too Fast (Pull Over)",
			Artist: "Jay Rock, Anderson .Paak",
		})

		if len(found) == 0 {
			t.Error("Failed to find track: Too Fast (Pull Over)")
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

	t.Run("can find Too Fast", func(t *testing.T) {
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

		found := s.FindTrack(FindTrackInput{
			Title:  "Too Fast",
			Artist: "Jay Rock, Anderson .Paak",
		})

		if len(found) == 0 {
			t.Error("Failed to find track: Too Fast")
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

	t.Run("can find Captive of the Sun (Remix)", func(t *testing.T) {
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

		found := s.FindTrack(FindTrackInput{
			Title:  "Captive of the Sun (Remix)",
			Artist: "Parquet Courts",
		})

		if len(found) == 0 {
			t.Error("Failed to find track: Captive of the Sun (Remix)")
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

	t.Run("can find New Level (Remix)", func(t *testing.T) {
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

		found := s.FindTrack(FindTrackInput{
			Title:  "New Level (Remix)",
			Artist: "A$AP Ferg",
		})

		if len(found) == 0 {
			t.Error("Failed to find track: New Level (Remix)")
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

	t.Run("can find Hit Me Where It Hurts (Toro Y Moi Remix)", func(t *testing.T) {
		t.Skip("skip tonys typo")
		err := godotenv.Load("../../.env")
		if err != nil {
			panic(err)
		}

		s := NewClient(Secrets{
			ClientId:     os.Getenv("SPOTIFY_CLIENT_ID"),
			ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
			RefreshToken: os.Getenv("SPOTIFY_REFRESH_TOKEN"),
		})

		found := s.FindTrack(FindTrackInput{
			Title:  "Hit Me Where It Hurts (Toro Y Moi Remix)",
			Artist: "Caroline Polacheck",
		})

		if len(found) == 0 {
			t.Error("Failed to find track: Hit Me Where It Hurts (Toro Y Moi Remix)")
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

	t.Run("can find The Lord of Lightning Vs. Balrog", func(t *testing.T) {
		t.Skip("skip tonys typo")
		err := godotenv.Load("../../.env")
		if err != nil {
			panic(err)
		}

		s := NewClient(Secrets{
			ClientId:     os.Getenv("SPOTIFY_CLIENT_ID"),
			ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
			RefreshToken: os.Getenv("SPOTIFY_REFRESH_TOKEN"),
		})

		found := s.FindTrack(FindTrackInput{
			Title:  "The Lord of Lightning Vs. Balrog",
			Artist: "King Gizzard & The Lizard Wizard",
		})

		if len(found) == 0 {
			t.Error("Failed to find: The Lord of Lightning Vs. Balrog")
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

	t.Run("can find Dress by Buke \u0026 Gase", func(t *testing.T) {
		t.Skip("skip test calling live spotify api.")
		err := godotenv.Load("../../.env")
		if err != nil {
			panic(err)
		}

		s := NewClient(Secrets{
			ClientId:     os.Getenv("SPOTIFY_CLIENT_ID"),
			ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
			RefreshToken: os.Getenv("SPOTIFY_REFRESH_TOKEN"),
		})

		// original txt input
		// Dress (PJ Harvey Cover) by Buke & Gase
		// maybe i need to remove the & char?
		// oh actually the & was encoded as \u0026 ! Maybe that's the issue
		found := s.FindTrack(FindTrackInput{
			Title:  "Dress (PJ Harvey Cover)",
			Artist: "Buke \u0026 Gase",
		})

		if len(found) == 0 {
			t.Error("Failed to find: Dress by Buke & Gase")
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

	t.Run("can find Aquamarine ft. Michael Kiwanuka by Dangermouse & Black Thought", func(t *testing.T) {
		t.Skip("skip test calling live spotify api. Tony!!!!")
		err := godotenv.Load("../../.env")
		if err != nil {
			panic(err)
		}

		s := NewClient(Secrets{
			ClientId:     os.Getenv("SPOTIFY_CLIENT_ID"),
			ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
			RefreshToken: os.Getenv("SPOTIFY_REFRESH_TOKEN"),
		})

		// fails
		// found := s.FindTrack(FindTrackInput{
		// 	Title:  "Aquamarine ft. Michael Kiwanuka",
		// 	Artist: "Dangermouse & Black Thought",
		// })
		// fails
		// found := s.FindTrack(FindTrackInput{
		// 	Title:  "Aquamarine",
		// 	Artist: "Dangermouse & Black Thought",
		// })
		// fails
		// found := s.FindTrack(FindTrackInput{
		// 	Title:  "Aquamarine",
		// 	Artist: "Dangermouse Black Thought",
		// })
		// WORKS
		// found := s.FindTrack(FindTrackInput{
		// 	Title:  "Aquamarine",
		// 	Artist: "Danger Mouse",
		// })
		// WORKS
		// found := s.FindTrack(FindTrackInput{
		// 	Title:  "Aquamarine",
		// 	Artist: "Danger Mouse & Black Thought",
		// })
		// ALSO WORKS
		// so it's a tony issue. Danger Mouse needs a space
		found := s.FindTrack(FindTrackInput{
			Title:  "Aquamarine ft. Michael Kiwanuka",
			Artist: "Danger Mouse & Black Thought",
		})

		if len(found) == 0 {
			t.Error("Failed to find: Aquamarine by Dangermouse & Black Thought")
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

	// no idea why this fails
	t.Run("can find Spirit Possession by Spirit Possession", func(t *testing.T) {
		t.Skip("skip test calling live spotify api. No idea why")
		err := godotenv.Load("../../.env")
		if err != nil {
			panic(err)
		}

		s := NewClient(Secrets{
			ClientId:     os.Getenv("SPOTIFY_CLIENT_ID"),
			ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
			RefreshToken: os.Getenv("SPOTIFY_REFRESH_TOKEN"),
		})

		found := s.FindTrack(FindTrackInput{
			Title:  "Spirit Possession",
			Artist: "Spirit Possession",
		})
		// also this fails
		// found := s.FindTrack(FindTrackInput{
		// 	Title:  "Second Possession",
		// 	Artist: "Spirit Possession",
		// })
		// this works?
		// Something about the word "possession" being in title and in artist causing above queries to fail?
		// found := s.FindTrack(FindTrackInput{
		// 	Title:  "Orthodox Weapons",
		// 	Artist: "Spirit Possession",
		// })

		if len(found) == 0 {
			t.Error("Failed to find: Spirit Possession by Spirit Possession")
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
}
