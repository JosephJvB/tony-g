package googlesearch

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"

	"google.golang.org/api/customsearch/v1"
	"google.golang.org/api/option"
)

// only 100 free requests a day ay cuz
// https://developers.google.com/custom-search/v1/introduction?authuser=1
// https://github.com/googleapis/google-api-go-client/blob/main/examples/customsearch.go
// oops I'm supposed to move to vertex ai search
// https://cloud.google.com/generative-ai-app-builder/docs/migrate-from-cse
// but then, I'm not using a site restricted search cx atm.
type GoogleSearchClient struct {
	svc *customsearch.Service
	cx  string
}

type Secrets struct {
	ApiKey string
	Cx     string
}

func NewClient(cfg Secrets) GoogleSearchClient {
	ctx := context.Background()
	svc, err := customsearch.NewService(ctx, option.WithAPIKey(cfg.ApiKey))
	if err != nil {
		panic(err)
	}

	return GoogleSearchClient{
		svc: svc,
		cx:  cfg.Cx,
	}
}

type FindTrackInput struct {
	Title  string
	Artist string
}

func (c *GoogleSearchClient) FindSpotifyTrack(t FindTrackInput) []*customsearch.Result {
	if os.Getenv("GOOGLE_SEARCH_DISABLED") != "" {
		fmt.Println("Google Search disabled by .env var \"GOOGLE_SEARCH_DISABLED\"")
		return []*customsearch.Result{}
	}

	q := fmt.Sprintf("%s %s", t.Artist, t.Title)
	q = url.QueryEscape(q)

	resp, err := c.svc.Cse.List().
		Cx(c.cx).
		SiteSearchFilter("i").
		SiteSearch("https://open.spotify.com/track").
		Num(1).
		Q(q).
		Do()

	if err != nil {
		log.Fatalf("google/customsearch/list failed %v\n", err)
	}

	// d, err := resp.MarshalJSON()
	// if err != nil {
	// 	panic(err)
	// }

	// err = os.WriteFile("../../data/google-customsearch.json", d, 0666)
	// if err != nil {
	// 	panic(err)
	// }

	return resp.Items
}
