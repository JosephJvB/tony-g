package youtube

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
)

const BaseUrl = "https://www.googleapis.com/youtube/v3"
const PlaylistId = "PLP4CSgl7K7or84AAhr7zlLNpghEnKWu2c"

type IYoutubeClient interface {
	LoadPlaylistItems()
}
type YoutubeClient struct {
	apiKey        string
	playlistItems []PlaylistItem
}

type PlaylistItem struct {
	Id      string `json:"id"`
	Snippet struct {
		Description         string `json:"description"`
		VideoOwnerChannelId string `json:"videoOwnerChannelId"`
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

	yt.playlistItems = append(yt.playlistItems, resp.Items...)

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
		queryPart.Set(pageToken, pageToken)
	}

	apiUrl += "?" + queryPart.Encode()

	resp, err := http.Get(apiUrl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// TODO: handle status codes?
	// if resp.StatusCode > 299 {
	// }

	responseBody := ApiResponse{}
	json.NewDecoder(resp.Body).Decode(&responseBody)

	return responseBody
}

func NewClient() YoutubeClient {
	return YoutubeClient{
		apiKey:        os.Getenv("YOUTUBE_API_KEY"),
		playlistItems: []PlaylistItem{},
	}
}
