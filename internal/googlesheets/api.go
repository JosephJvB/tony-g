package googlesheets

import (
	"context"

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
	Id:          279196507,
	AllRowRange: "A2:F",
}

var YoutubeVideosSheet = SheetConfig{
	Name:        "Youtube Videos",
	Id:          1649352476,
	AllRowRange: "A2:F",
}
var YoutubeTrackSheet = SheetConfig{
	Name:        "Youtube Tracks",
	Id:          1330220669,
	AllRowRange: "A2:F",
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
	sheetRange := AppleTrackSheet.Name + "!" + AppleTrackSheet.AllRowRange

	resp, err := gs.sheetsService.Spreadsheets.Values.
		Get(SpreadsheetId, sheetRange).
		Do()
	if err != nil {
		panic(err)
	}

	tracks := []AppleTrackRow{}
	for _, row := range resp.Values {
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

	// set next rows
	valueRange := sheets.ValueRange{
		MajorDimension: "ROWS",
		Values:         rows,
	}
	// is this range gonna append rows the way I want?
	rowRange := AppleTrackSheet.Name + "!" + AppleTrackSheet.AllRowRange
	req := gs.sheetsService.Spreadsheets.Values.Append(SpreadsheetId, rowRange, &valueRange)
	// is this the only way to add these params?
	req.ValueInputOption("RAW")
	req.InsertDataOption("INSERT_ROWS")

	req.Do()
}
