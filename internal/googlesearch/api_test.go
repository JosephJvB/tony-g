package googlesearch

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"google.golang.org/api/customsearch/v1"
)

func TestGoogleSearchApi(t *testing.T) {
	t.Run("search for clipping blood on the fang", func(t *testing.T) {
		t.Skip("Skip calling real CustomSearch API")

		err := godotenv.Load("../../.env")
		if err != nil {
			t.Errorf("Error loading .env file")
		}

		cfg := Secrets{
			ApiKey: os.Getenv("GOOGLE_SEARCH_API_KEY"),
			Cx:     os.Getenv("GOOGLE_SEARCH_CX"),
		}
		client := NewClient(cfg)

		result := client.FindSpotifyTrack(FindTrackInput{
			Title:  "Blood on the Fang",
			Artist: "clipping.",
		})
		if len(result) == 0 {
			t.Error("Failed to find spotify track clipping blood on the fang")
		}

		if result[0].Link != "https://open.spotify.com/track/0jv5VgdENAPV7lHtBlsaXE" {
			t.Errorf("Expected \"https://open.spotify.com/track/0jv5VgdENAPV7lHtBlsaXE\" received \"%s\"", result[0].Link)
		}
	})

	t.Run("search for symbol adrienne lenker", func(t *testing.T) {
		t.Skip("Skip calling real CustomSearch API")

		err := godotenv.Load("../../.env")
		if err != nil {
			t.Errorf("Error loading .env file")
		}

		cfg := Secrets{
			ApiKey: os.Getenv("GOOGLE_SEARCH_API_KEY"),
			Cx:     os.Getenv("GOOGLE_SEARCH_CX"),
		}
		client := NewClient(cfg)

		result := client.FindSpotifyTrack(FindTrackInput{
			Title:  "Symbol",
			Artist: "Adrienne Lenker",
		})
		if len(result) == 0 {
			t.Fatal("Failed to find spotify track Symbol Adrienne Lenker")
		}
		if result[0].Link != "https://open.spotify.com/track/5UvgTF3oGUxRwi96UZJd4I" {
			t.Errorf("Expected \"https://open.spotify.com/track/5UvgTF3oGUxRwi96UZJd4I\", received %s\n", result[0].Link)
		}
	})

	// TODO: test this:
	t.Run("search for Aquamarine ft. Michael Kiwanuka Dangermouse & Black Thought", func(t *testing.T) {
		t.Skip("Skip calling real CustomSearch API")

		err := godotenv.Load("../../.env")
		if err != nil {
			t.Errorf("Error loading .env file")
		}

		cfg := Secrets{
			ApiKey: os.Getenv("GOOGLE_SEARCH_API_KEY"),
			Cx:     os.Getenv("GOOGLE_SEARCH_CX"),
		}
		client := NewClient(cfg)

		result := client.FindSpotifyTrack(FindTrackInput{
			Title:  "Aquamarine ft. Michael Kiwanuka",
			Artist: "Dangermouse & Black Thought",
		})
		if len(result) == 0 {
			t.Fatal("Failed to find spotify track Aquamarine ft. Michael Kiwanuka Dangermouse & Black Thought")
		}
		if result[0].Link != "https://open.spotify.com/track/4so0c8plnJDYDCZtZCfsHr" {
			t.Errorf("Expected \"https://open.spotify.com/track/4so0c8plnJDYDCZtZCfsHr\", received %s\n", result[0].Link)
		}
	})

	// Results: it's a shared quota for the project.
	// I could make another project and cycle that way but let's stop here.
	t.Run("try with cse service", func(t *testing.T) {
		t.Skip("Skip calling real CustomSearch API")

		err := godotenv.Load("../../.env")
		if err != nil {
			t.Errorf("Error loading .env file")
		}

		cfg := Secrets{
			ApiKey: os.Getenv("GOOGLE_SEARCH_API_KEY"),
			Cx:     os.Getenv("GOOGLE_SEARCH_CX"),
		}

		client := NewClient(cfg)

		cse := customsearch.NewCseService(client.svc)

		q := "Blood on the Fang clipping."
		resp, err := cse.List().
			Cx(os.Getenv("GOOGLE_SEARCH_CX_SPECIFIC")).
			Q(q).
			Num(1).
			Do()

		if err != nil {
			t.Fatalf("cse custom search failed %s", err.Error())
		}

		b, err := json.Marshal(resp)
		if err != nil {
			t.Fatalf("failed to marshal response %s", err.Error())
		}

		err = os.WriteFile("../../data/google-customsearch.json", b, 0666)
		if err != nil {
			t.Fatalf("failed to write file %s", err.Error())
		}

		if len(resp.Items) == 0 {
			t.Fatalf("failed to find any search results")
		}

		if resp.Items[0].Link != "https://open.spotify.com/track/0jv5VgdENAPV7lHtBlsaXE" {
			t.Fatalf("failed to find correct track url")
		}
	})
}
