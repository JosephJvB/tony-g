package googlesheets

import (
	"strconv"
	"tony-g/internal/stringutil"
)

type AppleTrackRow struct {
	Title      string
	Artist     string
	Album      string
	SpotifyUrl string
	Year       int
	AddedAt    string
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
	Title            string
	Artist           string
	Source           string
	FoundTrackInfo   string
	SpotifyUrl       string
	Link             string
	VideoId          string
	VideoPublishDate string
	AddedAt          string
}

func RowToAppleTrack(row []interface{}) AppleTrackRow {
	yearStr := row[4].(string)
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		year = -1
	}

	return AppleTrackRow{
		Title:      row[0].(string),
		Artist:     row[1].(string),
		Album:      row[2].(string),
		SpotifyUrl: row[3].(string),
		Year:       year,
		AddedAt:    row[5].(string),
	}
}

func AppleTrackToRow(track AppleTrackRow) []interface{} {
	r := make([]interface{}, 6)
	r[0] = track.Title
	r[1] = track.Artist
	r[2] = track.Album
	r[3] = track.SpotifyUrl
	r[4] = strconv.Itoa(track.Year)
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
		Title:            row[0].(string),
		Artist:           row[1].(string),
		Source:           row[2].(string),
		FoundTrackInfo:   row[3].(string),
		SpotifyUrl:       row[4].(string),
		Link:             row[5].(string),
		VideoId:          row[6].(string),
		VideoPublishDate: row[7].(string),
		AddedAt:          row[8].(string),
	}
}

func YoutubeTrackToRow(track YoutubeTrackRow) []interface{} {
	r := make([]interface{}, 9)
	r[0] = track.Title
	r[1] = track.Artist
	r[2] = track.Source
	r[3] = track.FoundTrackInfo
	r[4] = track.SpotifyUrl
	r[5] = track.Link
	r[6] = track.VideoId
	r[7] = track.VideoPublishDate
	r[8] = track.AddedAt

	return r
}
