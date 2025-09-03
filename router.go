package main

import (
	"net/http"
	"nyarrent/pages"
)

func setup_routes() {
    http.HandleFunc("GET /",                    pages.Root)
    http.HandleFunc("GET /index",               pages.Root)
    http.HandleFunc("GET /index.html",          pages.Root)

    http.HandleFunc("GET /downloads/{title}",           pages.Download)
    http.HandleFunc("GET /downloads/{title}/{sub}",     pages.Download)
    http.HandleFunc("GET /addtorrent",                  pages.AddTorrent)
    http.HandleFunc("GET /delete/{id}",                 pages.DeleteTorrent)

    http.HandleFunc("GET /listall",                     pages.ListAllTorrents)
    http.HandleFunc("GET /delallanime",                 pages.DeleteAllAnime)
    http.HandleFunc("GET /searchanime",                 pages.SearchNewAnimes)
    http.HandleFunc("GET /timetables",                  pages.ListTimetables)
    http.HandleFunc("GET /addanime/{route}",            pages.AddAnime)

    http.HandleFunc("GET /listanime/{route}",           pages.ListAnime)
    http.HandleFunc("GET /delanime/{route}",            pages.DeleteAnime)
    http.HandleFunc("GET /addepisode",                  pages.AddEpisode)
    http.HandleFunc("GET /delepisode",                  pages.DelEpisode)
    http.HandleFunc("GET /refreshnyaa/{route}/{index}", pages.RefreshNyaa)

    http.HandleFunc("GET /imgcache/{imgroute}",         pages.DownloadCachedImages)

    http.HandleFunc("GET /api/gettorrent/{hash}",   pages.ApiGetTorrentJson)
}
