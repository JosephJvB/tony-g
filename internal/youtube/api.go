package youtube

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const BaseUrl = "https://www.googleapis.com/youtube/v3"
const PlaylistId = "PLP4CSgl7K7or84AAhr7zlLNpghEnKWu2c"

type YoutubeClient struct {
	apiKey string
}

type PlaylistItem struct {
	Snippet struct {
		Title               string `json:"title"`
		Description         string `json:"description"`
		PublishedAt         string `json:"publishedAt"`
		VideoOwnerChannelId string `json:"videoOwnerChannelId"`
		ChannelId           string `json:"channelId"`
		ResourceId          struct {
			Kind    string `json:"kind"`
			VideoId string `json:"videoId"`
		} `json:"resourceId"`
	} `json:"snippet"`
	Status struct {
		PrivacyStatus string `json:"privacyStatus"`
	} `json:"status"`
}

type ApiResponse struct {
	NextPageToken string         `json:"nextPageToken"`
	Items         []PlaylistItem `json:"items"`
}

func (yt *YoutubeClient) LoadPlaylistItems() []PlaylistItem {
	resp := getPlaylistItems(
		yt.apiKey,
		PlaylistId,
		"",
	)

	items := resp.Items
	pageToken := resp.NextPageToken

	for pageToken != "" {
		resp := getPlaylistItems(
			yt.apiKey,
			PlaylistId,
			pageToken,
		)

		items = append(items, resp.Items...)
		pageToken = resp.NextPageToken
	}

	return filterPlaylistItem(items)
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

func NewClient(apiKey string) YoutubeClient {
	return YoutubeClient{
		apiKey: apiKey,
	}
}

func filterPlaylistItem(items []PlaylistItem) []PlaylistItem {
	filtered := []PlaylistItem{}

	for _, item := range items {
		if item.Status.PrivacyStatus == "private" {
			continue
		}
		if item.Snippet.VideoOwnerChannelId != item.Snippet.ChannelId {
			continue
		}
		// album review, ep review, compilation review, mixtape review
		// why are these videos in his playlist? chaotic
		if strings.HasSuffix(strings.TrimSpace(item.Snippet.Title), "REVIEW") {
			continue
		}

		filtered = append(filtered, item)
	}

	return filtered
}
