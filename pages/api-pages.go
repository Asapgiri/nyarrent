package pages

import (
	"io"
	"net/http"
	"nyarrent/logic"
)

func ApiGetTorrentJson(w http.ResponseWriter, r *http.Request) {
    hash := r.PathValue("hash")

    io.WriteString(w, logic.ApiGetTorrentJson(hash))
}

func InvalidateCache(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("hash")
	cache.Invalidate(key)

    io.WriteString(w, "OK")
}
