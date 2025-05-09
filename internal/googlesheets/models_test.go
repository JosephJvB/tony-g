package googlesheets

import (
	"testing"
)

func TestGoogleSheetsModels(t *testing.T) {
	t.Run("apple track row makes expected id", func(t *testing.T) {
		expectedId := "sea of trees__king gizz__12 bb__2012"

		scrapedTrack := AppleTrackRow{
			Title:  "Sea of Trees",
			Artist: "King Gizz",
			Album:  "12 BB",
			Year:   2012,
		}

		a := scrapedTrack.GetAppleTrackId()
		if a != expectedId {
			t.Errorf("expected %s to equal %s", a, expectedId)
		}
	})

	t.Run("Row to Apple Track", func(t *testing.T) {
		r := make([]interface{}, 6)
		r[0] = "Metal"
		r[1] = "D beffs"
		r[2] = "new album"
		r[3] = "2025"
		r[4] = "TRUE"
		r[5] = "soon"

		at := RowToAppleTrack(r)

		if at.Title != "Metal" {
			t.Errorf("expected title to be Metal. Got %s", at.Title)
		}
		if at.Artist != "D beffs" {
			t.Errorf("expected artist to be D beffs. Got %s", at.Artist)
		}
		if at.Album != "new album" {
			t.Errorf("expected album to be new album. Got %s", at.Album)
		}
		if at.Year != 2025 {
			t.Errorf("expected year to be 2025. Got %d", at.Year)
		}
		if at.Found != true {
			t.Errorf("expected found to be true. Got %t", at.Found)
		}
		if at.AddedAt != "soon" {
			t.Errorf("expected adAddedAt to be soon. Got %s", at.AddedAt)
		}
	})

	t.Run("Apple Track To Row", func(t *testing.T) {
		at := AppleTrackRow{
			Title:   "little things",
			Artist:  "adrianne",
			Album:   "live at revvy hall",
			Year:    2025,
			Found:   true,
			AddedAt: "now",
		}

		row := AppleTrackToRow(at)

		if len(row) != 6 {
			t.Errorf("expected row to have 6 elements. Got %d", len(row))
		}

		if row[0] != "little things" {
			t.Errorf("expected row[0] to be little things. Got %s", row[0])
		}
		if row[1] != "adrianne" {
			t.Errorf("expected row[1] to be adrianne. Got %s", row[1])
		}
		if row[2] != "live at revvy hall" {
			t.Errorf("expected row[2] to be live at revvy hall. Got %s", row[2])
		}
		if row[3] != "2025" {
			t.Errorf("expected row[3] to be 2025. Got %s", row[3])
		}
		if row[4] != "TRUE" {
			t.Errorf("expected row[4] to be TRUE. Got %s", row[4])
		}
		if row[5] != "now" {
			t.Errorf("expected row[5] to be now. Got %s", row[5])
		}
	})

	t.Run("Row to Youtube Video", func(t *testing.T) {
		r := make([]interface{}, 6)
		r[0] = "id-123"
		r[1] = "weekly tracko roundo"
		r[2] = "2025-01-01"
		r[3] = "10"
		r[4] = "5"
		r[5] = "ages ago"

		v := RowToYoutubeVideo(r)

		if v.Id != "id-123" {
			t.Errorf("expected id to be id-123. Got %s", v.Id)
		}
		if v.Title != "weekly tracko roundo" {
			t.Errorf("expected title to be weekly tracko roundo. Got %s", v.Title)
		}
		if v.PublishedAt != "2025-01-01" {
			t.Errorf("expected publishedAt to be 2025-01-01. Got %s", v.PublishedAt)
		}
		if v.TotalTracks != 10 {
			t.Errorf("expected totalTracks to be 10. Got %d", v.TotalTracks)
		}
		if v.FoundTracks != 5 {
			t.Errorf("expected foundTracks to be 5. Got %d", v.FoundTracks)
		}
		if v.AddedAt != "ages ago" {
			t.Errorf("expected addedAt to be ages ago. Got %s", v.AddedAt)
		}
	})

	t.Run("Youtube Video To Row", func(t *testing.T) {
		v := YoutubeVideoRow{
			Id:          "id-123",
			Title:       "weekly tracko roundo",
			PublishedAt: "2025-01-01",
			TotalTracks: 10,
			FoundTracks: 5,
			AddedAt:     "ages ago",
		}

		row := YoutubeVideoToRow(v)

		if row[0] != "id-123" {
			t.Errorf("expected row[0] to be id-123. Got %s", row[0])
		}
		if row[1] != "weekly tracko roundo" {
			t.Errorf("expected row[1] to be weekly tracko roundo. Got %s", row[1])
		}
		if row[2] != "2025-01-01" {
			t.Errorf("expected row[2] to be 2025-01-01. Got %s", row[2])
		}
		if row[3] != 10 {
			t.Errorf("expected row[3] to be 10. Got %d", row[3])
		}
		if row[4] != 5 {
			t.Errorf("expected row[4] to be 5. Got %d", row[4])
		}
		if row[5] != "ages ago" {
			t.Errorf("expected row[5] to be ages ago. Got %s", row[5])
		}
	})

	t.Run("Youtube Track to Row", func(t *testing.T) {
		yt := YoutubeTrackRow{
			Title:            "little things",
			Artist:           "adrianne",
			Found:            true,
			Link:             "https://www.youtube.com/watch?v=123",
			VideoId:          "123",
			VideoPublishDate: "2025-01-01",
			AddedAt:          "ages ago",
		}

		row := YoutubeTrackToRow(yt)

		if len(row) != 7 {
			t.Errorf("expected row to have 7 elements. Got %d", len(row))
		}

		if row[0] != "little things" {
			t.Errorf("expected row[0] to be little things. Got %s", row[0])
		}
		if row[1] != "adrianne" {
			t.Errorf("expected row[1] to be adrianne. Got %s", row[1])
		}
		if row[2] != true {
			t.Errorf("expected row[2] to be true. Got %t", row[2])
		}
		if row[3] != "https://www.youtube.com/watch?v=123" {
			t.Errorf("expected row[3] to be https://www.youtube.com/watch?v=123. Got %s", row[3])
		}
		if row[4] != "123" {
			t.Errorf("expected row[4] to be 123. Got %s", row[4])
		}
		if row[5] != "2025-01-01" {
			t.Errorf("expected row[5] to be 2025-01-01. Got %s", row[5])
		}
		if row[6] != "ages ago" {
			t.Errorf("expected row[6] to be ages ago. Got %s", row[6])
		}
	})

	t.Run("Row to Youtube Track", func(t *testing.T) {
		r := make([]interface{}, 7)
		r[0] = "little things"
		r[1] = "adrianne"
		r[2] = "TRUE"
		r[3] = "https://www.youtube.com/watch?v=123"
		r[4] = "123"
		r[5] = "2025-01-01"
		r[6] = "ages ago"

		yt := RowToYoutubeTrack(r)

		if yt.Title != "little things" {
			t.Errorf("expected title to be little things. Got %s", yt.Title)
		}
		if yt.Artist != "adrianne" {
			t.Errorf("expected artist to be adrianne. Got %s", yt.Artist)
		}
		if yt.Found != true {
			t.Errorf("expected found to be true. Got %t", yt.Found)
		}
		if yt.Link != "https://www.youtube.com/watch?v=123" {
			t.Errorf("expected link to be https://www.youtube.com/watch?v=123. Got %s", yt.Link)
		}
		if yt.VideoId != "123" {
			t.Errorf("expected videoId to be 123. Got %s", yt.VideoId)
		}
		if yt.VideoPublishDate != "2025-01-01" {
			t.Errorf("expected videoPublishDate to be 2025-01-01. Got %s", yt.VideoPublishDate)
		}
		if yt.AddedAt != "ages ago" {
			t.Errorf("expected addedAt to be ages ago. Got %s", yt.AddedAt)
		}
	})
}
