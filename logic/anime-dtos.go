package logic

import (
	"nyarrent/dbase"

	"github.com/er-azh/go-animeschedule"
)

type AnimeSearchPage struct {
    Page        int
    TotalAmount int
    SearchText  string
    Anime       []animeschedule.ShowDetail
    Added       []bool
}

type TimetableShow struct {
    Anime   animeschedule.TimetableShow
    Added   bool
    Filled  bool
    Aired   bool
}

type AnimeWeek [7]TimetableShow

type AnimeTimetableFilter struct {
    OnlyOnList  bool
    SendBack    bool
}

type AnimeTimetablePage struct {
    AnimeWeek   []AnimeWeek
    Filter      AnimeTimetableFilter
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
    Nyaa        []dbase.NyaaData
}

type Anime struct {
    Anime       dbase.Anime
    Episodes    []Episodes
}

type DtoAnime struct {
    SearchText  string
    Anime       []Anime
}
