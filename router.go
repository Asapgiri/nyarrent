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
}
