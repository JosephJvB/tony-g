package spotify

import "testing"

func TestUtil(t *testing.T) {
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
}
