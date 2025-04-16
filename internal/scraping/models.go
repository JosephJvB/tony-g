package scraping

type ScrapedTrack struct {
	Id         string
	Title      string
	Artist     string
	Album      string
	DurationMs int
	Year       int
}

type AppleTrackListItem struct {
	Title         string `json:"title"`
	ArtistName    string `json:"artistName"`
	Duration      int    `json:"duration"`
	TertiaryLinks []struct {
		Title string `json:"title"`
		Segue struct {
			Destination struct {
				ContentDescriptor struct {
					Kind string `json:"kind"`
				} `json:"contentDescriptor"`
			} `json:"destination"`
		} `json:"segue"`
	} `json:"tertiaryLinks"`
}
type AppleServerData struct {
	Intent struct {
		ContentDescriptor struct {
			Kind string `json:"kind"`
		} `json:"contentDescriptor"`
	} `json:"intent"`
	Data struct {
		Sections []struct {
			Id       string               `json:"id"`
			ItemKind string               `json:"itemKind"`
			Items    []AppleTrackListItem `json:"items"`
		} `json:"sections"`
	} `json:"data"`
}
