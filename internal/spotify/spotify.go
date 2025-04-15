package spotify

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"tony-gony/internal/scraping"
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

type Secrets struct {
	ClientId     string
	ClientSecret string
	RefreshToken string
}

func NewClient(secrets Secrets) SpotifyClient {
	return SpotifyClient{
		clientId:      secrets.ClientId,
		clientSecret:  secrets.ClientSecret,
		refreshToken:  secrets.RefreshToken,
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

func (s *SpotifyClient) GetPlaylistByName(name string) SpotifyPlaylist {
	if len(s.tonyPlaylists) == 0 {
		if s.accessToken == "" {
			s.LoadAccessToken()
		}

		s.tonyPlaylists = getPlaylists(s.accessToken)
	}

	// again is this return OK
	for _, p := range s.tonyPlaylists {
		if p.Name == name {
			return p
		}
	}

	// feels kinda garbage
	// maybe this instead: https://stackoverflow.com/questions/2050391/how-to-check-if-a-map-contains-a-key-in-go
	return SpotifyPlaylist{}
}

func getPlaylists(token string) []SpotifyPlaylist {

	apiUrl := ApiBaseUrl + "/me/playlists"

	queryPart := url.Values{}
	queryPart.Set("limit", "50")

	apiUrl += "?" + queryPart.Encode()

	playlists := getPaginatedItems[SpotifyPlaylist](apiUrl, token)

	return playlists
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

func (s *SpotifyClient) FindTrack(t scraping.ScrapedTrack) []SpotifyTrack {
	if s.accessToken == "" {
		s.LoadAccessToken()
	}

	trackQuery := "track:" + t.Title
	trackQuery += " artist:" + t.Artist
	trackQuery += " year:" + strconv.Itoa(t.Year)
	if t.Album != "" {
		// apple music adds " - EP" | " - Single" to album suffix sometimes
		// I think that would break the spotify query
		// could trim that, but prefer to remove it
		// trackQuery += " album:" + t.Album
	}

	queryPart := url.Values{}
	queryPart.Set("q", trackQuery)
	queryPart.Set("type", "track")
	queryPart.Set("limit", "1")

	apiUrl := ApiBaseUrl + "/search?" + queryPart.Encode()

	req, _ := http.NewRequest("GET", apiUrl, nil)

	authHeaderValue := "Bearer " + s.accessToken
	req.Header.Set("Authorization", authHeaderValue)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode > 299 {
		log.Fatalf("\nFindTrack failed: \"%s\"\n%s\n", resp.Status, apiUrl)
	}

	trackResponse := SpotifyTrackSearchResults{}
	json.NewDecoder(resp.Body).Decode(&trackResponse)

	return trackResponse.Tracks.Items
}
