package logic

import (
	"io/fs"
	"nyarrent/logger"
	"os"
	"os/exec"
	"slices"
	"strconv"
	"strings"
)

const downloadFolder = "/mnt/d/Download"
const downloadUrl = "downloads"
var zipMap = []string{}

var log = logger.Logger {
    Color: logger.Colors.Cyan,
    Pretext: "logic",
}

var color_list = []string{
    "dark",
    "secondary",
    "info",
    "success",
}

func GenerateProgress(percentage float64) Progress {
    var index int

    if percentage >= 1.0 {
        index = len(color_list) - 1
    } else {
        index = int(float64(len(color_list)) * percentage)
    }


    return Progress{
        Color: color_list[index],
        Percentage: int(percentage * 100),
    }
}

func removeElement(slice []string, index int) []string {
    for i := index; i < len(slice)-1; i++ {
        slice[i] = slice[i+1]
    }
    return slice[:len(slice)-1]
}

func zipIt(path string) string {
    zipName := strings.Join([]string{path, "zip"}, ".")
    _, err := getFileInfo(zipName)
    if nil != err {
        if slices.Contains(zipMap, zipName) {
            return "In progress..."
        }

        zipMap = append(zipMap, zipName)
        out, err := exec.Command("zip", "-r", "-D", zipName, path).Output()

        log.Println(err)
        log.Println(string(out))
        log.Println("finished!")

        index := slices.Index(zipMap, zipName)
        zipMap = removeElement(zipMap, index)
    }

    return zipName
}

func getFileInfo(path string) (fs.FileInfo, error) {
    file, err := os.Open(path)
    if nil != err {
        log.Println(err.Error())
        return nil, err
    }

    fileInfo, err := file.Stat()
    if err != nil {
        log.Println(err.Error())
        return nil, err
    }

    return fileInfo, nil
}

func GetTorrentFile(title string, publicPath bool) string {
    path := strings.Join([]string{downloadFolder, title}, "/")

    log.Println(path)
    fileInfo, err := getFileInfo(path)
    if nil != err {
        log.Println(err.Error())
        return err.Error()
    }

    var file string
    if fileInfo.IsDir() {
        file = zipIt(path)
    } else {
        file = path
    }

    if publicPath {
        file = strings.Replace(file, downloadFolder, downloadUrl, 1)
    }

    return file
}

func GetTorrents() []Torrent {
    cmd := exec.Command("transmission-remote", "--list")

    stdout, err := cmd.Output()
    if nil != err {
        log.Println(err.Error())
        return []Torrent{}
    }

    lines := strings.Split(string(stdout), "\n")
    torrents := []Torrent{}

    for i, line := range lines {
        // Trim ID and trailing rows
        if 0 < i && len(lines) > i + 2 {
            fields := strings.Fields(line)

            percent, _ := strconv.ParseFloat(strings.Trim(fields[1], "%"), 64)

            newTorrent := Torrent{
                Id:         fields[0],
                Title:      strings.Join(fields[9:], " "),
                Size:       strings.Join(fields[2:4], " "),
                Eta:        fields[4],
                Status:     fields[8],
                Progress:   GenerateProgress(percent / 100),
            }
            newTorrent.Url = GetTorrentFile(newTorrent.Title, true)

            cmd := exec.Command("transmission-remote", "-t", newTorrent.Id, "-i")

            stdout, err := cmd.Output()
            if nil != err {
                log.Println(err.Error())
                return []Torrent{}
            }
            //newTorrent.Info = strings.Replace(string(stdout), "\n", "<br>", 0)
            newTorrent.Info = string(stdout)

            torrents = append(torrents, newTorrent)
        }
    }

    return torrents
}

func AddTorrent(link string) error {
    cmd := exec.Command("transmission-remote", "--add", link)

    stdout, err := cmd.Output()
    if nil != err {
        log.Println(err.Error())
        log.Println(string(stdout))
        return err
    }
    log.Println(string(stdout))

    return nil
}

func DeleteTorrent(id string) {
    cmd := exec.Command("transmission-remote", "-t", id, "-rad")

    stdout, err := cmd.Output()
    if nil != err {
        log.Println(err.Error())
        log.Println(string(stdout))
        return
    }
    log.Println(string(stdout))
}
