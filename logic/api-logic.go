package logic

func ApiGetTorrentJson(hash string) string {
    return getTorrentInfoString(hash, true)
}
