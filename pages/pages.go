package pages

import (
	"io"
	"net/http"
	"nyarrent/logger"
	"nyarrent/logic"
)

var log = logger.Logger {
    Color: logger.Colors.Red,
    Pretext: "pages",
}

func Root(w http.ResponseWriter, r *http.Request) {
    if "/" == r.URL.Path {
        fu := r.URL.Query().Get("fu")
        dto := logic.ListAllAnime("true" == fu)
        dto.SearchText = r.URL.Query().Get("query")

        if "" != dto.SearchText {
            http.Redirect(w, r, "/searchanime?query="+dto.SearchText, http.StatusSeeOther)
            return
        }

        fil, _ := read_artifact("index.html", w.Header())
        Render(w, fil, dto)
    } else {
        Unexpected(w, r)
    }
}

func ListTimetables(w http.ResponseWriter, r *http.Request) {
    filter := getMap(r)

    filter.AnimeTimetable.OnlyOnList = r.URL.Query().Get("onlyonlist") == "on"
    filter.AnimeTimetable.SendBack = r.URL.Query().Has("sendback")

    refreshMap(&filter, r, "sendback", "onlyonlist")

    fil, _ := read_artifact("timetables.html", w.Header())
    Render(w, fil, logic.ListTimetables(filter.AnimeTimetable))
}

func SearchNewAnimes(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query().Get("query")
    page := r.URL.Query().Get("page")

    fil, _ := read_artifact("searchanime.html", w.Header())
    Render(w, fil, logic.FindNewAnimes(query, page))
}

func ListAllTorrents(w http.ResponseWriter, r *http.Request) {
    tl := logic.DtoBase {
        TorrentList: logic.GetTorrents(),
        DiskUsage: logic.GetDiskUsage(),
    }

    fil, _ := read_artifact("listall.html", w.Header())
    Render(w, fil, tl)
}

func Unexpected(w http.ResponseWriter, r *http.Request) {
    fil, typ := read_artifact(r.URL.Path, w.Header())

    if "text" == typ {
        Render(w, fil, nil)
    } else {
        io.WriteString(w, fil)
    }
}

// For torrents

func Download(w http.ResponseWriter, r *http.Request) {
    title := r.PathValue("title")
    file := logic.GetTorrentFile(title, false)

    log.Printf("path:  %s\n", r.URL.Path)
    log.Printf("title: %s\n", title)
    log.Printf("file:  %s\n", file)

    http.ServeFile(w, r, file)
}

func AddTorrent(w http.ResponseWriter, r *http.Request) {
    link := r.URL.Query().Get("link")

    logic.AddTorrent(link)

    http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

func DeleteTorrent(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")

    log.Printf("path:  %s\n", r.URL.Path)

    logic.DelTorrent(id)

    http.Redirect(w, r, "/", http.StatusSeeOther)
}

// For anime caching

func AddAnime(w http.ResponseWriter, r *http.Request) {
    route := r.PathValue("route")
    sendback := r.URL.Query().Has("sendback")

    logic.AddorUpdateAnime(route)

    if sendback {
        log.Println(r.Header.Get("Referer"))
        http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
    } else {
        http.Redirect(w, r, "/listanime/"+route, http.StatusSeeOther)
    }
}

func ListAnime(w http.ResponseWriter, r *http.Request) {
    route := r.PathValue("route")

    filter := getMap(r)

    filter.Episode.Nyaa.Group = r.URL.Query().Get("group")
    filter.Episode.Nyaa.Category = r.URL.Query().Get("category")
    filter.Episode.Nyaa.SubCategory = r.URL.Query().Get("subcategory")
    filter.Episode.Nyaa.ResultCount = r.URL.Query().Get("resultcount")
    filter.Episode.Nyaa.NameParams = r.URL.Query().Get("nameparams")
    filter.Episode.Nyaa.Resolution = r.URL.Query().Get("resolution")
    filter.Episode.Nyaa.ForseRefrsh = r.URL.Query().Has("forcerefresh")

    refreshMap(&filter, r, "category", "subcategory", "resultcount", "nameparams", "forcerefresh")

    fil, _ := read_artifact("listanime.html", w.Header())
    Render(w, fil, logic.ListAnime(route, &filter.Episode))
    refreshMap(&filter, r)
}

func AddEpisode(w http.ResponseWriter, r *http.Request) {
    route := r.URL.Query().Get("route")
    index := r.URL.Query().Get("index")
    link := r.URL.Query().Get("link")

    log.Println(link)

    title, hash, err := logic.AddTorrent(link)
    if nil == err {
        err = logic.AddEpisode(route, index, title, link, hash)
        if nil != err {
            log.Println(err.Error())
        }
    }

    http.Redirect(w, r, "/listanime/"+route, http.StatusSeeOther)
}

func DelEpisode(w http.ResponseWriter, r *http.Request) {
    route := r.URL.Query().Get("route")
    hash := r.URL.Query().Get("hash")

    log.Println(hash)

    _, hash, err := logic.DelTorrent(hash)
    if nil == err {
        err = logic.DelEpisode(route, hash)
        if nil != err {
            log.Println(err.Error())
        }
    }

    http.Redirect(w, r, "/listanime/"+route, http.StatusSeeOther)
}

func RefreshNyaa(w http.ResponseWriter, r *http.Request) {
    route := r.PathValue("route")
    index := r.PathValue("index")

    filter := getMap(r)
    refreshMap(&filter, r)

    logic.RefreshNyaa(route, index, filter.Episode.Nyaa)

    http.Redirect(w, r, "/listanime/"+route, http.StatusSeeOther)
}
