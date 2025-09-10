package main

import (
	"net/http"
	"nyarrent/config"
	"nyarrent/dbase"
	"nyarrent/logger"
	"os"
	"strings"
)

var msg = logger.Logger {
    Color: logger.Colors.Green,
    Pretext: "main",
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
