package apple

import (
	"testing"
)

func TestAppleModels(t *testing.T) {
	t.Run("apple scraped track makes expected id", func(t *testing.T) {
		expectedId := "sea of trees__king gizz__12 bb__2012"

		scrapedTrack := ScrapedTrack{
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
}
