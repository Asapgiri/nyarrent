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

type DiskUsage struct {
    Size    string
    Used    string
    Avail   string
    Percent string
    Usage   Progress
}

type DtoBase struct {
    TorrentList []Torrent
    DiskUsage   DiskUsage
}

// Transfer related stuff

type TorrentInfo struct {
    ActivityDate        int         `json:activityDate`
    AddedDate           int         `json:addedDate`
    BandwidthPriority   int         `json:bandwidthPriority`
    Comment             string      `json:comment`
    CorruptEver         int         `json:corruptEver`
    Creator             string      `json:creator`
    DateCreated         int         `json:dateCreated`
    DesiredAvailable    int         `json:desiredAvailable`
    DoneDate            int         `json:doneDate`
    DownloadDir         string      `json:downloadDir`
    DownloadLimit       int         `json:downloadLimit`
    DownloadLimited     bool        `json:downloadLimited`
    DownloadedEver      int         `json:downloadedEver`
    Error               int         `json:error`
    ErrorString         string      `json:errorString`
    Eta                 int         `json:eta`
    Group               string      `json:group`
    HashString          string      `json:hashString`
    HaveUnchecked       int         `json:haveUnchecked`
    HaveValid           int         `json:haveValid`
    HonorsSessionLimits bool        `json:honorsSessionLimits`
    Id                  int         `json:id`
    IsFinished          bool        `json:isFinished`
    IsPrivate           bool        `json:isPrivate`
    Labels              []string    `json:labels`
    LeftUntilDone       int         `json:leftUntilDone`
    MagnetLink          string      `json:magnetLink`
    Name                string      `json:name`
    PeerLimit           int         `json:peer-limit`
    PeersConnected      int         `json:peersConnected`
    PeersGettingFromUs  int         `json:peersGettingFromUs`
    PeersSendingToUs    int         `json:peersSendingToUs`
    PieceCount          int         `json:pieceCount`
    PieceSize           int         `json:pieceSize`
    RateDownload        int         `json:rateDownload`
    RateUpload          int         `json:rateUpload`
    RecheckProgress     int         `json:recheckProgress`
    SecondsDownloading  int         `json:secondsDownloading`
    SecondsSeeding      int         `json:secondsSeeding`
    SeedRatioLimit      int         `json:seedRatioLimit`
    SeedRatioMode       int         `json:seedRatioMode`
    SizeWhenDone        int         `json:sizeWhenDone`
    Source              string      `json:source`
    StartDate           int         `json:startDate`
    Status              int         `json:status`
    TotalSize           int         `json:totalSize`
    UploadLimit         int         `json:uploadLimit`
    UploadLimited       bool        `json:uploadLimited`
    UploadRatio         float64     `json:uploadRatio`
    UploadedEver        int         `json:uploadedEver`
    Webseeds            []string    `json:webseeds`
    WebseedsSendingToUs int         `json:webseedsSendingToUs`
}

type TrArguments struct {
    Torrents    []TorrentInfo       `json:torrents`
}

type TransmissionResults struct {
    Arguments   TrArguments         `json:arguments`
    Result      string              `json:results`
    Tag         int                 `json:tag`
}
