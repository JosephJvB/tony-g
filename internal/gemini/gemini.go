package gemini

import (
	"context"
	"encoding/json"
	"log"

	"google.golang.org/genai"
)

type ParsedTrack struct {
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Url    string `json:"url"`
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

func (c *GeminiClient) ParseYoutubeDescription(description string) []ParsedTrack {
	input := "Return the Best Tracks mentioned in the following text snippet"
	input += "\nformat \"{artist} - {title}\n{url}\""
	// handle multi track for one artist case
	input += "\nIf title has one or more slash character and there is more than one url, return multiple tracks and split the titles by slash character"
	input += "\n"
	input += description

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
						"url":    {Type: genai.TypeString},
					},
				},
			},
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	parsedTracks := []ParsedTrack{}
	err = json.Unmarshal([]byte(result.Text()), &parsedTracks)
	if err != nil {
		log.Fatalf("ParseYoutubeDescription: Failed to parse response JSON")
	}

	return parsedTracks
}
