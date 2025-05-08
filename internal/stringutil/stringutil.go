package stringutil

import (
	"strconv"
	"strings"
)

type AppleIdParts struct {
	Title  string
	Artist string
	Album  string
	Year   int
}

func MakeAppleId(parts AppleIdParts) string {
	id := ""

	idParts := []string{
		parts.Title,
		parts.Artist,
		parts.Album,
		strconv.Itoa(parts.Year),
	}

	for _, p := range idParts {
		if id != "" {
			id += "__"
		}

		id += strings.ToLower(strings.TrimSpace(p))
	}

	return id
}
