package spotify

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const ApiBaseUrl = "https://api.spotify.com/v1"
const AccountsBaseUrl = "https://accounts.spotify.com/api"

type ISpotifyClient interface {
	LoadBasicToken()
	LoadMyPlaylists()
}

type SpotifyClient struct {
	clientId      string
	clientSecret  string
	refreshToken  string
	basicToken    string
	accessToken   string
	tonyPlaylists []SpotifyPlaylist
}

func NewClient() SpotifyClient {
	return SpotifyClient{
		clientId:      os.Getenv("SPOTIFY_CLIENT_ID"),
		clientSecret:  os.Getenv("SPOTIFY_CLIENT_SECRET"),
		refreshToken:  os.Getenv("SPOTIFY_REFRESH_TOKEN"),
		tonyPlaylists: []SpotifyPlaylist{},
	}
}

func (s *SpotifyClient) LoadBasicToken() {
	apiUrl := AccountsBaseUrl + "/token"

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", s.clientId)
	data.Set("client_secret", s.clientSecret)

	postData := strings.NewReader(data.Encode())

	resp, err := http.Post(apiUrl, "application/x-www-form-urlencoded", postData)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode > 299 {
		log.Fatalf("\nLoadBasicToken failed: \"%s\"", resp.Status)
	}

	tokenResponse := SpotifyTokenResponse{}
	json.NewDecoder(resp.Body).Decode(&tokenResponse)

	s.basicToken = tokenResponse.AccessToken
}

func (s *SpotifyClient) LoadAccessToken() {
	apiUrl := AccountsBaseUrl + "/token"

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", s.refreshToken)
	postData := strings.NewReader(data.Encode())

	req, _ := http.NewRequest("POST", apiUrl, postData)

	req.SetBasicAuth(s.clientId, s.clientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode > 299 {
		log.Fatalf("\nLoadAccessToken failed: \"%s\"", resp.Status)
	}

	tokenResponse := SpotifyTokenResponse{}
	json.NewDecoder(resp.Body).Decode(&tokenResponse)

	s.accessToken = tokenResponse.AccessToken
}

func (s *SpotifyClient) LoadTonyPlaylists() {
	if s.accessToken == "" {
		s.LoadAccessToken()
	}

	apiUrl := ApiBaseUrl + "/me/playlists"

	queryPart := url.Values{}
	queryPart.Set("limit", "50")

	apiUrl += "?" + queryPart.Encode()

	playlists := getPaginatedItems[SpotifyPlaylist](apiUrl, s.accessToken)

	for _, playlist := range playlists {
		isTony := strings.HasPrefix(playlist.Name, TonyPlaylistPrefix)

		if isTony {
			s.tonyPlaylists = append(s.tonyPlaylists, playlist)
		}
	}
}

func getPaginatedItems[T any](startUrl string, token string) []T {
	apiUrl := startUrl

	items := []T{}

	for apiUrl != "" {
		req, _ := http.NewRequest("GET", apiUrl, nil)

		authHeaderValue := "Bearer " + token
		req.Header.Set("Authorization", authHeaderValue)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}

		if resp.StatusCode > 299 {
			log.Fatalf("\ngetPaginatedItems failed: \"%s\"\n%s\n", resp.Status, apiUrl)
		}

		responseBody := PaginatedResponse[T]{}
		json.NewDecoder(resp.Body).Decode(&responseBody)

		items = append(items, responseBody.Items...)
		apiUrl = responseBody.Next
	}

	return items
}

func (s *SpotifyClient) GetPlaylistItems(playlistId string) []SpotifyPlaylistItem {
	if s.accessToken == "" {
		s.LoadAccessToken()
	}

	apiUrl := ApiBaseUrl + "/playlists/" + playlistId + "/tracks"

	items := getPaginatedItems[SpotifyPlaylistItem](apiUrl, s.accessToken)

	return items
}
