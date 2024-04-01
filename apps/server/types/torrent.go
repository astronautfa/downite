package types

import (
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
)

type TorrentStatus int

const (
	TorrentStatusPaused TorrentStatus = iota
	TorrentStatusDownloading
	TorrentStatusCompleted
	TorrentStatusSeeding
)

type PieceProgress struct {
	Index               int
	DownloadedByteCount int
	Length              int
}
type FileMeta struct {
	Length int64
	Path   []string
}
type Torrent struct {
	Name          string                      `json:"name"`
	InfoHash      string                      `json:"infoHash"`
	Files         metainfo.FileTree           `json:"files"`
	TotalSize     int                         `json:"totalSize"`
	Status        TorrentStatus               `json:"status"`
	PieceProgress []PieceProgress             `json:"pieceProgress"`
	Peers         map[string]torrent.PeerInfo `json:"peers"`
	Progress      float32                     `json:"progress"`
	PeersCount    int                         `json:"peersCount"`
	Eta           int                         `json:"eta"`
	Category      string                      `json:"category"`
	DownloadPath  string                      `json:"downloadPath"`
	DownloadDir   string                      `json:"downloadDir"`
	Tags          []string                    `json:"tags"`
	AddedOn       int64                       `json:"addedOn"`
	Availability  float32                     `json:"availability"`
	Ratio         float32                     `json:"ratio"`
	Seeds         int                         `json:"seeds"`
	DownloadSpeed int                         `json:"downloadSpeed"`
	UploadSpeed   int                         `json:"uploadSpeed"`
}