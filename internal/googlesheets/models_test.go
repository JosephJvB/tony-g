package googlesheets

import (
	"testing"
)

func TestYoutubeModels(t *testing.T) {
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
}
