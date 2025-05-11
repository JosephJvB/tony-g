package spotify

import (
	"slices"
	"testing"
)

func TestSpotifyModels(t *testing.T) {
	t.Run("can remove feature from \"Misery (feat. Kenny Segal)\"", func(t *testing.T) {
		title := "Misery (feat. Kenny Segal)"

		title = CleanSongTitle(title)

		if title != "Misery" {
			t.Errorf("Expected trimmed title to be \"Misery\". Received \"%s\"", title)
		}
	})

	t.Run("can remove feature from \"Flood (feat. Obongjayar & Moonchild Sanelly)\"", func(t *testing.T) {
		title := "Flood (feat. Obongjayar & Moonchild Sanelly)"

		title = CleanSongTitle(title)

		if title != "Flood" {
			t.Errorf("Expected trimmed title to be \"Flood\". Received \"%s\"", title)
		}
	})

	t.Run("can remove feature from \"Too Fast (Pull Over) [feat. Latto]\"", func(t *testing.T) {
		title := "Too Fast (Pull Over) [feat. Latto]"

		title = CleanSongTitle(title)

		if title != "Too Fast (Pull Over)" {
			t.Errorf("Expected trimmed title to be \"Too Fast (Pull Over)\". Received \"%s\"", title)
		}
	})

	t.Run("can remove feature from \"New Level (Remix) ft. A$AP Rocky, Future, Lil Uzi Vert\"", func(t *testing.T) {
		title := "New Level (Remix) ft. A$AP Rocky, Future, Lil Uzi Vert"

		title = CleanSongTitle(title)

		if title != "New Level (Remix)" {
			t.Errorf("Expected trimmed title to be \"New Level (Remix)\". Received \"%s\"", title)
		}
	})

	t.Run("clean song title \"Hit Me Where It Hurts (Toro Y Moi Remix) ft. Chino Moreno by Caroline Polacheck\"", func(t *testing.T) {
		title := "Hit Me Where It Hurts (Toro Y Moi Remix) ft. Chino Moreno by Caroline Polacheck"

		title = CleanSongTitle(title)

		if title != "Hit Me Where It Hurts (Toro Y Moi Remix)" {
			t.Errorf("Expected trimmed title to be \"Hit Me Where It Hurts (Toro Y Moi Remix)\". Received \"%s\"", title)
		}
	})

	t.Run("turns spotify link to uri", func(t *testing.T) {
		input := "https://open.spotify.com/track/0jv5VgdENAPV7lHtBlsaXE"

		result, ok := LinkToTrackUri(input)

		if !ok {
			t.Error("linkToTrackUri failed")
		}

		if result != "spotify:track:0jv5VgdENAPV7lHtBlsaXE" {
			t.Errorf("Failed to get uri from link. Received %s", result)
		}
	})

	type DateObj struct {
		date string
	}
	t.Run("sort array of dates", func(t *testing.T) {
		dates := []DateObj{
			{
				date: "2025-05-05T04:34:10Z",
			},
			{
				date: "2024-04-28T04:29:04Z",
			},
			{
				date: "2024-03-28T04:29:04Z",
			},
			{
				date: "2024-03-22T04:29:04Z",
			},
			{
				date: "2016-03-22T04:29:04Z",
			},
		}

		slices.SortFunc(dates, func(a, z DateObj) int {
			if a.date < z.date {
				return -1
			}
			if a.date > z.date {
				return 1
			}
			return 0
		})

		if dates[0].date != "2016-03-22T04:29:04Z" {
			t.Errorf("Expected oldest first")
		}
		if dates[1].date != "2024-03-22T04:29:04Z" {
			t.Errorf("Expected second second")
		}
		if dates[2].date != "2024-03-28T04:29:04Z" {
			t.Errorf("Expected third third")
		}
		if dates[3].date != "2024-04-28T04:29:04Z" {
			t.Errorf("Expected fourth fourth")
		}
		if dates[4].date != "2025-05-05T04:34:10Z" {
			t.Errorf("Expected fifth fifth")
		}
	})
}
