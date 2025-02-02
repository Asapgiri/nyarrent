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

// Transfer related stuff

type trName struct {
    Id      string
    Name    string
    Hash    string
    Magnet  string
    Labesl  string
}

type trTransfer struct {
    State           string //: Seeding
    Location        string //: /mnt/d/Download
    PercentDone     string //: 100%
    ETA             string //: 0 seconds (0 seconds)
    DownloadSpeed   string //: 0 kB/s
    UploadSpeed     string //: 13 kB/s
    Have            string //: 522.6 MB (522.6 MB verified)
    Availability    string //: 100%
    TotalSize       string //: 522.6 MB (522.6 MB wanted)
    Downloaded      string //: 527.7 MB
    Uploaded        string //: 19.49 MB
    Ratio           string //: 0.0
    CorruptDL       string //: 2.10 MB
    Peers           string //: connected to 1, uploading to 1, downloading from 0
}

type trHistory struct {
    DateAdded       string //:       Sun Feb  2 17:55:23 2025
    DateFinished    string //:    Sun Feb  2 17:57:16 2025
    DateStarted     string //:     Sun Feb  2 17:55:25 2025
    LatestActivity  string //:  Sun Feb  2 18:19:10 2025
    DownloadingTime string //: 2 minutes, 2 seconds (122 seconds)
    SeedingTime     string //:     23 minutes (1418 seconds)
}

type trOrigins struct {
    DateCreated     string //: Sun Feb  2 12:39:57 2025
    PublicTorrent   string //: Yes
    Comment         string //: https://nyaa.si/view/1929950
    Creator         string //: NyaaV2
    PieceCount      string //: 499
    PieceSize       string //: 1.00 MiB
}

type trLimitBandwidth struct {
    DownloadLimit       string //: Unlimited
    UploadLimit         string //: Unlimited
    RatioLimit          string //: Default
    HonorsSessionLimits string //: Yes
    PeerLimit           string //: 50
    BandwidthPriority   string //: Normal
}

type TorrentInfo struct {
    Name        trName
    Transfer    trTransfer
    Hystory     trHistory
    Origins     trOrigins
    LimitsB     trLimitBandwidth
}
