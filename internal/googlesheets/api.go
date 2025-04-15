package googlesheets

import (
	"context"
	"strconv"
	"strings"
	util "tony-gony/internal"

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

var ScrapedTracksSheet = SheetConfig{
	Name:        "Scraped Tracks",
	Id:          279196507,
	AllRowRange: "A2:F",
}

type ScrapedTrackRow struct {
	Id      string
	Title   string
	Artist  string
	Album   string
	Year    int
	Found   bool
	AddedAt string
}

type IGoogleSheetsClient interface {
	LoadSheetData()
}

type GoogleSheetsClient struct {
	sheetsService    *sheets.Service
	ScrapedTracks    []ScrapedTrackRow
	ScrapedTracksMap map[string]bool
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

func (gs *GoogleSheetsClient) LoadScrapedTracks() {
	sheetRange := ScrapedTracksSheet.Name + "!" + ScrapedTracksSheet.AllRowRange

	resp, err := gs.sheetsService.Spreadsheets.Values.
		Get(SpreadsheetId, sheetRange).
		Do()
	if err != nil {
		panic(err)
	}

	for _, row := range resp.Values {
		yearStr := row[3].(string)
		year, err := strconv.Atoi(yearStr)
		if err != nil {
			year = -1
		}

		r := ScrapedTrackRow{
			Id:      "",
			Title:   row[0].(string),
			Artist:  row[1].(string),
			Album:   row[2].(string),
			Year:    year,
			Found:   strings.ToUpper(row[4].(string)) == "TRUE",
			AddedAt: row[5].(string),
		}

		r.Id = util.MakeTrackId(util.IdParts{
			Title:  r.Title,
			Artist: r.Artist,
			Album:  r.Album,
			Year:   yearStr,
		})

		gs.ScrapedTracks = append(gs.ScrapedTracks, r)
		gs.ScrapedTracksMap[r.Id] = true
	}
}
