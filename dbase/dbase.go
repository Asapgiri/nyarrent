package dbase

import (
	"context"
	"errors"
	"nyarrent/logger"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoDBuri = os.Getenv("NYARRENT_URI")
var mongoDBdatabase = os.Getenv("NYARRENT_DATABASE_NAME")
var mongo_client *mongo.Client
var db *mongo.Database

var dbSAVED_ANIMES          *mongo.Collection
var dbDOWNLOADED_EPISODES   *mongo.Collection
var dbNYAA_CACHE            *mongo.Collection

var log = logger.Logger {
    Color: logger.Colors.Purple,
    Pretext: "database",
}

// =====================================================================================================================
// Basic connect and stuff

func check_envs(envs ...string) error {
    for i, e := range envs {
        // FIXME: Remove logs after secrets put in place and debug mode is off
        log.Printf("%02d: Loading DB env: %s=%s\n", i + 1, e, os.Getenv(e))
        if "" == os.Getenv(e) {
            msg := "Environmental variable does not exists: " + e
            log.Println(msg)
            return errors.New(msg)
        }
    }

    return nil
}

func Connect() error {
    var err error
    err = check_envs("NYARRENT_URI", "NYARRENT_DATABASE_NAME")
    if nil != err {
        return err
    }

    // Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoDBuri).SetServerAPIOptions(serverAPI)

    // Create a new client and connect to the server
    mongo_client, err = mongo.Connect(context.TODO(), opts)
	if err != nil {
        return err
	}
    db = mongo_client.Database(mongoDBdatabase)

	// Send a ping to confirm a successful connection
	var result bson.M
	if err := db.RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		panic(err)
	}
	log.Println("Pinged your deployment. You successfully connected to MongoDB!")

    dbSAVED_ANIMES = db.Collection("saved-anime")
    dbDOWNLOADED_EPISODES = db.Collection("downloaded-episodes")
    dbNYAA_CACHE = db.Collection("nyaa-cache")

    return nil
}

// =====================================================================================================================
// Internal Anime Listing CRUD

func (anime *Anime) List() ([]Anime, error) {
    var anyime []Anime

    cursor, err := dbSAVED_ANIMES.Find(context.TODO(), bson.D{{}})
    if nil != err {
        return anyime, err
    }

    err = cursor.All(context.TODO(), &anyime)

    return anyime, err
}

func (anime *Anime) Select(route string) error {
    filter := bson.D{{"route", route}}
    err := dbSAVED_ANIMES.FindOne(context.TODO(), filter).Decode(anime)
    return err
}

func (anime *Anime) Add() error {
    _, err := dbSAVED_ANIMES.InsertOne(context.TODO(), anime)
    return err
}

func (anime *Anime) Update() error {
    _, err := dbSAVED_ANIMES.ReplaceOne(context.TODO(), bson.D{{"_id", anime.Id}}, anime)
    return err

}

func (anime *Anime) Delete() error {
    filter := bson.D{{"_id", anime.Id}}
    _, err := dbSAVED_ANIMES.DeleteOne(context.TODO(), filter)
    return err
}


// =====================================================================================================================
// Downloads episode Listing CRUD

func (danime *AnimeDownload) List(animeId primitive.ObjectID) ([]AnimeDownload, error) {
    var downloads []AnimeDownload

    filter := bson.D{{"animeid", animeId}}
    cursor, err := dbDOWNLOADED_EPISODES.Find(context.TODO(), filter)
    if nil != err {
        return downloads, err
    }

    err = cursor.All(context.TODO(), &downloads)

    return downloads, err
}

func (danime *AnimeDownload) Select(id primitive.ObjectID) error {
    filter := bson.D{{"_id", id}}
    err := dbDOWNLOADED_EPISODES.FindOne(context.TODO(), filter).Decode(danime)
    return err
}

func (danime *AnimeDownload) Add() error {
    _, err := dbDOWNLOADED_EPISODES.InsertOne(context.TODO(), danime)
    return err
}

func (danime *AnimeDownload) Update() error {
    _, err := dbDOWNLOADED_EPISODES.ReplaceOne(context.TODO(), bson.D{{"_id", danime.Id}}, danime)
    return err

}

func (danime *AnimeDownload) Delete() error {
    filter := bson.D{{"_id", danime.Id}}
    _, err := dbDOWNLOADED_EPISODES.DeleteOne(context.TODO(), filter)
    return err
}

// =====================================================================================================================
// Nyaa CRUD

func (nyaached *NyaaCached) Select(title string, episode int) error {
    filter := bson.D{{"episode", episode}, {"title", title}}
    err := dbNYAA_CACHE.FindOne(context.TODO(), filter).Decode(nyaached)
    return err
}

func (nyaached *NyaaCached) Add() error {
    _, err := dbNYAA_CACHE.InsertOne(context.TODO(), nyaached)
    return err
}

func (nyaached *NyaaCached) Update() error {
    _, err := dbNYAA_CACHE.ReplaceOne(context.TODO(), bson.D{{"_id", nyaached.Id}}, nyaached)
    return err

}

func (nyaached *NyaaCached) Delete() error {
    filter := bson.D{{"_id", nyaached.Id}}
    _, err := dbNYAA_CACHE.DeleteOne(context.TODO(), filter)
    return err
}
