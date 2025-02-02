package logic

import (
	"nyarrent/dbase"

	"github.com/er-azh/go-animeschedule"
)

type AnimeSearchPage struct {
    Page        int
    TotalAmount int
    Anime       []animeschedule.ShowDetail
}

type EpisodeTorrent struct {
    Torrent     dbase.AnimeDownload
    Info        TorrentInfo
    Progress    Progress
    Url         string
}

type Episodes struct {
    Index       int
    Title       string
    Torrents    []EpisodeTorrent
}

type Anime struct {
    Anime       dbase.Anime
    Episodes    []Episodes
}

type DtoAnime struct {
    SearchText  string
    Anime       []Anime
}
