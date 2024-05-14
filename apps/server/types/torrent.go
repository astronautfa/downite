package types

import (
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/types"
)

var PiecePriorityStringMap = map[string]types.PiecePriority{
	"none":    types.PiecePriorityNone,
	"maximum": types.PiecePriorityNow,
	"high":    types.PiecePriorityHigh,
	"normal":  types.PiecePriorityNormal,
}

type TorrentStatus int

const (
	TorrentStatusPaused TorrentStatus = iota
	TorrentStatusDownloading
	TorrentStatusCompleted
	TorrentStatusSeeding
	TorrentStatusMetadata
)

var TorrentStatusStringMap = map[TorrentStatus]string{
	TorrentStatusPaused:      "paused",
	TorrentStatusDownloading: "downloading",
	TorrentStatusCompleted:   "completed",
	TorrentStatusSeeding:     "seeding",
	TorrentStatusMetadata:    "metadata",
}

type PieceProgress struct {
	Index               int
	DownloadedByteCount int
	Length              int
}
type TreeNodeMeta struct {
	Length   int64            `json:"length"`
	Name     string           `json:"name"`
	Path     []string         `json:"path"`
	Children *[]*TreeNodeMeta `json:"children"`
}
type TorrentMeta struct {
	TotalSize     int64           `json:"totalSize"`
	Files         []*TreeNodeMeta `json:"files"`
	Name          string          `json:"name"`
	InfoHash      string          `json:"infoHash"`
	TorrentMagnet string          `json:"torrentMagnet"`
}
type Torrent struct {
	Name          string                      `json:"name"`
	InfoHash      string                      `json:"infoHash"`
	Files         metainfo.FileTree           `json:"files"`
	TotalSize     int64                       `json:"totalSize" db:"total_size"`
	AmountLeft    int64                       `json:"amountLeft"`
	Uploaded      int64                       `json:"uploaded"`
	Downloaded    int64                       `json:"downloaded"`
	Magnet        string                      `json:"magnet"`
	Status        string                      `json:"status"`
	PieceProgress []PieceProgress             `json:"pieceProgress"`
	Peers         map[string]torrent.PeerInfo `json:"peers"`
	Progress      float32                     `json:"progress"`
	PeersCount    int                         `json:"peersCount"`
	Eta           int                         `json:"eta"`
	Category      string                      `json:"category"`
	SavePath      string                      `json:"savePath" db:"save_path"`
	Tags          []string                    `json:"tags"`
	AddedOn       int64                       `json:"addedOn" db:"added_on"`
	Availability  float32                     `json:"availability"`
	Ratio         float32                     `json:"ratio"`
	Seeds         int                         `json:"seeds"`
	DownloadSpeed int                         `json:"downloadSpeed"`
	UploadSpeed   int                         `json:"uploadSpeed"`
}

type TorrentFileOptions struct {
	Path             string `json:"path"`
	Name             string `json:"name"`
	DownloadPriority string `json:"downloadPriority" enum:"None,Low,Normal,High,Maximum"`
}

type DownloadPriority int

const (
	DownloadPriorityNone DownloadPriority = iota
	DownloadPriorityLow
	DownloadPriorityNormal
	DownloadPriorityHigh
	DownloadPriorityMaximum
)
