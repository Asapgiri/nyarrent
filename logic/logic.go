package logic

import (
	"encoding/json"
	"io"
	"io/fs"
	"net/http"
	"nyarrent/logger"
	"os"
	"os/exec"
	"regexp"
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
        file = "/" + strings.Replace(file, downloadFolder, downloadUrl, 1)
    }

    return file
}

func getTorrentInfoString(idHash string) string {
    cmd := exec.Command("transmission-remote", "-t", idHash, "-i")

    stdout, err := cmd.Output()
    if nil != err {
        log.Println(err.Error())
        return err.Error()
    }
    return string(stdout)
}

func getTorrentInfo(idHash string) TorrentInfo {
    infoStr := getTorrentInfoString(idHash)

    lines := strings.Split(infoStr, "\n")
    if len(lines) < 3 {
        return TorrentInfo{}
    }

    re := regexp.MustCompile("^.*: ")
    info := TorrentInfo{}

    info.Name.Id =                  re.ReplaceAllString(lines[1], "$2")
    info.Name.Name =                re.ReplaceAllString(lines[2], "$2")
    info.Name.Hash =                re.ReplaceAllString(lines[3], "$2")
    info.Name.Magnet =              re.ReplaceAllString(lines[4], "$2")
    info.Name.Labesl =              re.ReplaceAllString(lines[5], "$2")

    info.Transfer.State =           re.ReplaceAllString(lines[8], "$2")
    info.Transfer.Location =        re.ReplaceAllString(lines[9], "$2")
    info.Transfer.PercentDone =     re.ReplaceAllString(lines[10], "$2")
    info.Transfer.ETA =             re.ReplaceAllString(lines[11], "$2")
    info.Transfer.DownloadSpeed =   re.ReplaceAllString(lines[12], "$2")
    info.Transfer.UploadSpeed =     re.ReplaceAllString(lines[13], "$2")
    info.Transfer.Have =            re.ReplaceAllString(lines[14], "$2")
    info.Transfer.Availability =    re.ReplaceAllString(lines[15], "$2")
    info.Transfer.TotalSize =       re.ReplaceAllString(lines[16], "$2")
    info.Transfer.Downloaded =      re.ReplaceAllString(lines[17], "$2")
    info.Transfer.Uploaded =        re.ReplaceAllString(lines[18], "$2")
    info.Transfer.Ratio =           re.ReplaceAllString(lines[19], "$2")
    info.Transfer.CorruptDL =       re.ReplaceAllString(lines[20], "$2")
    info.Transfer.Peers =           re.ReplaceAllString(lines[21], "$2")

    info.Hystory.DateAdded =        strings.Join(strings.Fields(lines[24])[2:], " ")
    info.Hystory.DateFinished =     strings.Join(strings.Fields(lines[25])[2:], " ")
    info.Hystory.DateStarted =      strings.Join(strings.Fields(lines[26])[2:], " ")
    info.Hystory.LatestActivity =   strings.Join(strings.Fields(lines[27])[2:], " ")
    info.Hystory.DownloadingTime =  strings.Join(strings.Fields(lines[28])[2:], " ")
    info.Hystory.SeedingTime =      strings.Join(strings.Fields(lines[29])[2:], " ")

    info.Origins.DateCreated            = re.ReplaceAllString(lines[32], "$2")
    info.Origins.PublicTorrent          = re.ReplaceAllString(lines[33], "$2")
    info.Origins.Comment                = re.ReplaceAllString(lines[34], "$2")
    info.Origins.Creator                = re.ReplaceAllString(lines[35], "$2")
    info.Origins.PieceCount             = re.ReplaceAllString(lines[36], "$2")
    info.Origins.PieceSize              = re.ReplaceAllString(lines[37], "$2")

    info.LimitsB.DownloadLimit          = re.ReplaceAllString(lines[40], "$2")
    info.LimitsB.UploadLimit            = re.ReplaceAllString(lines[41], "$2")
    info.LimitsB.RatioLimit             = re.ReplaceAllString(lines[42], "$2")
    info.LimitsB.HonorsSessionLimits    = re.ReplaceAllString(lines[43], "$2")
    info.LimitsB.PeerLimit              = re.ReplaceAllString(lines[44], "$2")
    info.LimitsB.BandwidthPriority      = re.ReplaceAllString(lines[45], "$2")

    return info
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

            newTorrent.Info = getTorrentInfoString(newTorrent.Id)

            torrents = append(torrents, newTorrent)
        }
    }

    return torrents
}

func AddTorrent(link string) (string, string, error) {
    stdout, err := exec.Command("transmission-remote", "--list").Output()
    if nil != err {
        log.Println(err.Error())
        log.Println(string(stdout))
        return "", "", err
    }
    preLines := strings.Split(string(stdout), "\n")

    stdout, err = exec.Command("transmission-remote", "--add", link).Output()
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
    if len(preLines) + 1 != len(lines) {
        log.Println("Line calculated len mismatch...")
        return "", "", err
    }
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

func FindNewAnimes(query string, page string) AnimeSearchPage {
    var aniList = AnimeSearchPage{}
    q := strings.ReplaceAll(query, " ", "+")

    resp, err := http.Get("https://animeschedule.net/api/v3/anime?page="+page+"&q="+q)
    if nil != err {
        log.Println(err.Error())
        return AnimeSearchPage{}
    }
    defer resp.Body.Close()

    aniListJsonBarr, err := io.ReadAll(resp.Body)
    if nil != err {
        log.Println(err.Error())
        return AnimeSearchPage{}
    }

    err = json.Unmarshal(aniListJsonBarr, &aniList)
    if nil != err {
        log.Println(err.Error())
        return AnimeSearchPage{}
    }

    log.Println(aniList)

    return aniList
}

func ListAnimes() DtoAnime {
    return listAllAnime()
}
