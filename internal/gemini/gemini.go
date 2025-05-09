package gemini

import (
	"context"
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

func (c *GeminiClient) ParseYoutubeDescription(description string) *genai.GenerateContentResponse {
	input := "Return the Best Tracks mentioned in the following text snippet:\n" + description

	result, err := c.client.Models.GenerateContent(
		c.ctx,
		"gemini-2.0-flash",
		genai.Text(input),
		&genai.GenerateContentConfig{
			// Tools: []*genai.Tool{
			// 	{GoogleSearch: &genai.GoogleSearch{}},
			// },
			// can return JSON but not with a google search!
			ResponseMIMEType: "application/json",
			ResponseSchema: &genai.Schema{
				Type: genai.TypeArray,
				Items: &genai.Schema{

					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"title":  {Type: genai.TypeString},
						"artist": {Type: genai.TypeString},
					},
				},
			},
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	return result
}

// Try to fix any typos that come from Youtube Video Description Text Snippet
//
// Deprecated: doesn't work
// using Google Search API to find correct properties
func (c *GeminiClient) ValidateSongProperties(songString string) *genai.GenerateContentResponse {
	input := "Perform a web search for the following song and return the correct song title and artist name in case the input is spelled incorrectly:\nSong:\n" + songString

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

	return result
}

// Try find Spotify URLs for tracks from Youtube Video Description Text Snippet
//
// Deprecated: doesn't work
// using Google Search API to find Spotify URL
func (c *GeminiClient) FindSpotifyUrls(tracks []ParsedTrack) *genai.GenerateContentResponse {
	input := "Perform a google search to find valid Spotify Track URLs for following tracks.\nTracks:\n"

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

	return result
}
