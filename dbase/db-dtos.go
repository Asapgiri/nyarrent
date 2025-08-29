package dbase

import (
	"time"

	"github.com/er-azh/go-animeschedule"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Anime struct {
    Id              primitive.ObjectID `bson:"_id"`
    Title           string
    JpTitle         string
    Route           string
    Banner          string
    EpisodeCurrent  int
    EpisodeRelease  time.Time
    EpisodeCount    int
    Status          string
    JpnTime         time.Time
    SubTime         time.Time
    FullInfo        animeschedule.ShowDetail
}

type AnimeDownload struct {
    Id              primitive.ObjectID `bson:"_id"`
    AnimeId         primitive.ObjectID
    Episode         int
    Title           string
    Link            string
    Hash            string
}

type NyaaCached struct {
    Id              primitive.ObjectID `bson:"_id"`
    Episode         int
    Title           string
    Nyaa            NyaaJson
	Shoto			[]AnimeShotoTorrent
}

// Nyaa.si stuff

type AnimeShotoTorrent struct {
    Id                          int
    Title                       string
    Link                        string
    Timestamp                   int
    Status                      string
    Tosho_id                    any
    Nyaa_id                     int
    Nyaa_subdom                 any
    Anidex_id                   any
    Torrent_url                 string
    Torrent_name                string
    Info_hash                   string
    Info_hash_v2                any
    Magnet_uri                  string
    Seeders                     int
    Leechers                    int
    Torrent_downloaded_count    int
    Tracker_updated             int
    Nzb_url                     string
    Total_size                  int
    SizeStr                     string
    Num_files                   int
    Anidb_aid                   int
    Anidb_eid                   int
    Anidb_fid                   any
    Article_url                 any
    Article_title               any
    Website_url                 any
}

type NyaaData struct {
    Category    string
    Title       string
    Link        string
    Torrent     string
    Magnet      string
    Size        string
    Time        string
    Seeders     int
    Leechers    int
    Downloads   int
}

type NyaaJson struct {
    Count   int
    Data    []NyaaData
}
