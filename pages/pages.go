package pages

import (
	"io"
	"net/http"
	"nyarrent/dbase"
	"nyarrent/logger"
	"nyarrent/logic"
	"nyarrent/session"
	"os"
	"strconv"
	"strings"
	"time"
)

var log = logger.Logger {
    Color: logger.Colors.Red,
    Pretext: "pages",
}

var cache = dbase.Cache {
	Module: "pages",
	Time: time.Now(),
	Interval: 30 * time.Minute,
}

func ReturnPkiValidation(w http.ResponseWriter, r *http.Request) {
    file := r.PathValue("pkifile")

    // TODO: Create config path for it
    base_path := "/var/www/html/.well-known/pki-validation/"
    file_read, err := os.ReadFile(base_path + file)
    if nil != err {
	    http.Error(w, "Not found!", http.StatusNotFound)
    }

    io.WriteString(w, string(file_read))
}

func Root(w http.ResponseWriter, r *http.Request) {
    session := session.GetCurrentSession(r)
    if !checkForAccess(session, []string{logic.ROLES.USER}) && r.URL.Path != "/common.css" {
       http.Redirect(w, r, "/login", http.StatusSeeOther)
       return
    }

    if "/" == r.URL.Path {
        fu := r.URL.Query().Get("fu")
        dto := cache.GetWithInterval("root", func() any {
			return logic.ListAllAnime("true" == fu)
		}, time.Second).(logic.DtoAnime)
        dto.SearchText = r.URL.Query().Get("query")

        if "" != dto.SearchText {
            http.Redirect(w, r, "/searchanime?query="+dto.SearchText, http.StatusSeeOther)
            return
        }

        fil, _ := read_artifact("index.html", w.Header())
        Render(session, w, fil, dto)
    } else {
        Unexpected(session, w, r)
    }
}

func ListTimetables(w http.ResponseWriter, r *http.Request) {
    session := session.GetCurrentSession(r)
    if !checkForAccess(session, []string{logic.ROLES.USER}) {
       http.Redirect(w, r, "/login", http.StatusSeeOther)
       return
    }

    filter := getMap(r)

    filter.AnimeTimetable.OnlyOnList = r.URL.Query().Get("onlyonlist") == "on"
    filter.AnimeTimetable.SendBack = r.URL.Query().Has("sendback")

    refreshMap(&filter, r, "sendback", "onlyonlist")

    fil, _ := read_artifact("timetables.html", w.Header())

	logic_cached := cache.Get(filter.AnimeTimetable, func() any {
		return logic.ListTimetables(filter.AnimeTimetable)
	}).(logic.AnimeTimetablePage)

    Render(session, w, fil, logic_cached)
}

func SearchNewAnimes(w http.ResponseWriter, r *http.Request) {
    session := session.GetCurrentSession(r)
    if !checkForAccess(session, []string{logic.ROLES.USER}) {
       http.Redirect(w, r, "/login", http.StatusSeeOther)
       return
    }

    query := r.URL.Query().Get("query")
    page := r.URL.Query().Get("page")

    fil, _ := read_artifact("searchanime.html", w.Header())
    Render(session, w, fil, logic.FindNewAnimes(query, page))
}

func ListAllTorrents(w http.ResponseWriter, r *http.Request) {
    session := session.GetCurrentSession(r)
    if !checkForAccess(session, []string{logic.ROLES.USER}) {
       http.Redirect(w, r, "/login", http.StatusSeeOther)
       return
    }

    filter := getMap(r)

    page, err := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
    if nil == err {
        filter.Torrents.Page = int(page)
    } else {
        filter.Torrents.Page = defaultFilter.Torrents.Page
	}
    eps, err := strconv.ParseInt(r.URL.Query().Get("episodesperpage"), 10, 64)
    if nil == err {
        filter.Torrents.EpisodesPerPage = int(eps)
    } else {
		filter.Torrents.EpisodesPerPage = defaultFilter.Torrents.EpisodesPerPage
	}

	logic_cached := cache.Get(filter.Torrents, func() any {
		return logic.GetTorrents(&filter.Torrents)
	}).([]logic.Torrent)

    tl := logic.DtoBase {
        TorrentList: logic_cached,
        DiskUsage: logic.GetDiskUsage(),
		Filter: filter.Torrents,
    }

    fil, _ := read_artifact("listall.html", w.Header())
    Render(session, w, fil, tl)
}

func Unexpected(session session.Sessioner, w http.ResponseWriter, r *http.Request) {
    fil, typ := read_artifact(r.URL.Path, w.Header())

    if "text" == typ {
        Render(session, w, fil, nil)
    } else {
        io.WriteString(w, fil)
    }
}

// For torrents

func Download(w http.ResponseWriter, r *http.Request) {
    session := session.GetCurrentSession(r)
    if !checkForAccess(session, []string{logic.ROLES.USER}) {
       http.Redirect(w, r, "/login", http.StatusSeeOther)
       return
    }

    title := r.PathValue("title")
    sub := r.PathValue("sub")

    file, _, subs := logic.GetTorrentFile(title, false)

	if "" != sub {
		for _, s := range(subs) {
			if strings.Contains(s.Path, sub) {
				file = s.Path
				break
			}
		}
	}

    // log.Printf("path:  %s\n", r.URL.Path)
    // log.Printf("title: %s\n", title)
    // log.Printf("file:  %s\n", file)

    //w.Header().Add("Content-Type", "video/mp4")
    http.ServeFile(w, r, file)
}

func AddTorrent(w http.ResponseWriter, r *http.Request) {
    session := session.GetCurrentSession(r)
    if !checkForAccess(session, []string{logic.ROLES.ADMIN, logic.ROLES.EDITOR, logic.ROLES.MODERATOR}) {
       renderPageWithAccessViolation(w, r)
       return
    }

    link := r.URL.Query().Get("link")
    route := r.URL.Query().Get("route")

    logic.AddTorrent(link)

	filter := getMap(r)
	cache.Invalidate(filter.Episode.Hash+route)
	cache.Invalidate(filter.Torrents)

    http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

func DeleteTorrent(w http.ResponseWriter, r *http.Request) {
    session := session.GetCurrentSession(r)
    if !checkForAccess(session, []string{logic.ROLES.ADMIN, logic.ROLES.EDITOR, logic.ROLES.MODERATOR}) {
       renderPageWithAccessViolation(w, r)
       return
    }

    id := r.PathValue("id")
    route := r.URL.Query().Get("route")

    log.Printf("path:  %s\n", r.URL.Path)

    logic.DelTorrent(id)

	filter := getMap(r)
	cache.Invalidate(filter.Episode.Hash+route)
	cache.Invalidate(filter.Torrents)

    http.Redirect(w, r, "/", http.StatusSeeOther)
}

// For anime caching

func AddAnime(w http.ResponseWriter, r *http.Request) {
    session := session.GetCurrentSession(r)
    if !checkForAccess(session, []string{logic.ROLES.ADMIN, logic.ROLES.EDITOR, logic.ROLES.MODERATOR}) {
       renderPageWithAccessViolation(w, r)
       return
    }

    route := r.PathValue("route")
    sendback := r.URL.Query().Has("sendback")

    logic.AddorUpdateAnime(route)

    if sendback {
        log.Println(r.Header.Get("Referer"))
        http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
    } else {
        http.Redirect(w, r, "/listanime/"+route, http.StatusSeeOther)
    }
}

func ListAnime(w http.ResponseWriter, r *http.Request) {
    session := session.GetCurrentSession(r)
    if !checkForAccess(session, []string{logic.ROLES.USER}) {
       http.Redirect(w, r, "/login", http.StatusSeeOther)
       return
    }

    route := r.PathValue("route")

    filter := getMap(r)

    filter.Episode.Nyaa.Group = r.URL.Query().Get("group")
    filter.Episode.Nyaa.Category = r.URL.Query().Get("category")
    filter.Episode.Nyaa.SubCategory = r.URL.Query().Get("subcategory")
    filter.Episode.Nyaa.ResultCount = r.URL.Query().Get("resultcount")
    filter.Episode.Nyaa.NameParams = r.URL.Query().Get("nameparams")
    filter.Episode.Nyaa.Resolution = r.URL.Query().Get("resolution")
    filter.Episode.Nyaa.ForseRefrsh = r.URL.Query().Has("forcerefresh")

    refreshMap(&filter, r, "category", "subcategory", "resultcount", "nameparams", "forcerefresh")

	logic_cached := cache.Get(filter.Episode.Hash+route, func() any {
		return logic.ListAnime(route, &filter.Episode)
	}).(logic.Anime)

    fil, _ := read_artifact("listanime.html", w.Header())
    Render(session, w, fil, logic_cached)
    refreshMap(&filter, r)
}

func AddEpisode(w http.ResponseWriter, r *http.Request) {
    session := session.GetCurrentSession(r)
    if !checkForAccess(session, []string{logic.ROLES.ADMIN, logic.ROLES.EDITOR, logic.ROLES.MODERATOR}) {
       renderPageWithAccessViolation(w, r)
       return
    }

    route := r.URL.Query().Get("route")
    index := r.URL.Query().Get("index")
    link := r.URL.Query().Get("link")

    log.Println(link)

    title, hash, err := logic.AddTorrent(link)
    if nil == err {
        err = logic.AddEpisode(route, index, title, link, hash)
        if nil != err {
            log.Println(err.Error())
        }
    }

	filter := getMap(r)
	cache.Invalidate(filter.Episode.Hash+route)
	cache.Invalidate(filter.Torrents)

    http.Redirect(w, r, "/listanime/"+route, http.StatusSeeOther)
}

func DelEpisode(w http.ResponseWriter, r *http.Request) {
    session := session.GetCurrentSession(r)
    if !checkForAccess(session, []string{logic.ROLES.ADMIN, logic.ROLES.EDITOR, logic.ROLES.MODERATOR}) {
       renderPageWithAccessViolation(w, r)
       return
    }

    route := r.URL.Query().Get("route")
    hash := r.URL.Query().Get("hash")

    log.Println(hash)

    _, hash, err := logic.DelTorrent(hash)
    if nil == err {
        err = logic.DelEpisode(route, hash)
        if nil != err {
            log.Println(err.Error())
        }
    }

	filter := getMap(r)
	cache.Invalidate(filter.Episode.Hash+route)
	cache.Invalidate(filter.Torrents)

    http.Redirect(w, r, "/listanime/"+route, http.StatusSeeOther)
}

func RefreshNyaa(w http.ResponseWriter, r *http.Request) {
    session := session.GetCurrentSession(r)
    if !checkForAccess(session, []string{logic.ROLES.ADMIN, logic.ROLES.EDITOR, logic.ROLES.MODERATOR}) {
       renderPageWithAccessViolation(w, r)
       return
    }

    route := r.PathValue("route")
    index := r.PathValue("index")

    filter := getMap(r)
    refreshMap(&filter, r)

    logic.RefreshNyaa(route, index, filter.Episode.Nyaa)

	cache.Invalidate(getMap(r).Episode.Hash)

    http.Redirect(w, r, "/listanime/"+route, http.StatusSeeOther)
}

func DeleteAnime(w http.ResponseWriter, r *http.Request) {
    session := session.GetCurrentSession(r)
    if !checkForAccess(session, []string{logic.ROLES.ADMIN, logic.ROLES.EDITOR, logic.ROLES.MODERATOR}) {
       renderPageWithAccessViolation(w, r)
       return
    }

    route := r.PathValue("route")
    sendback := r.URL.Query().Has("sendback")

	logic.DeleteAnime(route)

    if sendback {
        http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
    } else {
    	http.Redirect(w, r, "/", http.StatusSeeOther)
    }
}

func DeleteAllAnime(w http.ResponseWriter, r *http.Request) {
    session := session.GetCurrentSession(r)
    if !checkForAccess(session, []string{logic.ROLES.ADMIN, logic.ROLES.EDITOR, logic.ROLES.MODERATOR}) {
       renderPageWithAccessViolation(w, r)
       return
    }

    sendback := r.URL.Query().Has("sendback")

	logic.DeleteAllAnime()

    if sendback {
        http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
    } else {
    	http.Redirect(w, r, "/", http.StatusSeeOther)
    }
}

func NotFound(w http.ResponseWriter, r *http.Request) {
    session := session.GetCurrentSession(r)

    fil, _ := read_artifact("not-found.html", w.Header())
    Render(session, w, fil, nil)
}

func AccessViolation(w http.ResponseWriter, r *http.Request) {
    session := session.GetCurrentSession(r)

    fil, _ := read_artifact("access-violation.html", w.Header())
    Render(session, w, fil, nil)
}

func renderPageWithAccessViolation(w http.ResponseWriter, r *http.Request) {
    AccessViolation(w, r)
}
