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
        dto := logic.ListAnimes()
        dto.SearchText = r.URL.Query().Get("query")
        log.Println(dto.Anime[0])
        fil, _ := read_artifact("index.html", w.Header())
        Render(w, fil, dto)
    } else {
        Unexpected(w, r)
    }
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
    link := r.FormValue("link")

    logic.AddTorrent(link)

    http.Redirect(w, r, "/", http.StatusSeeOther)
}

func DeleteTorrent(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")

    log.Printf("path:  %s\n", r.URL.Path)

    logic.DeleteTorrent(id)

    http.Redirect(w, r, "/", http.StatusSeeOther)
}

// For anime caching

func AddAnime(w http.ResponseWriter, r *http.Request) {
    route := r.PathValue("route")

    logic.AddAnime(route)

    http.Redirect(w, r, "/listanime/"+route, http.StatusSeeOther)
}

func ListAnime(w http.ResponseWriter, r *http.Request) {
    route := r.PathValue("route")
    anime := logic.ListAnime(route)

    log.Println(len(anime.Episodes[0].Torrents))
    log.Println(anime.Episodes[0].Torrents[0].Torrent)
    log.Println(anime.Episodes[0].Torrents[0].Info)

    fil, _ := read_artifact("listanime.html", w.Header())
    Render(w, fil, anime)
}

func AddEpisode(w http.ResponseWriter, r *http.Request) {
    route := r.FormValue("route")
    index := r.FormValue("index")
    link := r.FormValue("link")

    title, hash, err := logic.AddTorrent(link)
    if nil == err {
        logic.AddEpisode(route, index, title, link, hash)
    }

    http.Redirect(w, r, "/listanime/"+route, http.StatusSeeOther)
}
