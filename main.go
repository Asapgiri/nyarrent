package main

import (
	"net/http"
	"nyarrent/dbase"
	"nyarrent/logger"
	"nyarrent/pages"
	"os"
	"strings"
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
    var port = "3000";

    dbase.Connect()
    setup_routes()

    args := os.Args[1:]
    if 0 < len(args) {
        port = args[0];
    }

    err := http.ListenAndServe(strings.Join([]string{":", port}, ""), nil)
    msg.Println(err)
}
