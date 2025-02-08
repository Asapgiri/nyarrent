package logic

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"nyarrent/config"
	"nyarrent/dbase"
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
            getCurrentEpisode(&a)
            a.Update()
        }
        lAnime[i].Anime = a
    }

    slices.SortFunc(lAnime, func(a, b Anime) int {
        if "Ongoing" != a.Anime.Status && "Ongoing" == b.Anime.Status {
            return 1
        } else if "Ongoing" == a.Anime.Status && "Ongoing" != b.Anime.Status {
            return -1
        } else if time.Now().Compare(a.Anime.EpisodeRelease) <= 0 && time.Now().Compare(b.Anime.EpisodeRelease) > 0 {
            return 1
        } else if time.Now().Compare(a.Anime.EpisodeRelease) > 0 && time.Now().Compare(b.Anime.EpisodeRelease) <= 0 {
            return -1
        } else {
            return b.Anime.EpisodeRelease.Compare(a.Anime.EpisodeRelease)
        }
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

func aSHttpGet(query string) (*http.Response, error) {
    req, err := http.NewRequest("GET", config.Config.AnimeAPI.ApiUrl + query, nil)
    if nil != err {
        return nil, err
    }
    if "" != config.Config.AnimeAPI.ApiKey {
        req.Header.Add("Authorization", "Bearer " + config.Config.AnimeAPI.ApiKey)
    }
    client := http.Client{}
    return client.Do(req)
}

func sJsonHttpUnmarshall(object interface{}, query string) error {
    resp, err := aSHttpGet(query)
    if nil != err {
        log.Println(err.Error())
        return err
    }
    defer resp.Body.Close()

    objectJsonBarr, err := io.ReadAll(resp.Body)
    if nil != err {
        log.Println(err.Error())
        return err
    }

    err = json.Unmarshal(objectJsonBarr, &object)
    if nil != err {
        log.Println(resp)
        log.Println(err.Error())
        return err
    }

    return nil
}

func getLongestDayShowCount(timetable animeschedule.Timetable) ([7]int, int) {
    weekdays := [7]int{}
    dayCount := 0
    current := timetable[0].EpisodeDate

    for _, anime := range(timetable) {
        if current.Day() != anime.EpisodeDate.Day() {
            for current.Day() != anime.EpisodeDate.Day() {
                current = current.Add(time.Hour * 24)
                dayCount++
            }
        }

        weekdays[dayCount]++
    }

    highestCount := weekdays[0]
    for _, count := range(weekdays) {
        if count > highestCount {
            highestCount = count
        }
    }

    return weekdays, highestCount
}

func ListTimetables(filter AnimeTimetableFilter) AnimeTimetablePage {
    timetable := animeschedule.Timetable{}

    err := sJsonHttpUnmarshall(&timetable, "/timetables/sub")
    if nil != err {
        return AnimeTimetablePage{}
    }

    dbanime := dbase.Anime{}
    list, _ := dbanime.List()

    if (filter.OnlyOnList) {
        filteredTimetable := animeschedule.Timetable{}
        for _, anime := range(timetable) {
            for _, l := range(list) {
                if anime.Route == l.Route {
                    filteredTimetable = append(filteredTimetable, anime)
                }
            }
        }
        timetable = filteredTimetable
    }

    weekdays, maxShowPerDay := getLongestDayShowCount(timetable)
    attPage := AnimeTimetablePage{
        AnimeWeek: make([]AnimeWeek, maxShowPerDay),
        Filter: filter,
    }

    ttCounter := 0
    dtNow := time.Now()
    for day, days := range(weekdays) {
        for i := 0; i < days; i++ {
            attPage.AnimeWeek[i][day].Anime = timetable[ttCounter + i]
            attPage.AnimeWeek[i][day].Filled = true
            attPage.AnimeWeek[i][day].Aired = timetable[ttCounter + i].EpisodeDate.Compare(dtNow) <= 0
            for _, l := range(list) {
                if attPage.AnimeWeek[i][day].Anime.Route == l.Route {
                    attPage.AnimeWeek[i][day].Added = true
                }
            }
        }
        ttCounter += days
    }

    return attPage
}

func FindNewAnimes(query string, page string) AnimeSearchPage {
    var aniList = AnimeSearchPage{}
    q := strings.ReplaceAll(query, " ", "+")

    err := sJsonHttpUnmarshall(&aniList, "/anime?page="+page+"&q="+q)
    if nil != err {
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

    err := sJsonHttpUnmarshall(&anime, "/anime/"+route)

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

    getCurrentEpisode(&dbAnime)

    if nil != err {
        return dbAnime, dbAnime.Add()
    } else {
        return dbAnime, dbAnime.Update()
    }
}

func ListAnime(route string, filter *EpisodeFilter) Anime {
    dbAnime := dbase.Anime{}
    dbAnime.Select(route)
    getCurrentEpisode(&dbAnime)

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
            episodes[i].NyaaText, episodes[i].Nyaa = GetNyaaList(dbAnime.Title, idx, filter.Nyaa)
        }
    }

    slices.Reverse(episodes)
    // Doing it at every load would be annoying...
    filter.Nyaa.ForseRefrsh = false

    anime := Anime{
        Anime: dbAnime,
        Episodes: episodes,
        Filter: *filter,
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

    err = dbEpisode.Add()
    if nil != err {
        return err
    }

    return DeleteNyaaCached(title, dbEpisode.Episode)
}

func DelEpisode(route string, hash string) error {
    anime := dbase.Anime{}
    err := anime.Select(route)
    if nil != err {
        return err
    }

    episodes := dlSelectAll(anime)

    for _, e := range(episodes) {
        if e.Hash == hash {
            return e.Delete()
        }
    }

    return errors.New("Can't find episode...")
}

const WEEK_HOUR_DIFF = 24 * 7
var lastTimetableCheck = time.Now()
var lastTimetable = animeschedule.Timetable{}

func getCurrentEpisode(anime *dbase.Anime) {
    if "Finished" == anime.Status {
        anime.EpisodeCurrent = anime.EpisodeCount
        anime.EpisodeRelease = anime.SubTime.Add(time.Hour * time.Duration(WEEK_HOUR_DIFF * anime.EpisodeCount))
        return
    }

    if time.Now().Sub(lastTimetableCheck).Minutes() > 1 {
        err := sJsonHttpUnmarshall(&lastTimetable, "/timetables/sub")
        if nil != err {
            return
        }
        lastTimetableCheck = time.Now()
    }

    timetable := lastTimetable

    for _, tanime := range(timetable) {
        if anime.Route == tanime.Route {
            anime.EpisodeCurrent = tanime.EpisodeNumber
            anime.EpisodeRelease = tanime.EpisodeDate
            return
        }
    }
}
