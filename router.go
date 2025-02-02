package main

import (
	"net/http"
	"nyarrent/pages"
)

func setup_routes() {
    http.HandleFunc("GET /",                    pages.Root)
    http.HandleFunc("GET /index",               pages.Root)
    http.HandleFunc("GET /index.html",          pages.Root)

    http.HandleFunc("GET /downloads/{title}",   pages.Download)
    http.HandleFunc("POST /addtorrent",         pages.AddTorrent)
    http.HandleFunc("GET /delete/{id}",         pages.DeleteTorrent)

    http.HandleFunc("GET /listall",             pages.ListAllTorrents)
    http.HandleFunc("GET /searchanime",         pages.SearchNewAnimes)
    http.HandleFunc("GET /addanime/{route}",    pages.AddAnime)

    http.HandleFunc("GET /listanime/{route}",   pages.ListAnime)
    http.HandleFunc("POST /addepisode",         pages.AddEpisode)
}
