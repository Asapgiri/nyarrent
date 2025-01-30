package main

import (
	"net/http"
	"nyarrent/logger"
	"nyarrent/pages"
)

var msg = logger.Logger {
    Color: logger.Colors.Green,
    Pretext: "main",
}

type Hello struct{}


func (h Hello) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    var path = r.URL.Path
    msg.Printf("Serving request: %#v\n", path)

    pages.Unexpected(w, r)
}

func main() {
    setup_routes()

    err := http.ListenAndServe(":80", nil)
    msg.Println(err)
}
