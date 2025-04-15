package scraping

import (
	"encoding/json"
	"strconv"
	"strings"
	util "tony-gony/internal"

	"github.com/gocolly/colly/v2"
)

type ScrapedTrack struct {
	Id         string
	Title      string
	Artist     string
	Album      string
	DurationMs int
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

type IScrapingClient interface {
	LoadTracksForYear(year int)
}

type ScrapingClient struct {
	TracksByYear map[int][]ScrapedTrack
}

func NewClient() ScrapingClient {
	return ScrapingClient{
		TracksByYear: map[int][]ScrapedTrack{},
	}
}

func (sc *ScrapingClient) LoadTracksForYear(year int) {
	sc.TracksByYear[year] = []ScrapedTrack{}

	playlistUrl := scrapeApplePlaylistUrlFromTony(year)
	if playlistUrl == "" {
		return
	}

	trackList := scrapeTrackListFromApple(playlistUrl)

	yearStr := strconv.Itoa(year)

	for i, t := range trackList {
		trackList[i].Id = util.MakeTrackId(util.IdParts{
			Title:  t.Title,
			Artist: t.Artist,
			Album:  t.Album,
			Year:   yearStr,
		})
	}

	sc.TracksByYear[year] = trackList
}

// someone recommended this tutorial
// haven't looked at it yet but good to know
// https://www.google.com/search?client=firefox-b-d&q=Akhil+Sharma+golang+scraper

func scrapeApplePlaylistUrlFromTony(year int) string {
	tonysUrl := "https://theneedledrop.com/loved-list/" + strconv.Itoa(year)

	playlistUrl := ""

	c := colly.NewCollector()

	c.OnHTML("iframe[src^=\"https://embed.music.apple\"]", func(e *colly.HTMLElement) {
		playlistUrl = e.Attr("src")
	})

	c.Visit(tonysUrl)
	c.Wait()

	return strings.Replace(playlistUrl, "embed.music.apple", "music.apple", 1)
}

func scrapeTrackListFromApple(playlistUrl string) []ScrapedTrack {
	trackList := []ScrapedTrack{}

	c := colly.NewCollector()

	// c.OnScraped(func(r *colly.Response) {
	// 	r.Save("../../data/test.html")
	// })

	c.OnHTML("#serialized-server-data", func(e *colly.HTMLElement) {
		serverData := []AppleServerData{}
		json.Unmarshal([]byte(e.Text), &serverData)

		trackList = getTracklistFromServerData(serverData)
	})

	c.Visit(playlistUrl)
	c.Wait()

	return trackList
}

// I could(should) test it separate if I wanted
func getTracklistFromServerData(serverData []AppleServerData) []ScrapedTrack {
	appleTrackList := []AppleTrackListItem{}
	for _, serverDataItem := range serverData {
		if serverDataItem.Intent.ContentDescriptor.Kind == "playlist" {
			for _, contentSection := range serverDataItem.Data.Sections {
				if contentSection.ItemKind == "trackLockup" {
					appleTrackList = contentSection.Items
				}
			}
		}
	}

	trackList := []ScrapedTrack{}
	for _, appleTrack := range appleTrackList {
		album := ""
		for _, link := range appleTrack.TertiaryLinks {
			if link.Segue.Destination.ContentDescriptor.Kind == "album" {
				album = link.Title
			}
		}

		t := ScrapedTrack{
			Title:      appleTrack.Title,
			Artist:     appleTrack.ArtistName,
			Album:      album,
			DurationMs: appleTrack.Duration,
		}

		trackList = append(trackList, t)
	}

	return trackList
}
