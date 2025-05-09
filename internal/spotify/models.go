package spotify

import (
	"regexp"
	"strings"
)

const ApplePlaylistPrefix = "Now That's What I Call Melon Music: "
const YoutubePlaylistPrefix = "Internet's nerdiest businessman: "
const JvbSpotifyId = "xnmacgqaaa6a1xi7uy2k1fe7w"

type SpotifyArtist struct {
	Id   string `json:"id"`
	Uri  string `json:"uri"`
	Href string `json:"href"`
	Name string `json:"name"`
}

type SpotifyTrack struct {
	Id      string          `json:"id"`
	Uri     string          `json:"uri"`
	Href    string          `json:"href"`
	Name    string          `json:"name"`
	Artists []SpotifyArtist `json:"artists"`
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

// https://stackoverflow.com/questions/4292468/javascript-regex-remove-text-between-parentheses#answer-4292483
// .replace(/*\([^)]*\)*/g, ”)
// maybe need to handle "ft." ? but this seems enough for now
func CleanSongTitle(songTitle string) string {
	rmParens := regexp.MustCompile(`\\*\(feat.[^)]*\)*`)
	rmSquareBrackets := regexp.MustCompile(`\\*\[feat.[^)]*\]*`)
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
