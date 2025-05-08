package apple

import (
	"testing"
	"tony-g/internal/googlesheets"
)

func TestAppleModels(t *testing.T) {
	t.Run("apple scraped track id matches apple track row", func(t *testing.T) {
		scrapedTrack := ScrapedTrack{
			Title:  "Sea of Trees",
			Artist: "King Gizz",
			Album:  "12 BB",
			Year:   2012,
		}
		rowTrack := googlesheets.AppleTrackRow{
			Title:  "Sea of Trees",
			Artist: "King Gizz",
			Album:  "12 BB",
			Year:   2012,
		}

		a := scrapedTrack.GetAppleTrackId()
		b := rowTrack.GetAppleTrackId()
		if a != b {
			t.Errorf("expected %s to equal %s", a, b)
		}
	})
}
