package spotify

const TonyPlaylistPrefix = "Now That's What I Call Melon Music: "
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

type PaginatedResponse[T any] struct {
	Items []T    `json:"items"`
	Next  string `json:"next"`
}

type SpotifyTokenResponse struct {
	AccessToken string `json:"access_token"`
}
