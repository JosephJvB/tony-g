package util

import "strings"

type IdParts struct {
	Title  string
	Artist string
	Album  string
	Year   string
}

func MakeTrackId(parts IdParts) string {
	id := ""

	partsList := []string{
		parts.Title,
		parts.Artist,
		parts.Album,
		parts.Year,
	}

	for _, p := range partsList {
		p = strings.ToLower(strings.TrimSpace(p))

		if id != "" {
			id += "__"
		}

		id += p
	}

	return id
}
