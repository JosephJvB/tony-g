package spotify

import (
	"encoding/json"
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
	clientId     string
	clientSecret string
	refreshToken string
	basicToken   string
	myPlaylists  []SpotifyPlaylist
}

func NewClient() SpotifyClient {
	return SpotifyClient{
		clientId:     os.Getenv("SPOTIFY_CLIENT_ID"),
		clientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
		refreshToken: os.Getenv("SPOTIFY_REFRESH_TOKEN"),
		myPlaylists:  []SpotifyPlaylist{},
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

	tokenResponse := SpotifyTokenResponse{}
	json.NewDecoder(resp.Body).Decode(&tokenResponse)

	s.basicToken = tokenResponse.AccessToken
}

func (s *SpotifyClient) LoadMyPlaylists() {
	apiUrl := ApiBaseUrl + "/me/playlists"

	queryPart := url.Values{}
	queryPart.Set("limit", "50")

	apiUrl += "?" + queryPart.Encode()

	for apiUrl != "" {
		resp := getPlaylists(apiUrl, s.basicToken)

		s.myPlaylists = append(s.myPlaylists, resp.Items...)

		apiUrl = resp.Next
	}
}

func getPlaylists(apiUrl string, token string) PaginatedResponse[SpotifyPlaylist] {
	req, _ := http.NewRequest("GET", apiUrl, nil)

	bearerToken := "Bearer " + token
	req.Header.Set("Authorization", bearerToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	responseBody := PaginatedResponse[SpotifyPlaylist]{}
	json.NewDecoder(resp.Body).Decode(&responseBody)

	return responseBody
}
