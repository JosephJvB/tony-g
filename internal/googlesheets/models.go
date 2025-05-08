package googlesheets

import (
	"strconv"
	"strings"
	"tony-g/internal/stringutil"
)

type AppleTrackRow struct {
	Title   string
	Artist  string
	Album   string
	Year    int
	Found   bool
	AddedAt string
}

func (atr *AppleTrackRow) GetAppleTrackId() string {
	idParts := stringutil.AppleIdParts{
		Title:  atr.Title,
		Artist: atr.Artist,
		Album:  atr.Album,
		Year:   atr.Year,
	}

	return stringutil.MakeAppleId(idParts)
}

type YoutubeVideoRow struct {
	Id          string
	Title       string
	PublishedAt string
	TotalTracks int
	FoundTracks int
	AddedAt     string
}
type YoutubeTrackRow struct {
	Id               string
	Title            string
	Artist           string
	Found            bool
	Link             string
	VideoId          string
	VideoPublishDate string
	AddedAt          string
}

func RowToAppleTrack(row []interface{}) AppleTrackRow {
	yearStr := row[3].(string)
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		year = -1
	}

	return AppleTrackRow{
		Title:   row[0].(string),
		Artist:  row[1].(string),
		Album:   row[2].(string),
		Year:    year,
		Found:   strings.ToUpper(row[4].(string)) == "TRUE",
		AddedAt: row[5].(string),
	}
}

func AppleTrackToRow(track AppleTrackRow) []interface{} {
	r := make([]interface{}, 6)
	r[0] = track.Title
	r[1] = track.Artist
	r[2] = track.Album
	r[3] = track.Year
	r[4] = track.Found
	r[5] = track.AddedAt

	return r
}

func RowToYoutubeVideo(row []interface{}) YoutubeVideoRow {
	ttStr := row[3].(string)
	tt, err := strconv.Atoi(ttStr)
	if err != nil {
		tt = -1
	}
	ftStr := row[4].(string)
	ft, err := strconv.Atoi(ftStr)
	if err != nil {
		ft = -1
	}

	return YoutubeVideoRow{
		Id:          row[0].(string),
		Title:       row[1].(string),
		PublishedAt: row[2].(string),
		TotalTracks: tt,
		FoundTracks: ft,
		AddedAt:     row[5].(string),
	}
}

func YoutubeVideoToRow(video YoutubeVideoRow) []interface{} {
	r := make([]interface{}, 6)
	r[0] = video.Id
	r[1] = video.Title
	r[2] = video.PublishedAt
	r[3] = video.TotalTracks
	r[4] = video.FoundTracks
	r[5] = video.AddedAt

	return r
}

func RowToYoutubeTrack(row []interface{}) YoutubeTrackRow {
	return YoutubeTrackRow{
		Id:               "",
		Title:            row[0].(string),
		Artist:           row[1].(string),
		Found:            strings.ToUpper(row[2].(string)) == "TRUE",
		Link:             row[3].(string),
		VideoId:          row[4].(string),
		VideoPublishDate: row[5].(string),
		AddedAt:          row[6].(string),
	}
}

func YoutubeTrackToRow(track YoutubeTrackRow) []interface{} {
	r := make([]interface{}, 7)
	r[0] = track.Title
	r[1] = track.Artist
	r[2] = track.Found
	r[3] = track.Link
	r[4] = track.VideoId
	r[5] = track.VideoPublishDate
	r[6] = track.AddedAt

	return r
}
