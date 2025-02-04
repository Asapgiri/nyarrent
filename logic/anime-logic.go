package logic

import (
	"encoding/json"
	"io"
	"math"
	"net/http"
	"nyarrent/dbase"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/er-azh/go-animeschedule"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// =====================================================================================================================
// Anime related internal logic...

func ListAllAnime(forceUpdate bool) DtoAnime {
    anime := dbase.Anime{}
    list, _ := anime.List()

    lAnime := make([]Anime, len(list))

    for i, a := range(list) {
        if forceUpdate {
            a, _ = AddorUpdateAnime(a.Route)
        } else {
            calculateCurrentEpisodeFromTime(&a, a.SubTime)
            a.Update()
        }
        lAnime[i].Anime = a
    }

    slices.SortFunc(lAnime, func(a, b Anime) int {
        return b.Anime.EpisodeRelease.Compare(a.Anime.EpisodeRelease)
    })

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

var AS_API_URL = "https://animeschedule.net/api/v3"
var AS_API_KEY = os.Getenv("AS_API_KEY")

func aSHttpGet(query string) (*http.Response, error) {
    req, err := http.NewRequest("GET", AS_API_URL + query, nil)
    if nil != err {
        return nil, err
    }
    if "" != AS_API_KEY {
        req.Header.Add("Authorization", "Bearer " + AS_API_KEY)
    }
    client := http.Client{}
    return client.Do(req)
}

func FindNewAnimes(query string, page string) AnimeSearchPage {
    var aniList = AnimeSearchPage{}
    q := strings.ReplaceAll(query, " ", "+")

    resp, err := aSHttpGet("/anime?page="+page+"&q="+q)
    if nil != err {
        log.Println(err.Error())
        return AnimeSearchPage{}
    }
    defer resp.Body.Close()

    aniListJsonBarr, err := io.ReadAll(resp.Body)
    if nil != err {
        log.Println(err.Error())
        return AnimeSearchPage{}
    }

    err = json.Unmarshal(aniListJsonBarr, &aniList)
    if nil != err {
        log.Println(resp)
        log.Println(err.Error())
        return AnimeSearchPage{}
    }

    aniList.SearchText = query
    aniList.Added = make([]bool, len(aniList.Anime))

    dbanime := dbase.Anime{}
    list, _ := dbanime.List()

    for _, l := range(list) {
        for i, anime := range(aniList.Anime) {
            if anime.Route == l.Route {
                aniList.Added[i] = true
            }
        }
    }

    return aniList
}


func AddorUpdateAnime(route string) (dbase.Anime, error) {
    anime := animeschedule.ShowDetail{}

    resp, err := aSHttpGet("/anime/"+route)
    if nil != err {
        log.Println(err.Error())
        return dbase.Anime{}, err
    }
    defer resp.Body.Close()

    aniListJsonBarr, err := io.ReadAll(resp.Body)
    if nil != err {
        log.Println(err.Error())
        return dbase.Anime{}, err
    }

    err = json.Unmarshal(aniListJsonBarr, &anime)
    if nil != err {
        log.Println(err.Error())
        return dbase.Anime{}, err
    }

    dbAnime := dbase.Anime{}
    err = dbAnime.Select(route)
    if nil != err {
        dbAnime.Id = primitive.NewObjectID()
    }

    dbAnime.Title =         anime.Title
    dbAnime.JpTitle =       anime.Names.Japanese
    dbAnime.Route =         anime.Route
    dbAnime.Banner =        anime.ImageVersionRoute
    dbAnime.EpisodeCount =  anime.Episodes
    dbAnime.Status =        anime.Status
    dbAnime.JpnTime =       anime.JpnTime
    dbAnime.SubTime =       anime.SubTime
    dbAnime.FullInfo =      anime

    calculateCurrentEpisodeFromTime(&dbAnime, dbAnime.SubTime)

    if nil != err {
        return dbAnime, dbAnime.Add()
    } else {
        return dbAnime, dbAnime.Update()
    }
}

func ListAnime(route string) Anime {
    dbAnime := dbase.Anime{}
    dbAnime.Select(route)
    calculateCurrentEpisodeFromTime(&dbAnime, dbAnime.SubTime)

    episodes := make([]Episodes, dbAnime.EpisodeCurrent)
    dbEpisodes := dlSelectAll(dbAnime)

    log.Println(dbEpisodes)

    for _, episode := range(dbEpisodes) {
        info := getTorrentInfo(episode.Hash)
        var percent float64
        if 0 != info.SizeWhenDone {
            percent = float64(info.HaveValid) / float64(info.SizeWhenDone)
        } else {
            percent = 0.0
        }

        if dbAnime.EpisodeCurrent >= episode.Episode {
            episodes[episode.Episode - 1].Title = episode.Title
            episodes[episode.Episode - 1].Torrents =
                append(episodes[episode.Episode - 1].Torrents, EpisodeTorrent{
                    Torrent: episode,
                    Info: getTorrentInfo(episode.Hash),
                    Progress: GenerateProgress(percent),
                    Url: GetTorrentFile(info.Name, true),
                })
        }
    }

    for i, _ := range(episodes) {
        idx := i + 1
        episodes[i].Index = idx
        if "" == episodes[i].Title {
            episodes[i].Title = "null"
        }
        if 0 == len(episodes[i].Torrents) {
            //log.Printf("Getting nyaa for episode %s - ep%d\n", dbAnime.Title, idx)
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

func calculateCurrentEpisodeFromTime(anime *dbase.Anime, startTime time.Time) {
    const WEEK_HOUR_DIFF = 24 * 7
    diff := time.Now().Sub(startTime)
    weeks := int(math.Floor(diff.Hours() / WEEK_HOUR_DIFF))

    if anime.EpisodeCount < 1 + weeks {
        weeks = anime.EpisodeCount - 1
    }

    // TODO: Calculate delays and dates stuff correctly...

    anime.EpisodeCurrent = 1 + weeks
    anime.EpisodeRelease = startTime.Add(time.Hour * time.Duration(weeks * WEEK_HOUR_DIFF))
}
