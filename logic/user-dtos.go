package logic

type Progress struct {
    Color       string
    Percentage  int
}

type Torrent struct {
    Id          string
    Title       string
    Url         string
    Size        string
    Eta         string
    Status      string
    Progress    Progress
    Info        string
}

type DtoBase struct {
    TorrentList []Torrent
}
