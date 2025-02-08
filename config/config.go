package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type HttpConfig struct {
    Url     string
    Port    string
}

type DbConfig struct {
    Url     string
    Name    string
}

type AnimeScheduleConfig struct {
    ApiUrl  string
    ApiKey  string
}

type TorrentConfig struct {
    Disk        string
    Folder      string
}


type ConfigT struct {
    Http        HttpConfig
    Dbase       DbConfig
    AnimeAPI    AnimeScheduleConfig
    Downloads   TorrentConfig
}


var Config = ConfigT{
    Http: HttpConfig{
        Url:    "",
        Port:   "3000",
    },
    Dbase: DbConfig{
        Url:    "mongodb://localhost:27017",
        Name:   "nyarrent",
    },
    AnimeAPI: AnimeScheduleConfig{
        ApiUrl: "https://animeschedule.net/api/v3",
        ApiKey: "",
    },
    Downloads: TorrentConfig{
        Disk:   "/mnt/d",
        Folder: "downloads",
    },
}

func InitConfig() {
    ex, err := os.Executable()
    if nil != err {
        panic(err)
    }
    expath := filepath.Dir(ex)
    configfile := expath + "/.config.json"

    dat, err := os.ReadFile(configfile)
    if nil != err {
        log.Println(err.Error())
        configdat, _ := json.MarshalIndent(Config, "", "  ")
        os.WriteFile(configfile, configdat, 0644)
    } else {
        err = json.Unmarshal(dat, &Config)
        if nil != err {
            log.Println(err.Error())
        }
    }
}
