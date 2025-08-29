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
	"fmt"
	"path/filepath"

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

    log.Println("start zip-it ===============================================")
    log.Println(zipName)
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
    log.Println("end zip-it =================================================")

    return zipName
}

// FIXME: Really hacky hacky solution...
func getSubtitles(path string) []Subtitle {
	var files []string

	subsdir := strings.Join([]string{path, "subs"}, ".")
	_, err := os.Open(subsdir)

	if nil != err {
		// Dir not exists
		os.Mkdir(subsdir, 0755)

        out, _ := exec.Command("ffmpeg", "-i", path).CombinedOutput()
    	lines := strings.Split(string(out), "\n")

		for _, line := range(lines) {
			if strings.Contains(line, "Subtitle") {
				// Stream #0:2(eng): Subtitle: ass (default) (forced)
				lang := strings.Split(strings.Split(line, "(")[1], ")")[0]

				files = append(files, strings.Join([]string{lang, "vtt"}, "."))
			}
		}

		// Export out files
		for i, file := range(files) {
			which := strings.Join([]string{"0", "s", strconv.FormatInt(int64(i), 10)}, ":")
			fpath := strings.Join([]string{subsdir, file}, "/")
			exec.Command("ffmpeg", "-i", path, "-f", "webvtt", "-map", which, fpath).Output()
		}
	} else {
		// Dir exists
		entries, _ := os.ReadDir(subsdir)
		for _, e := range(entries) {
			files = append(files, e.Name())
		}
	}

	subsret := []Subtitle{}
	for _, p := range(files) {
		title := strings.Split(p, ".")[0]
		subsret = append(subsret, Subtitle{
			Title: title,
			Lang: title,
			Path: strings.Join([]string{subsdir, p}, "/"),
		})
	}

	return subsret
}

func convertMp4(path string) (string, []Subtitle) {
    mp4Name := strings.Join([]string{path, "mp4"}, ".")
    _, err := getFileInfo(mp4Name)
    if nil != err {
        if slices.Contains(zipMap, mp4Name) {
            return "In progress...", []Subtitle{}
        }

        zipMap = append(zipMap, mp4Name)
        out, err := exec.Command("ffmpeg", "-i", path, "-acodec", "copy", "-vcodec", "copy", mp4Name).Output()

        log.Println(err)
        log.Println(string(out))
        log.Println("finished!")

        index := slices.Index(zipMap, mp4Name)
        zipMap = removeElement(zipMap, index)
    }

    return mp4Name, getSubtitles(path)
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

func removeDoubleSlash(path string) string {
    retPath := path
    for strings.Contains(retPath, "//") {
        retPath = strings.Replace(retPath, "//", "/", -1)
    }
    return retPath
}

func GetTorrentFile(title string, publicPath bool) (string, string, []Subtitle) {
	var mp4 string
	var subs []Subtitle

    dlpath := removeDoubleSlash(strings.Join([]string{config.Config.Downloads.Disk, config.Config.Downloads.Folder}, "/"))
    path := removeDoubleSlash(strings.Join([]string{dlpath, title}, "/"))

    fileInfo, err := getFileInfo(path)
    if nil != err {
        log.Println(err.Error())
        return err.Error(), "", []Subtitle{}
    }

    var file string
    if fileInfo.IsDir() {
        file = zipIt(path)
    } else if ".mp4" != filepath.Ext(path) {
        log.Println(filepath.Ext(path))
		file = path
        mp4, subs = convertMp4(path)
    } else {
        file = path
    }

    if publicPath {
        file = "/" + strings.Replace(file, dlpath, "downloads", 1)
        mp4 = "/" + strings.Replace(mp4, dlpath, "downloads", 1)
		for i := 0; i < len(subs); i++ {
			subs[i].Path = "/" + strings.Replace(subs[i].Path, dlpath, "downloads", 1)
			subs[i].Path = strings.Replace(subs[i].Path, ".subs", "", 1)
		}
    }

    return file, mp4, subs
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

func GetTorrents(filter *TorrentsFilter) []Torrent {
    cmd := exec.Command("transmission-remote", "-j", "--list")

    stdout, err := cmd.Output()
    if nil != err {
        log.Println(err.Error())
        return []Torrent{}
    }

	list := TorrentList{}
	err = json.Unmarshal(stdout, &list)
	if nil != err {
		log.Println(err.Error())
	}
	slices.Reverse(list.Arguments.Torrents)

	pageCount := len(list.Arguments.Torrents) / filter.EpisodesPerPage
	if 0 != len(list.Arguments.Torrents) % filter.EpisodesPerPage {
		pageCount++
	}
	filter.PageCounter = make([]PageCounter, pageCount)
	for i := 0; i < pageCount; i++ {
		filter.PageCounter[i].Index = i + 1
		filter.PageCounter[i].Value = i
		filter.PageCounter[i].Selected = i == filter.Page
	}

    torrents := []Torrent{}

	start := filter.Page * filter.EpisodesPerPage
	end := start + filter.EpisodesPerPage

	if end > len(list.Arguments.Torrents) {
		end = len(list.Arguments.Torrents)
	}

    log.Println("start range ===============================================")
	for _, t := range list.Arguments.Torrents[start:end] {
        percent := float64(t.SizeWhenDone - t.LeftUntilDone) / float64(t.SizeWhenDone)

        newTorrent := Torrent{
            Id:         strconv.Itoa(t.Id),
            Title:      t.Name,
            Size:       strconv.Itoa(t.SizeWhenDone),
            Eta:        strconv.Itoa(t.Eta),
            Status:     strconv.Itoa(t.Status),
            Progress:   GenerateProgress(percent),
        }
        newTorrent.Url, _, _ = GetTorrentFile(newTorrent.Title, true)

        newTorrent.Info = getTorrentInfoString(newTorrent.Id, false)
        newTorrent.FullInfo = getTorrentInfo(newTorrent.Id)
        newTorrent.Dates = getTorrentDatesReadable(newTorrent.FullInfo)

        torrents = append(torrents, newTorrent)
    }
    log.Println("end range =================================================")

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

func calculate_len(length int) (string) {
	var calculated float32
	const splitter = 1024

	calculated = float32(length) / splitter
	if calculated < splitter {
		return fmt.Sprintf("%.2f KB", calculated)
	}

	calculated = calculated / splitter
	if calculated < splitter {
		return fmt.Sprintf("%.2f MB", calculated)
	}

	calculated = calculated / splitter
	return fmt.Sprintf("%.2f GB", calculated)
}

func GetNyaaList(title string, episode int, filter NyaaFilter) (string, []dbase.AnimeShotoTorrent) {
    var q string
    nameLen, err := strconv.ParseInt(filter.NameParams, 10, 64)
    if nil != err {
        nameLen = 2
    }
    if len(strings.Split(title, " ")) > int(nameLen) {
        q = strings.Join(strings.Split(title, " ")[:nameLen], "+")
    } else {
        q = strings.ReplaceAll(title, " ", "+")
    }

    episodeStr := strconv.FormatInt(int64(episode), 10)
    if len(episodeStr) <= 1 {
        episodeStr = "0" + episodeStr
    }

    q = strings.Join([]string{
        filter.Group,
        q,
        episodeStr,
        filter.Resolution,
    }, "+")
    shoto := []dbase.AnimeShotoTorrent{}

    if "" == filter.Category {
        filter.Category = "anime"
    }
    query := []string{
        "q=", q,
    }
    queryStr := strings.Join(query, "")

    nyaacache := dbase.NyaaCached{}
    nyaacache.Select(title, episode)
    if "" == nyaacache.Title || filter.ForseRefrsh {
        log.Println("Getting episode for: " + q)
        //resp, err := http.Get("https://nyaaapi.onrender.com/nyaa?"+queryStr)
        resp, err := http.Get("https://feed.animetosho.org/json?"+queryStr)
        if nil != err {
            log.Println(err.Error())
            return queryStr, []dbase.AnimeShotoTorrent{}
        }
        defer resp.Body.Close()

        aniListJsonBarr, err := io.ReadAll(resp.Body)
        if nil != err {
            log.Println(err.Error())
            return queryStr, []dbase.AnimeShotoTorrent{}
        }

        err = json.Unmarshal(aniListJsonBarr, &shoto)
        if nil != err {
            log.Println(err.Error())
            return queryStr, []dbase.AnimeShotoTorrent{}
        }

        for i := 0; i < len(shoto); i++ {
            shoto[i].SizeStr = calculate_len(shoto[i].Total_size)
        }

        nyaacache.Episode = episode
        nyaacache.Title = title
        //nyaacache.Nyaa = nyaaJson
        nyaacache.Shoto = shoto
        if nyaacache.Id.IsZero() {
            nyaacache.Id = primitive.NewObjectID()
            nyaacache.Add()
        } else {
            nyaacache.Update()
        }
    } else {
        //log.Println("Episodes exists for: " + q)
        //nyaaJson = nyaacache.Nyaa
        shoto = nyaacache.Shoto
    }

    resultCount, err := strconv.ParseInt(filter.ResultCount, 10, 64)
    if nil == err && len(shoto) >= int(resultCount) {
        return queryStr, shoto[:resultCount]
    } else {
        return queryStr, shoto
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
