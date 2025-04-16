package scraping

import (
	"encoding/json"
	"strconv"
	"strings"
	util "tony-gony/internal"

	"github.com/gocolly/colly/v2"
)

// don't really need a client
// but I guess consistency with other internal API's?
type IScrapingClient interface {
	LoadTracksForYear(year int)
}
type ScrapingClient struct{}

func NewClient() ScrapingClient {
	return ScrapingClient{}
}

func (sc *ScrapingClient) GetTracksForYear(year int) []ScrapedTrack {
	playlistUrl := scrapeApplePlaylistUrlFromTony(year)
	if playlistUrl == "" {
		return []ScrapedTrack{}
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
		trackList[i].Year = year
	}

	return trackList
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

		trackList = getTracksFromServerData(serverData)
	})

	c.Visit(playlistUrl)
	c.Wait()

	return trackList
}

func getTracksFromServerData(serverData []AppleServerData) []ScrapedTrack {
	for _, serverDataItem := range serverData {
		if serverDataItem.Intent.ContentDescriptor.Kind == "playlist" {
			for _, contentSection := range serverDataItem.Data.Sections {
				if contentSection.ItemKind == "trackLockup" {
					return parseAppleTracklists(contentSection.Items)
				}
			}
		}
	}

	// not super happy to return empty struct
	// but in this case it's an empty list so that feels less bad
	return []ScrapedTrack{}
}

func parseAppleTracklists(items []AppleTrackListItem) []ScrapedTrack {
	trackList := []ScrapedTrack{}

	for _, appleTrack := range items {
		t := ScrapedTrack{
			Title:      appleTrack.Title,
			Artist:     appleTrack.ArtistName,
			Album:      getAlbumName(appleTrack),
			DurationMs: appleTrack.Duration,
		}

		trackList = append(trackList, t)
	}

	return trackList
}

func getAlbumName(t AppleTrackListItem) string {
	for _, link := range t.TertiaryLinks {
		if link.Segue.Destination.ContentDescriptor.Kind == "album" {
			return link.Title
		}
	}

	return ""
}
