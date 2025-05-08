package youtube

import (
	"strings"
)

type ScrapedTrack struct {
	Id               string
	Title            string
	Artist           string
	Found            bool
	Link             string
	VideoId          string
	VideoPublishDate string
}

func MakeTrackId(t ScrapedTrack) string {
	id := ""

	partsList := []string{
		t.Title,
		t.Artist,
		t.VideoId,
	}

	for _, p := range partsList {
		p = strings.ToLower(strings.TrimSpace(p))

		if id != "" {
			id += "__"
		}

		id += p
	}

	return id
}
