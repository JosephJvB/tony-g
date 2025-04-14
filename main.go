package main

import (
	"fmt"
	"tony-gony/internal/googlesheets"
	"tony-gony/internal/youtube"
)

func main() {
	yt := youtube.NewClient()
	gs := googlesheets.NewClient()

	yt.LoadPlaylistItems("")
	gs.LoadParsedVideos()

	toParse := []youtube.PlaylistItem{}
	for _, video := range yt.PlaylistItems {
		if !gs.ParsedVideosMap[video.Id] {
			toParse = append(toParse, video)
		}
	}

	if len(toParse) == 0 {
		fmt.Println("No new videos to parse. Exiting")
		return
	}
}
