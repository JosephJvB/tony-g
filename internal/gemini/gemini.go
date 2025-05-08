package gemini

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/genai"
)

type ParsedTrack struct {
	Title     string `json:"title"`
	Artist    string `json:"artist"`
	SpotifyId string `json:"spotifyId"`
}

type GeminiClient struct {
	client genai.Client
	ctx    context.Context
}

func NewClient(apiKey string) GeminiClient {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatal(err)
	}

	return GeminiClient{
		client: *client,
		ctx:    ctx,
	}
}

func (c *GeminiClient) ParseYoutubeDescription(description string) {
	input := "Return the Best Tracks mentioned in the following text. Each track item should have properties title and artist. Please return the list as valid JSON.\n" + description

	result, err := c.client.Models.GenerateContent(
		c.ctx,
		"gemini-2.0-flash",
		genai.Text(input),
		&genai.GenerateContentConfig{
			Tools: []*genai.Tool{
				{GoogleSearch: &genai.GoogleSearch{}},
			},
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result.Text())
}

func (c *GeminiClient) FindSpotifyUrls(tracks []ParsedTrack) {
	input := "Perform a web search to find following tracks in Songlink and return the Spotify URL from that Songlink page.\n"

	for i, t := range tracks {
		if i != 0 {
			input += ","
		}
		input += " " + t.Title + " by " + t.Artist
	}

	result, err := c.client.Models.GenerateContent(
		c.ctx,
		"gemini-2.0-flash",
		genai.Text(input),
		&genai.GenerateContentConfig{
			Tools: []*genai.Tool{
				{GoogleSearch: &genai.GoogleSearch{}},
			},
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result.Text())
}
