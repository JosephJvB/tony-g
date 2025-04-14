package youtube

import (
	"strconv"
	"strings"
	"time"
)

var BestTrackPrefixes = []string{"!!!BEST", "!!BEST", "!BEST"}

type BestTrack struct {
	Id                 string
	Name               string
	Artist             string
	Link               string
	Year               int
	VideoPublishedDate string
}

func ParseVideo(video PlaylistItem) []BestTrack {
	tracks := []BestTrack{}

	videoDate, err := time.Parse(time.RFC3339, video.Snippet.PublishedAt)
	publishedAtYear := -1
	if err == nil {
		publishedAtYear = videoDate.Year()
	}

	bestTrackList := getBestTracksList(video.Snippet.Description)
	if len(bestTrackList) == 0 {
		return tracks
	}

	for _, trackStr := range bestTrackList {
		ts := strings.Split(trackStr, "\n")

		if len(ts) < 2 {
			continue
		}

		artistTrackName := strings.Split(ts[0], " - ")
		if len(artistTrackName) < 2 {
			continue
		}

		t := BestTrack{
			Name:               strings.TrimSpace(artistTrackName[1]),
			Artist:             strings.TrimSpace(artistTrackName[0]),
			Link:               strings.TrimSpace(ts[1]),
			Year:               publishedAtYear,
			VideoPublishedDate: video.Snippet.PublishedAt,
		}
		// don't really need id anymore but the column still exists in the spreadsheet
		t.Id = strings.Join([]string{t.Artist, t.Name, strconv.Itoa(t.Year)}, "__")

		tracks = append(tracks, t)
	}

	return tracks
}

func getBestTracksList(description string) []string {
	replacer := strings.NewReplacer(
		"–", "-",
		"\r", "",
		"\n \n", "\n\n",
	)

	d := replacer.Replace(description)

	bestSectionStr := ""
	for _, sectionStr := range strings.Split(d, "\n\n\n") {
		s := strings.ReplaceAll(sectionStr, "!", "")
		if strings.HasPrefix(s, "BEST") {
			bestSectionStr = sectionStr
		}
	}

	if bestSectionStr == "" {
		return []string{}
	}

	tracks := strings.Split(bestSectionStr, "\n\n")

	if len(tracks) <= 1 {
		return []string{}
	}

	return tracks[1:]
}
