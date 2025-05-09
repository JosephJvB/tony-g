package googlesearch

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"google.golang.org/api/customsearch/v1"
	"google.golang.org/api/option"
)

// only 100 free requests a day ay cuz
// https://developers.google.com/custom-search/v1/introduction?authuser=1
// https://github.com/googleapis/google-api-go-client/blob/main/examples/customsearch.go
type GoogleSearchClient struct {
	svc *customsearch.Service
	cx  string
}

type Config struct {
	ApiKey string
	Cx     string
}

func NewClient(cfg Config) GoogleSearchClient {
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

func (c *GoogleSearchClient) FindSpotifyTrackHref(t FindTrackInput) (spotifyTrackUrl string, ok bool) {
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
		log.Fatalf("google/customsearch/list failed %v", err)
	}

	// d, err := resp.MarshalJSON()
	// if err != nil {
	// 	panic(err)
	// }

	// err = os.WriteFile("../../data/google-customsearch.json", d, 0666)
	// if err != nil {
	// 	panic(err)
	// }

	if len(resp.Items) == 0 {
		return "", false
	}

	return resp.Items[0].Link, true
}
