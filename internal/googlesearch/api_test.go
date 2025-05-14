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

	t.Run("search for symbol adrienne lenker", func(t *testing.T) {
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
			Title:  "Symbol",
			Artist: "Adrienne Lenker",
		})
		if !ok {
			t.Fatal("Failed to find spotify track Symbol Adrienne Lenker")
		}
		if result != "https://open.spotify.com/track/5UvgTF3oGUxRwi96UZJd4I" {
			t.Errorf("Expected \"https://open.spotify.com/track/5UvgTF3oGUxRwi96UZJd4I\", received %s\n", result)
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

		cfg := Config{
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
