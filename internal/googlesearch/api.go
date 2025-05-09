package googlesearch

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"

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

func (c *GoogleSearchClient) FindSpotifyTrackUri(t FindTrackInput) (spotifyTrackUrl string, ok bool) {
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

	return linkToTrackUri(resp.Items[0].Link)
}

// input: https://open.spotify.com/track/0jv5VgdENAPV7lHtBlsaXE
// output: spotify:track:0jv5VgdENAPV7lHtBlsaXE
func linkToTrackUri(link string) (string, bool) {
	split := strings.Split(link, "/")
	l := len(split)

	if l < 2 {
		return "", false
	}

	id := split[l-1]
	if len(id) != 22 {
		return "", false
	}

	t := split[l-2]
	if t != "track" {
		return "", false
	}

	uri := "spotify:track:" + id

	return uri, true
}
