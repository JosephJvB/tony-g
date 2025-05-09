package googlesearch

import "testing"

func TestGoogleSearch(t *testing.T) {
	t.Run("can find clipping blood on the fang", func(t *testing.T) {
		query := "spotify track clipping blood on the fang"

		Search(query)
	})
}
