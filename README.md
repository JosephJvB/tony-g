# Tony Gony

rebuild Tony with Go

Gonna end up deploying him to AWS as per. So I'll use this repo as ref https://github.com/JosephJvB/spotify-users-backend

need:

- aws parameter store api
- youtube api
- google sheets api
- spotify api

steps:

- load secrets from parameter store
  - currently other clients are reading from os.Getenv()
  - I could have param store do os.Setenv() if I wanted
  - but I think I'd rather pass from paramClient into other clients on creation
- load youtube items
- load google sheets
- load existing playlists
  - don't need their items immediately
  - only need to know their items if we are adding to those playlists

https://edu.anarcho-copy.org/Programming%20Languages/Go/Concurrency%20in%20Go.pdf

go clean -testcache: expires all test results

hey actually. Tony and his mates keep an up to date Apple playlist no?
https://music.apple.com/us/playlist/my-fav-singles-of-2025/pl.u-ayeZTygbKDy

Why not use that as the source, rather than youtube video descriptions?

Better flow:
playlist description says it gets updated every friday
so every saturday
get tony's playlists from apple music
find the current one by title (current year)
get all the playlists songs from apple music
(filter for just the recent ones?)

get my spotify playlist for the same year (create if not exists)
get all songs currently in my playlist

find the tracks that need to be added (song_artist_album)
find those tracks in spotify
add them to playlist
Seems easy enough?

nah wait apple music api is garbage

lets do it this way:
scrape https://theneedledrop.com/loved-list/${year}

- years >= 2022 use apple playlists
- years < 2022 use spotify or nothing at all
- with this one, I'm gonna just handle future playlists
  - maybe recreate playlists from 2022 onwards too with new name and decommish the old service

find the apple music link in html:
https://embed.music.apple.com/us/playlist/my-fav-singles-of-2024/pl.u-e2ZmtK9VM5K?wmode=opaque
https://music.apple.com/us/playlist/my-fav-singles-of-2024/pl.u-e2ZmtK9VM5K

scrape the apple music playlist page for tracks: songname, artist, album

get google sheets tracks that I've already tried to add

find those tracks in spotify that I haven't tried to add

get my spotify playlists and their items

search in spotify for those tracks

rather than scraping the playlist url every week should I save it in google sheets?
We'll see

deleted tony2 stack cos it would have found the new Tony #2 playlists due to prefix overlap
I can change the prefix maybe but should be fine now that the old fn is not running

just noticed too actually that the curated Apple Playlists are a lot shorter than the ones I made
see: https://music.apple.com/us/playlist/my-fav-singles-of-2024/pl.u-e2ZmtK9VM5K
vs: https://open.spotify.com/playlist/3cIeEpjP3PhNiFD6aKfyD6?si=a277b1e20a8043f7
ie: Apple Playlists don't have every liked song from Youtube Best Tracks
eep!
Maybe I should keep my guy running!

1. turn tony2 back online (delete web assets stack tho that's garbage)
2. tony-gony needs a new prefix "2onys 2op 2racks 2024" ?

- aso I should rename tony-gony to tony-g that's way better
