package youtube

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const BaseUrl = "https://www.googleapis.com/youtube/v3"
const PlaylistId = "PLP4CSgl7K7or84AAhr7zlLNpghEnKWu2c"

type YoutubeClient struct {
	apiKey        string
	PlaylistItems []PlaylistItem
}

type PlaylistItem struct {
	Id      string `json:"id"`
	Snippet struct {
		Description         string `json:"description"`
		PublishedAt         string `json:"publishedAt"`
		VideoOwnerChannelId string `json:"videoOwnerChannelId"`
		ChannelId           string `json:"channelId"`
	} `json:"snippet"`
	Status struct {
		PrivacyStatus string `json:"privacyStatus"`
	} `json:"status"`
}

type ApiResponse struct {
	NextPageToken string         `json:"nextPageToken"`
	Items         []PlaylistItem `json:"items"`
}

func (yt *YoutubeClient) LoadPlaylistItems(pageToken string) {
	resp := getPlaylistItems(
		yt.apiKey,
		PlaylistId,
		pageToken,
	)

	yt.PlaylistItems = append(yt.PlaylistItems, resp.Items...)

	// recurse
	if resp.NextPageToken != "" {
		yt.LoadPlaylistItems(resp.NextPageToken)
	}
}

func getPlaylistItems(key string, playlistId string, pageToken string) ApiResponse {
	apiUrl := BaseUrl + "/playlistItems"

	queryPart := url.Values{}
	queryPart.Set("maxResults", "50")
	queryPart.Set("playlistId", playlistId)
	queryPart.Set("part", "snippet,status")
	queryPart.Set("key", key)
	if pageToken != "" {
		queryPart.Set("pageToken", pageToken)
	}

	apiUrl += "?" + queryPart.Encode()

	resp, err := http.Get(apiUrl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		b := new(strings.Builder)
		io.Copy(b, resp.Body)
		log.Print(b.String())
		log.Fatalf("\ngetPlaylistItems failed: \"%s\"", resp.Status)
	}

	responseBody := ApiResponse{}
	json.NewDecoder(resp.Body).Decode(&responseBody)

	return responseBody
}

func NewClient() YoutubeClient {
	return YoutubeClient{
		apiKey:        os.Getenv("YOUTUBE_API_KEY"),
		PlaylistItems: []PlaylistItem{},
	}
}
