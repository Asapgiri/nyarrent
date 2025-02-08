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
    Hash        string
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
    NyaaText    string
}

type NyaaFilter struct {
    Group       string
    NameParams  string
    Category    string
    SubCategory string
    ResultCount string
    Resolution  string
    EpisodeFmt  string
    SeasonFmt   string
    ForseRefrsh bool
}

type EpisodeFilter struct {
    Nyaa    NyaaFilter
    Hash    string
}

type Anime struct {
    Anime       dbase.Anime
    Episodes    []Episodes
    Filter      EpisodeFilter
}

type DtoAnime struct {
    SearchText  string
    Anime       []Anime
}
