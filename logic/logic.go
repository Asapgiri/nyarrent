package logic

import (
	"encoding/json"
	"io"
	"io/fs"
	"math"
	"net/http"
	"nyarrent/config"
	"nyarrent/dbase"
	"nyarrent/logger"
	"os"
	"os/exec"
	"slices"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
    } else if percentage < 0 {
        index = 0
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

//func convertMp4(path string) string {
//    mp4Name := strings.Join([]string{path, "mp4"}, ".")
//    _, err := getFileInfo(mp4Name)
//    if nil != err {
//        if slices.Contains(zipMap, mp4Name) {
//            return "In progress..."
//        }
//
//        zipMap = append(zipMap, mp4Name)
//        out, err := exec.Command("ffmpeg", "-i", path, "-c:v", "copy", "-map", "0", "-c:s", "mov_text", mp4Name).Output()
//
//        log.Println(err)
//        log.Println(string(out))
//        log.Println("finished!")
//
//        index := slices.Index(zipMap, mp4Name)
//        zipMap = removeElement(zipMap, index)
//    }
//
//    return mp4Name
//}

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
    dlpath := strings.Join([]string{config.Config.Downloads.Disk, config.Config.Downloads.Folder}, "/")
    path := strings.Join([]string{dlpath, title}, "/")

    fileInfo, err := getFileInfo(path)
    if nil != err {
        log.Println(err.Error())
        return err.Error()
    }

    var file string
    if fileInfo.IsDir() {
        file = zipIt(path)
    //} else if ".mp4" != filepath.Ext(path) {
    //    log.Println(filepath.Ext(path))
    //    file = convertMp4(path)
    } else {
        file = path
    }

    if publicPath {
        file = "/" + strings.Replace(file, dlpath, "downloads", 1)
    }

    return file
}

func getTorrentInfoString(idHash string, isJson bool) string {
    var cmd *exec.Cmd
    if isJson {
        cmd = exec.Command("transmission-remote", "-t", idHash, "-j", "-i")
    } else {
        cmd = exec.Command("transmission-remote", "-t", idHash, "-i")
    }

    stdout, err := cmd.Output()
    if nil != err {
        log.Println(err.Error())
        return err.Error()
    }
    return string(stdout)
}

func getTorrentInfo(idHash string) TorrentInfo {
    infoStr := getTorrentInfoString(idHash, true)

    infoJson := TransmissionResults{}
    err := json.Unmarshal([]byte(infoStr), &infoJson)
    if nil != err {
        log.Println(err.Error())
        return TorrentInfo{}
    }

    if len(infoJson.Arguments.Torrents) > 0 {
        return infoJson.Arguments.Torrents[0]
    } else {
        return TorrentInfo{}
    }
}

var formatStr = "2006-01-02 15:04:05"
func timeAgo(date int64) string {
    t := time.Unix(date, 0)
    diff := time.Now().Sub(t)

    pastForm := ""

    if diff.Hours() <= 24 {
        hours := int64(math.Floor(diff.Hours()))
        minutes := int64(math.Floor(diff.Minutes())) % 60
        seconds := int64(math.Floor(diff.Seconds())) % 60
        if hours > 0 {
            pastForm += " (" + strconv.FormatInt(hours, 10) + "h "
            pastForm +=        strconv.FormatInt(minutes, 10) + "m ago)"
        } else if minutes > 0 {
            pastForm += " (" + strconv.FormatInt(minutes, 10) + "m ago)"
        } else {
            pastForm += " (" + strconv.FormatInt(seconds, 10) + "s ago)"
        }
    }


    return t.Format(formatStr) + pastForm
}

func getTorrentDatesReadable(info TorrentInfo) TorrentDatesReadable {
    return TorrentDatesReadable{
        Activity: timeAgo(int64(info.ActivityDate)),
        Added: timeAgo(int64(info.AddedDate)),
        Created: timeAgo(int64(info.DateCreated)),
        Done: timeAgo(int64(info.DoneDate)),
        Start: timeAgo(int64(info.StartDate)),
    }
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

            newTorrent.Info = getTorrentInfoString(newTorrent.Id, false)
            newTorrent.FullInfo = getTorrentInfo(newTorrent.Id)
            newTorrent.Dates = getTorrentDatesReadable(newTorrent.FullInfo)

            torrents = append(torrents, newTorrent)
        }
    }

    slices.Reverse(torrents)

    return torrents
}

func GetDiskUsage() DiskUsage {
    stdout, err := exec.Command("df", "-H", config.Config.Downloads.Disk).Output()
    if nil != err {
        log.Println(err.Error())
        return DiskUsage{}
    }

    line := strings.Split(string(stdout), "\n")[1]
    fields := strings.Fields(line)
    percent, _ := strconv.ParseFloat(strings.Trim(fields[4], "%"), 64)

    diskUsage := DiskUsage{
        Size:       fields[1],
        Used:       fields[2],
        Avail:      fields[3],
        Percent:    fields[4],
        Usage:      GenerateProgress(percent / 100),
    }

    return diskUsage
}

func AddTorrent(link string) (string, string, error) {
    stdout, err := exec.Command("transmission-remote", "--add", link).Output()
    if nil != err {
        log.Println(err.Error())
        log.Println(string(stdout))
        return "", "", err
    }

    stdout, err = exec.Command("transmission-remote", "--list").Output()
    if nil != err {
        log.Println(err.Error())
        log.Println(string(stdout))
        return "", "", err
    }

    lines := strings.Split(string(stdout), "\n")
    fields := strings.Fields(lines[len(lines) - 3])

    stdout, err = exec.Command("transmission-remote", "-t", fields[0], "-i").Output()
    if nil != err {
        log.Println(err.Error())
        log.Println(string(stdout))
        return "", "", err
    }
    lines = strings.Split(string(stdout), "\n")

    title := strings.Join(strings.Fields(lines[2])[1:], " ")
    hash := strings.Fields(lines[3])[1]

    return title, hash, nil
}

func DelTorrent(id string) (string, string, error) {
    stdout, err := exec.Command("transmission-remote", "-t", id, "-i").Output()
    if nil != err {
        log.Println(err.Error())
        log.Println(string(stdout))
        return "", "", err
    }
    lines := strings.Split(string(stdout), "\n")

    title := strings.Join(strings.Fields(lines[2])[1:], " ")
    hash := strings.Fields(lines[3])[1]

    stdout, err = exec.Command("transmission-remote", "-t", id, "-rad").Output()
    if nil != err {
        log.Println(err.Error())
        log.Println(string(stdout))
        return "", "", err
    }
    log.Println(string(stdout))

    return title, hash, nil
}

// Nyaa.si related

func GetNyaaList(title string, episode int, filter NyaaFilter) (string, []dbase.NyaaData) {
    var q string
    nameLen, err := strconv.ParseInt(filter.NameParams, 10, 64)
    if nil != err {
        nameLen = 2
    }
    if len(strings.Split(title, " ")) > int(nameLen) {
        q = strings.Join(strings.Split(title, " ")[:nameLen], "+")
    } else {
        q = strings.Replace(title, " ", "+", 0)
    }

    q = strings.Join([]string{
        filter.Group,
        q,
        strconv.FormatInt(int64(episode), 10),
        filter.Resolution,
    }, "+")
    nyaaJson := dbase.NyaaJson{}

    if "" == filter.Category {
        filter.Category = "anime"
    }
    query := []string{
        "q=",               q,
        "&category=",       filter.Category,
        "&sub_category=",   filter.SubCategory,
    }
    queryStr := strings.Join(query, "")

    nyaacache := dbase.NyaaCached{}
    nyaacache.Select(title, episode)
    if "" == nyaacache.Title || filter.ForseRefrsh {
        log.Println("Getting episode for: " + q)
        resp, err := http.Get("https://nyaaapi.onrender.com/nyaa?"+queryStr)
        if nil != err {
            log.Println(err.Error())
            return queryStr, []dbase.NyaaData{}
        }
        defer resp.Body.Close()

        aniListJsonBarr, err := io.ReadAll(resp.Body)
        if nil != err {
            log.Println(err.Error())
            return queryStr, []dbase.NyaaData{}
        }

        err = json.Unmarshal(aniListJsonBarr, &nyaaJson)
        if nil != err {
            log.Println(err.Error())
            return queryStr, []dbase.NyaaData{}
        }

        nyaacache.Episode = episode
        nyaacache.Title = title
        nyaacache.Nyaa = nyaaJson
        if nyaacache.Id.IsZero() {
            nyaacache.Id = primitive.NewObjectID()
            nyaacache.Add()
        } else {
            nyaacache.Update()
        }
    } else {
        //log.Println("Episodes exists for: " + q)
        nyaaJson = nyaacache.Nyaa
    }

    resultCount, err := strconv.ParseInt(filter.ResultCount, 10, 64)
    if nil == err && len(nyaaJson.Data) >= int(resultCount) {
        return queryStr, nyaaJson.Data[:resultCount]
    } else {
        return queryStr, nyaaJson.Data
    }

}

func DeleteNyaaCached(title string, episode int) error {
    nyaacache := dbase.NyaaCached{}
    err := nyaacache.Select(title, episode)
    if nil != err {
        return err
    }
    return nyaacache.Delete()
}

func RefreshNyaa(route string, episode string, filter NyaaFilter) {
    dbAnime := dbase.Anime{}
    dbAnime.Select(route)

    epiInt, err := strconv.ParseInt(episode, 10, 64)
    if nil != err {
        return
    }

    filter.ForseRefrsh = true
    GetNyaaList(dbAnime.Title, int(epiInt), filter)
}
