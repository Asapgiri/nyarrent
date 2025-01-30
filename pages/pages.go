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
        tl := logic.DtoBase {
            TorrentList: logic.GetTorrents(),
        }

        fil, _ := read_artifact("index.html", w.Header())
        Render(w, fil, tl)
    } else {
        Unexpected(w, r)
    }
}

func Unexpected(w http.ResponseWriter, r *http.Request) {
    fil, typ := read_artifact(r.URL.Path, w.Header())

    if "text" == typ {
        Render(w, fil, nil)
    } else {
        io.WriteString(w, fil)
    }
}

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
