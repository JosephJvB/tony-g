package googlesearch

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"

	"google.golang.org/api/customsearch/v1"
	"google.golang.org/api/option"
)

// TODO REMOVE
const LockingFilePath = "./data/locking.txt"

// only 100 free requests a day ay cuz
// https://developers.google.com/custom-search/v1/introduction?authuser=1
// https://github.com/googleapis/google-api-go-client/blob/main/examples/customsearch.go
type GoogleSearchClient struct {
	svc          *customsearch.Service
	cx           string
	RequestCount int
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

	reqCount := loadRequestCount(LockingFilePath)

	return GoogleSearchClient{
		svc:          svc,
		cx:           cfg.Cx,
		RequestCount: reqCount,
	}
}

type FindTrackInput struct {
	Title  string
	Artist string
}

func (c *GoogleSearchClient) FindSpotifyTrackHref(t FindTrackInput) (spotifyTrackUrl string, ok bool) {
	if os.Getenv("GOOGLE_SEARCH_DISABLED") != "" {
		fmt.Println("Google Search disabled by .env var \"GOOGLE_SEARCH_DISABLED\"")
		return "\n", false
	}

	q := fmt.Sprintf("%s %s\n", t.Artist, t.Title)
	q = url.QueryEscape(q)

	resp, err := c.svc.Cse.List().
		Cx(c.cx).
		SiteSearchFilter("i").
		SiteSearch("https://open.spotify.com/track").
		Num(1).
		Q(q).
		Do()

		// TODO REMOVE
	c.IncrementRequestCount()

	if err != nil {
		log.Fatalf("google/customsearch/list failed %v\n", err)
	}

	// d, err := resp.MarshalJSON()
	// if err != nil {
	// 	panic(err)
	// }

	// err = os.WriteFile("../../data/google-customsearch.json\n", d, 0666))
	// if err != nil {
	// 	panic(err)
	// }

	if len(resp.Items) == 0 {
		return "", false
	}

	return resp.Items[0].Link, true
}

// TODO REMOVE
func (c *GoogleSearchClient) IncrementRequestCount() {
	c.RequestCount++

	cs := strconv.Itoa(c.RequestCount)

	err := os.WriteFile(LockingFilePath, []byte(cs), 0666)
	if err != nil {
		log.Fatalf("Failed to write to locking file %s: %v\n", LockingFilePath, err)
	}
}

// TODO REMOVE
// / TEMPORARY LOCAL CODE
func loadRequestCount(filePath string) int {
	if _, err := os.Stat(filePath); err != nil {
		err := os.WriteFile(filePath, []byte("0"), 0666)
		if err != nil {
			log.Fatalf("Failed to write to file %s: %v\n", filePath, err)
		}
		return 0
	}

	txt, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Failed to read locking file: %s\n%v\n", filePath, err)
		err := os.WriteFile(filePath, []byte("0"), 0666)
		if err != nil {
			log.Fatalf("Failed to write to file %s: %v\n", filePath, err)
		}
		return 0
	}

	v, err := strconv.Atoi(string(txt))
	if err != nil {
		fmt.Printf("Failed to convert locking value to int: %v\n", err)
		err := os.WriteFile(filePath, []byte("0"), 0666)
		if err != nil {
			log.Fatalf("Failed to write to file %s: %v\n", filePath, err)
		}
		return 0
	}

	return v
}
