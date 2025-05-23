package googlesheets

import (
	"context"
	"log"
	"os"

	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const SpreadsheetId = "1F5DXCTNZbDy6mFE3Sp1prvU2SfpoqK0dZRsXVHiiOfo"

type SheetConfig struct {
	Name        string
	Id          int
	AllRowRange string
}

var AppleTrackSheet = SheetConfig{
	Name:        "Apple Tracks",
	Id:          1745426463,
	AllRowRange: "A2:F",
}

var YoutubeVideoSheet = SheetConfig{
	Name:        "Youtube Videos",
	Id:          1649352476,
	AllRowRange: "A2:F",
}
var YoutubeTrackSheet = SheetConfig{
	Name:        "Youtube Tracks",
	Id:          1330220669,
	AllRowRange: "A2:I",
}
var TestYoutubeTrackSheet = SheetConfig{
	Name:        "TEST",
	Id:          1316778886,
	AllRowRange: "A2:I",
}

type GoogleSheetsClient struct {
	sheetsService *sheets.Service
}

type Secrets struct {
	Email      string
	PrivateKey string
}

// https://gist.github.com/karayel/1b915b61d3cf307ca23b14313848f3c4
func NewClient(secrets Secrets) GoogleSheetsClient {
	conf := &jwt.Config{
		Email:      secrets.Email,
		PrivateKey: []byte(secrets.PrivateKey),
		TokenURL:   "https://oauth2.googleapis.com/token",
		Scopes: []string{
			"https://www.googleapis.com/auth/spreadsheets",
		},
	}

	client := conf.Client(context.Background())

	sheetsService, err := sheets.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		panic(err)
	}

	return GoogleSheetsClient{
		sheetsService: sheetsService,
	}
}

func (gs *GoogleSheetsClient) GetAppleTracks() []AppleTrackRow {
	rows := gs.getRows(AppleTrackSheet)

	tracks := []AppleTrackRow{}
	for _, row := range rows {
		r := RowToAppleTrack(row)

		tracks = append(tracks, r)
	}

	return tracks
}

func (gs *GoogleSheetsClient) AddAppleTracks(nextRows []AppleTrackRow) {
	// sheets.ValueRange.Values needs interfaces
	rows := make([][]interface{}, len(nextRows))
	for _, t := range nextRows {
		r := AppleTrackToRow(t)

		rows = append(rows, r)
	}

	gs.addRows(AppleTrackSheet, rows)
}

func (gs *GoogleSheetsClient) GetYoutubeVideos() []YoutubeVideoRow {
	rows := gs.getRows(YoutubeVideoSheet)

	videos := []YoutubeVideoRow{}
	for _, row := range rows {
		r := RowToYoutubeVideo(row)

		videos = append(videos, r)
	}

	return videos
}

func (gs *GoogleSheetsClient) AddYoutubeVideos(nextRows []YoutubeVideoRow) {
	// sheets.ValueRange.Values needs interfaces
	rows := make([][]interface{}, len(nextRows))
	for _, t := range nextRows {
		r := YoutubeVideoToRow(t)

		rows = append(rows, r)
	}

	gs.addRows(YoutubeVideoSheet, rows)
}

func (gs *GoogleSheetsClient) GetYoutubeTracks() []YoutubeTrackRow {
	rows := gs.getRows(YoutubeTrackSheet)

	tracks := []YoutubeTrackRow{}
	for _, row := range rows {
		r := RowToYoutubeTrack(row)

		tracks = append(tracks, r)
	}

	return tracks
}

func (gs *GoogleSheetsClient) AddYoutubeTracks(nextRows []YoutubeTrackRow) {
	// sheets.ValueRange.Values needs interfaces
	rows := make([][]interface{}, len(nextRows))
	for _, t := range nextRows {
		r := YoutubeTrackToRow(t)

		rows = append(rows, r)
	}

	gs.addRows(YoutubeTrackSheet, rows)
}

func (gs *GoogleSheetsClient) GetTESTTracks() []YoutubeTrackRow {
	rows := gs.getRows(TestYoutubeTrackSheet)

	tracks := []YoutubeTrackRow{}
	for _, row := range rows {
		r := RowToYoutubeTrack(row)

		tracks = append(tracks, r)
	}

	return tracks
}

func (gs *GoogleSheetsClient) UpdateTESTSourceInfo(nextRows []YoutubeTrackRow) {
	values := make([][]interface{}, len(nextRows))
	for _, t := range nextRows {
		v := make([]interface{}, 2)
		v[0] = t.Source
		v[1] = t.FoundTrackInfo

		values = append(values, v)
	}

	gs.updateValues(TestYoutubeTrackSheet, "C2:D", values)
}

func (gs *GoogleSheetsClient) getRows(cfg SheetConfig) [][]interface{} {
	sheetRange := cfg.Name + "!" + cfg.AllRowRange

	resp, err := gs.sheetsService.Spreadsheets.Values.
		Get(SpreadsheetId, sheetRange).
		Do()
	if err != nil {
		// I think these errors were due to Office Wi-Fi dropping
		// happened twice in a row! And then stopped
		// something like this
		// googleapi: Error 500: Internal error encountered., backendError
		os.WriteFile("./data/googlesheets-error.txt", []byte(err.Error()), 0666)
		log.Fatal(err)
	}

	return resp.Values
}

func (gs *GoogleSheetsClient) addRows(cfg SheetConfig, rows [][]interface{}) {
	// set next rows
	valueRange := sheets.ValueRange{
		MajorDimension: "ROWS",
		Values:         rows,
	}
	// is this range gonna append rows the way I want?
	rowRange := cfg.Name + "!" + cfg.AllRowRange
	req := gs.sheetsService.Spreadsheets.Values.Append(SpreadsheetId, rowRange, &valueRange)
	// is this the only way to add these params?
	req.ValueInputOption("RAW")
	req.InsertDataOption("INSERT_ROWS")

	req.Do()
}

func (gs *GoogleSheetsClient) updateValues(cfg SheetConfig, cellRange string, values [][]interface{}) {
	valueRange := sheets.ValueRange{
		MajorDimension: "ROWS",
		Values:         values,
	}

	updateRange := cfg.Name + "!" + cellRange

	req := gs.sheetsService.Spreadsheets.Values.Update(
		SpreadsheetId,
		updateRange,
		&valueRange,
	)
	req.ValueInputOption("RAW")
	req.Do()
}
