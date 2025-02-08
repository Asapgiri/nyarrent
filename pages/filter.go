package pages

import (
	"net/http"
	"nyarrent/logic"
)

type Filters struct {
    AnimeTimetable  logic.AnimeTimetableFilter
    Episode         logic.EpisodeFilter
    Hash            string
}

var filterMap = map[string]Filters{}

func refreshHashes(filter *Filters) {
    filter.AnimeTimetable.Hash  = filter.Hash
    filter.Episode.Hash         = filter.Hash
}

var defaultFilter = Filters{
    Episode: logic.EpisodeFilter{
        Nyaa: logic.NyaaFilter{
            Group:          "",
            NameParams:     "3",
            Category:       "anime",
            SubCategory:    "eng",
            ResultCount:    "5",
            Resolution:     "1080p",
            ForseRefrsh:    false,
        },
    },
}

func getMap(r *http.Request) Filters {
    cookie, err := r.Cookie("formathash")
    if nil != err {
        return Filters{}
    }

    val, ok := filterMap[cookie.Value]
    if ok {
        return val
    } else {
        return defaultFilter
    }
}

func refreshMap(filter *Filters, r *http.Request, hasParams ...string) {
    filter.Hash = r.URL.Query().Get("hash")
    hasOne := r.URL.Query().Has("hash")

    cookie, err := r.Cookie("formathash")
    if nil != err {
        return
    }

    for _, param := range hasParams {
        hasOne = hasOne || r.URL.Query().Has(param)
    }

    val, ok := filterMap[cookie.Value]
    if !ok || hasOne {
        filter.Hash = cookie.Value
        refreshHashes(filter)
        filterMap[cookie.Value] = *filter
    } else {
        *filter = val
    }
}
