package googlesheets

import (
	"context"
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

var MissingTrackSheet = SheetConfig{
	Name:        "Missing Tracks",
	Id:          1814426117,
	AllRowRange: "A2:F",
}

type MissingTrack struct {
	Id         string
	Name       string
	Artist     string
	Date       string
	Link       string
	SpotifyIds string
}

var ParsedVideosSheet = SheetConfig{
	Name:        "Youtube Videos",
	Id:          404,
	AllRowRange: "A2:D",
}

type ParsedVideo struct {
	Id          string
	Title       string
	PublishedAt string
	TotalTracks string
}

type googlesheets interface {
	LoadSheetData()
}

type GoogleSheetsClient struct {
	sheetsService   *sheets.Service
	parsedVideos    []ParsedVideo
	parsedVideosMap map[string]bool
}

func (gs *GoogleSheetsClient) LoadSheetData() {
	gs.parsedVideos = loadParsedVideos(gs.sheetsService)
	for _, v := range gs.parsedVideos {
		gs.parsedVideosMap[v.Id] = true
	}
}

func loadParsedVideos(service *sheets.Service) []ParsedVideo {
	sheetRange := ParsedVideosSheet.Name + "!" + ParsedVideosSheet.AllRowRange

	resp, err := service.Spreadsheets.Values.Get(SpreadsheetId, sheetRange).Do()
	if err != nil {
		panic(err)
	}

	videos := []ParsedVideo{}
	for _, row := range resp.Values {
		v := ParsedVideo{
			Id:          row[0].(string),
			Title:       row[1].(string),
			PublishedAt: row[2].(string),
			TotalTracks: row[3].(string),
		}

		videos = append(videos, v)
	}

	return videos
}

// https://gist.github.com/karayel/1b915b61d3cf307ca23b14313848f3c4
func NewClient() GoogleSheetsClient {
	conf := &jwt.Config{
		Email:      os.Getenv("GOOGLE_SHEETS_EMAIL"),
		PrivateKey: []byte(os.Getenv("GOOGLE_SHEETS_PRIVATE_KEY")),
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
		sheetsService:   sheetsService,
		parsedVideos:    []ParsedVideo{},
		parsedVideosMap: map[string]bool{},
	}
}
