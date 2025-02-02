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
