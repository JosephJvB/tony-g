package apple

import (
	"tony-g/internal/stringutil"
)

type ScrapedTrack struct {
	Title      string
	Artist     string
	Album      string
	DurationMs int
	Year       int
}

func (st *ScrapedTrack) GetAppleTrackId() string {
	idParts := stringutil.AppleIdParts{
		Title:  st.Title,
		Album:  st.Album,
		Artist: st.Artist,
		Year:   st.Year,
	}

	return stringutil.MakeAppleId(idParts)
}
