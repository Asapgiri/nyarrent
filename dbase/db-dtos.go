package dbase

import (
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
    EpisodeCount    int
    Status          string
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
}

// Nyaa.si stuff

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
