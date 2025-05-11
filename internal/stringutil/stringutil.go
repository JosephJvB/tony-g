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
	idParts := []string{
		parts.Title,
		parts.Artist,
		parts.Album,
		strconv.Itoa(parts.Year),
	}

	return makeId(idParts)
}

func makeId(parts []string) string {
	id := ""

	for _, p := range parts {
		if id != "" {
			id += "__"
		}

		id += strings.ToLower(strings.TrimSpace(p))
	}

	return id
}
