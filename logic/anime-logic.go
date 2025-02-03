package logic

import (
	"encoding/json"
	"io"
	"net/http"
	"nyarrent/dbase"
	"slices"
	"strconv"

	"github.com/er-azh/go-animeschedule"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// =====================================================================================================================
// Anime related internal logic...

func listAllAnime() DtoAnime {
    anime := dbase.Anime{}
    list, _ := anime.List()

    lAnime := make([]Anime, len(list))

    for i, a := range(list) {
        lAnime[i].Anime = a

        episodes := dlSelectAll(a)
        lAnime[i].Episodes = make([]Episodes, a.EpisodeCurrent)
        for j, _ := range(lAnime[i].Episodes) {
            lAnime[i].Episodes[j].Index = i + 1
            lAnime[i].Episodes[j].Title = "null"
        }
        for _, episode := range(episodes) {
            info := getTorrentInfo(episode.Hash)
            var percent float64
            if 0 == info.SizeWhenDone {
                percent = 0.0
            } else {
                percent = float64(info.HaveValid) / float64(info.SizeWhenDone)
            }

            lAnime[i].Episodes[episode.Episode - 1].Title = episode.Title
            lAnime[i].Episodes[episode.Episode - 1].Torrents =
                append(lAnime[i].Episodes[episode.Episode - 1].Torrents, EpisodeTorrent{
                    Torrent: episode,
                    Info: info,
                    Progress: GenerateProgress(percent),
                    Url: GetTorrentFile(info.Name, true),
                })
        }
    }

    ret := DtoAnime{
        Anime: lAnime,
    }

    return ret
}

func deleteAnime() error {
    return nil
}

// =====================================================================================================================
// Anime related internal downloads logic...

func dlSelectAll(anime dbase.Anime) []dbase.AnimeDownload {
    var adl dbase.AnimeDownload
    list, _ := adl.List(anime.Id)
    return list
}

func dlSelect(id primitive.ObjectID) dbase.AnimeDownload { return dbase.AnimeDownload{} }

func dlAdd(anime dbase.Anime, dl dbase.AnimeDownload) error { return nil }
func dlDelete(dl dbase.AnimeDownload) error { return nil }

// =====================================================================================================================
// Anime related internal downloads logic...

func AddAnime(route string) error {
    anime := animeschedule.ShowDetail{}

    resp, err := http.Get("https://animeschedule.net/api/v3/anime/"+route)
    if nil != err {
        log.Println(err.Error())
        return err
    }
    defer resp.Body.Close()

    aniListJsonBarr, err := io.ReadAll(resp.Body)
    if nil != err {
        log.Println(err.Error())
        return err
    }

    err = json.Unmarshal(aniListJsonBarr, &anime)
    if nil != err {
        log.Println(err.Error())
        return err
    }

    log.Println(anime)

    dbAnime := dbase.Anime{
        Id:             primitive.NewObjectID(),
        Title:          anime.Title,
        JpTitle:        anime.Names.Japanese,
        Route:          anime.Route,
        Banner:         anime.ImageVersionRoute,
        EpisodeCurrent: anime.Episodes,
        EpisodeCount:   12,
        Status:         anime.Status,
        FullInfo:       anime,
    }

    return dbAnime.Add()
}

func ListAnime(route string) Anime {
    dbAnime := dbase.Anime{}
    dbAnime.Select(route)

    episodes := make([]Episodes, dbAnime.EpisodeCurrent)
    dbEpisodes := dlSelectAll(dbAnime)

    for _, episode := range(dbEpisodes) {
        info := getTorrentInfo(episode.Hash)
        percent := float64(info.HaveValid) / float64(info.SizeWhenDone)

        episodes[episode.Episode - 1].Title = episode.Title
        episodes[episode.Episode - 1].Torrents =
            append(episodes[episode.Episode - 1].Torrents, EpisodeTorrent{
                Torrent: episode,
                Info: getTorrentInfo(episode.Hash),
                Progress: GenerateProgress(percent),
                Url: GetTorrentFile(info.Name, true),
            })
    }

    for i, _ := range(episodes) {
        idx := i + 1
        episodes[i].Index = idx
        if "" == episodes[i].Title {
            episodes[i].Title = "null"
        }
        if 0 == len(episodes[i].Torrents) {
            log.Printf("Getting nyaa for episode %s - ep%d\n", dbAnime.Title, idx)
            episodes[i].Nyaa = GetNyaaList(dbAnime.Title, idx, 10, false)
        }
    }


    slices.Reverse(episodes)

    anime := Anime{
        Anime: dbAnime,
        Episodes: episodes,
    }

    return anime
}

func AddEpisode(route string, index string, title string, link string, hash string) error {
    anime := dbase.Anime{}
    err := anime.Select(route)
    if nil != err {
        return err
    }

    i, _ := strconv.ParseInt(index, 10, 64)
    dbEpisode := dbase.AnimeDownload{
        Id:         primitive.NewObjectID(),
        AnimeId:    anime.Id,
        Episode:    int(i),
        Title:      title,
        Link:       link,
        Hash:       hash,
    }
    return dbEpisode.Add()
}
