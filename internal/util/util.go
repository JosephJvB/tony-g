package util

import (
	"regexp"
	"strings"
)

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

// https://stackoverflow.com/questions/4292468/javascript-regex-remove-text-between-parentheses#answer-4292483
// .replace(/*\([^)]*\)*/g, ”)
// maybe need to handle "ft." ? but this seems enough for now
func RemoveFeatureString(songTitle string) string {
	rmParens := regexp.MustCompile(`\\*\(feat.[^)]*\)*`)
	rmSquareBrackets := regexp.MustCompile(`\\*\[feat.[^)]*\]*`)
	songTitle = rmParens.ReplaceAllLiteralString(songTitle, "")
	songTitle = rmSquareBrackets.ReplaceAllLiteralString(songTitle, "")

	return strings.TrimSpace(songTitle)
}
