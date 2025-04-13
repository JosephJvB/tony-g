package youtube

import (
	"encoding/json"
	"net/http"
	"os"
)

const BaseUrl = "https://www.googleapis.com/youtube/v3"

type IYoutubeClient interface {
	LoadPlaylistItems()
}
type YoutubeClient struct {
	apiKey        string
	playlistId    string
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
		yt.playlistId,
		pageToken,
	)

	yt.playlistItems = append(yt.playlistItems, resp.Items...)

	// recurse
	if resp.NextPageToken != "" {
		yt.LoadPlaylistItems(resp.NextPageToken)
	}
}

func getPlaylistItems(key string, playlistId string, pageToken string) ApiResponse {
	url := BaseUrl + "/playlistItems"
	url += "?maxResults=50"
	url += "&playlistId=" + playlistId
	url += "&part=snippet,status"
	url += "&key=" + key
	if pageToken != "" {
		url += "&pageToken=" + pageToken
	}

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	responseBody := ApiResponse{}
	json.NewDecoder(resp.Body).Decode(&responseBody)

	return responseBody
}

func NewYoutubeClient() YoutubeClient {
	return YoutubeClient{
		apiKey:     os.Getenv("YOUTUBE_API_KEY"),
		playlistId: "PLP4CSgl7K7or84AAhr7zlLNpghEnKWu2c",
	}
}
