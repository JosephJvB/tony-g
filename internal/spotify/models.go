package spotify

import (
	"regexp"
	"strings"
)

const ApplePlaylistPrefix = "Now That's What I Call Melon Music: "
const YoutubePlaylistPrefix = "Melon Music (Deluxe): "
const JvbSpotifyId = "xnmacgqaaa6a1xi7uy2k1fe7w"

type SpotifyArtist struct {
	Id   string `json:"id"`
	Uri  string `json:"uri"`
	Href string `json:"href"`
	Name string `json:"name"`
}

type SpotifyTrack struct {
	Id           string          `json:"id"`
	Uri          string          `json:"uri"`
	Href         string          `json:"href"`
	Name         string          `json:"name"`
	Artists      []SpotifyArtist `json:"artists"`
	ExternalUrls struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
}

type SpotifyPlaylistItem struct {
	AddedAt string       `json:"added_at"`
	Track   SpotifyTrack `json:"track"`
}

type SpotifyPlaylist struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Public        bool   `json:"public"`
	Collaborative bool   `json:"collaborative"`
	Tracks        struct {
		Total int                   `json:"total"`
		Items []SpotifyPlaylistItem `json:"items"`
	} `json:"tracks"`
}

type SpotifyTrackSearchResults struct {
	Tracks struct {
		Href  string         `json:"href"`
		Items []SpotifyTrack `json:"items"`
	} `json:"tracks"`
}

type SpotifyItem interface {
	SpotifyPlaylist | SpotifyPlaylistItem | SpotifyArtist | SpotifyTrack
}

type PaginatedResponse[T SpotifyItem] struct {
	Items []T    `json:"items"`
	Next  string `json:"next"`
}

type SpotifyTokenResponse struct {
	AccessToken string `json:"access_token"`
}

func CleanSongTitle(songTitle string) string {
	// Apple music main cases are (feat. ...) and [feat. ...]
	rmParens := regexp.MustCompile(`\\*\(feat.[^)]*\)*`)
	rmSquareBrackets := regexp.MustCompile(`\\*\[feat.[^)]*\]*`)
	songTitle = rmParens.ReplaceAllLiteralString(songTitle, "")
	songTitle = rmSquareBrackets.ReplaceAllLiteralString(songTitle, "")

	// Youtube description titles
	rmFtDot := regexp.MustCompile(`\\*( ft\..*)`)
	songTitle = rmFtDot.ReplaceAllLiteralString(songTitle, "")
	rmFeatDot := regexp.MustCompile(`\\*( feat\..*)`)
	songTitle = rmFeatDot.ReplaceAllLiteralString(songTitle, "")
	rmProdDot := regexp.MustCompile(`\\*( prod\..*)`)
	songTitle = rmProdDot.ReplaceAllLiteralString(songTitle, "")

	return strings.TrimSpace(songTitle)
}
func RmParens(songTitle string) string {
	rmParens := regexp.MustCompile(`\\*\([^)]*\)*`)
	rmSquareBrackets := regexp.MustCompile(`\\*\[[^)]*\]*`)
	songTitle = rmParens.ReplaceAllLiteralString(songTitle, "")
	songTitle = rmSquareBrackets.ReplaceAllLiteralString(songTitle, "")
	return strings.TrimSpace(songTitle)
}

// input: https://open.spotify.com/track/0jv5VgdENAPV7lHtBlsaXE
// output: spotify:track:0jv5VgdENAPV7lHtBlsaXE
func LinkToTrackUri(link string) (string, bool) {
	split := strings.Split(link, "/")
	l := len(split)

	if l < 2 {
		return "", false
	}

	id := split[l-1]
	if len(id) != 22 {
		return "", false
	}

	t := split[l-2]
	if t != "track" {
		return "", false
	}

	uri := "spotify:track:" + id

	return uri, true
}

func GetTrackInfo(track SpotifyTrack) string {
	artistsStr := ""
	for i, a := range track.Artists {
		if i > 0 {
			artistsStr += ", "
		}
		artistsStr += a.Name
	}

	return track.Name + " By " + artistsStr
}
