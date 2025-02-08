package main

import (
	"net/http"
	"nyarrent/config"
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
    config.InitConfig()

    dbase.Connect()
    setup_routes()

    args := os.Args[1:]
    if 0 < len(args) {
        config.Config.Http.Port = args[0];
    }

    err := http.ListenAndServe(strings.Join([]string{":", config.Config.Http.Port}, ""), nil)
    msg.Println(err)
}
