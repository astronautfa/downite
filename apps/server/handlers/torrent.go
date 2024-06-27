package handlers

import (
	"context"
	"downite/db"
	"downite/download/protocol/torr"
	"downite/types"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"sort"
	"time"

	gotorrent "github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/types/infohash"
	"github.com/danielgtaylor/huma/v2"
)

type TorrentHandler struct {
	Db     *db.Database
	Engine *torr.TorrentEngine
}

type TorrentsTotalSpeedData struct {
	DownloadSpeed float32 `json:"downloadSpeed"`
	UploadSpeed   float32 `json:"uploadSpeed"`
	Time          string  `json:"time"`
}
type GetTorrentsTotalSpeedRes struct {
	Body TorrentsTotalSpeedData
}

func (handler *TorrentHandler) GetTorrentsTotalSpeed(ctx context.Context, input *struct{}) (*GetTorrentsTotalSpeedRes, error) {
	res := &GetTorrentsTotalSpeedRes{}
	res.Body.DownloadSpeed = handler.Engine.GetTotalDownloadSpeed()
	res.Body.UploadSpeed = handler.Engine.GetTotalUploadSpeed()
	res.Body.Time = time.Now().Format("15:04:05")
	return res, nil
}

type GetTorrentsRes struct {
	Body struct {
		Torrents []types.Torrent `json:"torrents"`
	}
}

func (handler *TorrentHandler) GetTorrents(ctx context.Context, input *struct{}) (*GetTorrentsRes, error) {
	res := &GetTorrentsRes{}
	torrentsRes := []types.Torrent{}

	torrents := handler.Engine.Client.Torrents()
	for _, torrent := range torrents {
		dbTorrent, err := handler.Engine.GetTorrentDetails(torrent.InfoHash())
		if err != nil {
			return nil, err
		}
		torrentsRes = append(torrentsRes, *dbTorrent)
	}

	sort.Slice(torrentsRes, func(i, j int) bool {
		return torrentsRes[i].QueueNumber < torrentsRes[j].QueueNumber
	})
	res.Body.Torrents = torrentsRes
	return res, nil
}

type GetTorrentReq struct {
	Infohash string `path:"infohash" maxLength:"40" example:"2b66980093bc11806fab50cb3cb41835b95a0362" doc:"Infohash of the torrent"`
}
type GetTorrentRes struct {
	Body types.Torrent
}

func (handler *TorrentHandler) GetTorrent(ctx context.Context, input *GetTorrentReq) (*GetTorrentRes, error) {
	res := &GetTorrentRes{}
	torrent, ok := handler.Engine.Client.Torrent(infohash.FromHexString(input.Infohash))
	if !ok {
		return nil, fmt.Errorf("torrent with hash %s not found", input.Infohash)
	}
	dbTorrent, err := handler.Engine.GetTorrentDetails(torrent.InfoHash())
	if err != nil {
		return nil, err
	}

	res.Body = *dbTorrent

	return res, nil
}

type TorrentActionReq struct {
	Body struct {
		InfoHashes []string `json:"infoHashes" maxLength:"30" example:"2b66980093bc11806fab50cb3cb41835b95a0362" doc:"Hashes of torrents"`
	}
}
type TorrentActionRes struct {
	Body struct {
		Success bool `json:"result"`
	}
}

func (handler *TorrentHandler) PauseTorrent(ctx context.Context, input *TorrentActionReq) (*TorrentActionRes, error) {
	res := &TorrentActionRes{}
	foundTorrents, err := handler.Engine.FindTorrents(input.Body.InfoHashes)
	if err != nil {
		return nil, err
	}
	for _, foundTorrent := range foundTorrents {
		if foundTorrent.Info() != nil {
			foundTorrent.CancelPieces(0, foundTorrent.NumPieces())
			foundTorrent.SetMaxEstablishedConns(0)
			torrent, err := handler.Db.GetTorrent(foundTorrent.InfoHash().String())
			if err != nil {
				return nil, err
			}
			torrent.Status = types.TorrentStatusStringMap[types.TorrentStatusPaused]
			handler.Db.UpdateTorrentStatus(torrent.Infohash, types.TorrentStatusPaused)

		} else {
			return nil, fmt.Errorf("cannot modify torrent because metainfo is not yet received")
		}
	}
	res.Body.Success = true

	return res, nil
}
func (handler *TorrentHandler) ResumeTorrent(ctx context.Context, input *TorrentActionReq) (*TorrentActionRes, error) {
	res := &TorrentActionRes{}
	foundTorrents, err := handler.Engine.FindTorrents(input.Body.InfoHashes)
	if err != nil {
		return nil, err
	}
	for _, foundTorrent := range foundTorrents {
		// TODO(fatih): check if torrent is already started
		if foundTorrent.Info() != nil {
			foundTorrent.SetMaxEstablishedConns(80)
			handler.Engine.StartTorrent(foundTorrent)

			torrent, err := handler.Db.GetTorrent(foundTorrent.InfoHash().String())
			if err != nil {
				return nil, err
			}
			torrent.Status = types.TorrentStatusStringMap[types.TorrentStatusDownloading]
			handler.Db.UpdateTorrentStatus(torrent.Infohash, types.TorrentStatusDownloading)

		} else {
			return nil, fmt.Errorf("cannot modify torrent because metainfo is not yet received")
		}
	}
	res.Body.Success = true

	return res, nil
}
func (handler *TorrentHandler) RemoveTorrent(ctx context.Context, input *TorrentActionReq) (*TorrentActionRes, error) {
	res := &TorrentActionRes{}
	foundTorrents, err := handler.Engine.FindTorrents(input.Body.InfoHashes)
	if err != nil {
		return nil, err
	}
	for _, foundTorrent := range foundTorrents {
		// TODO(fatih): check if torrent is already started
		if foundTorrent.Info() != nil {
			foundTorrent.SetMaxEstablishedConns(0)
			foundTorrent.CancelPieces(0, foundTorrent.NumPieces())
			foundTorrent.Drop()
			handler.Db.DeleteTorrent(foundTorrent.InfoHash().String())
		} else {
			return nil, fmt.Errorf("cannot modify torrent because metainfo is not yet received")
		}
	}
	res.Body.Success = true

	return res, nil
}

// this is also deletes the torrent from disk
func (handler *TorrentHandler) DeleteTorrent(ctx context.Context, input *TorrentActionReq) (*TorrentActionRes, error) {
	res := &TorrentActionRes{}
	foundTorrents, err := handler.Engine.FindTorrents(input.Body.InfoHashes)
	if err != nil {
		return nil, err
	}
	for _, foundTorrent := range foundTorrents {
		// TODO(fatih): check if torrent is already started
		if foundTorrent.Info() == nil {
			return nil, fmt.Errorf("cannot modify torrent because metainfo is not yet received")
		}
		foundTorrent.SetMaxEstablishedConns(0)
		foundTorrent.CancelPieces(0, foundTorrent.NumPieces())
		foundTorrent.Drop()

		dbTorrent, err := handler.Db.GetTorrent(foundTorrent.InfoHash().String())
		err = handler.Db.DeleteTorrent(foundTorrent.InfoHash().String())
		err = handler.Db.DeleteTorrentFilesByInfohash(foundTorrent.InfoHash().String())
		err = os.RemoveAll(filepath.Join(dbTorrent.SavePath, foundTorrent.Name()))
		if err != nil {
			return nil, err
		}
	}
	res.Body.Success = true
	return res, nil
}

type DownloadTorrentData struct {
	TorrentFile multipart.File `form-data:"torrentFile" content-type:"application/x-bittorrent" required:"false"`
}
type DownloadTorrentReqBody struct {
	Magnet                      string                          `json:"magnet"`
	SavePath                    string                          `json:"savePath" validate:"required, dir"`
	IsIncompleteSavePathEnabled bool                            `json:"isIncompleteSavePathEnabled"`
	IncompleteSavePath          string                          `json:"incompleteSavePath,omitempty" validate:"dir"`
	Category                    string                          `json:"category,omitempty"`
	Tags                        []string                        `json:"tags,omitempty"`
	StartTorrent                bool                            `json:"startTorrent"`
	AddTopOfQueue               bool                            `json:"addTopOfQueue"`
	DownloadSequentially        bool                            `json:"downloadSequentially"`
	SkipHashCheck               bool                            `json:"skipHashCheck"`
	ContentLayout               string                          `json:"contentLayout" enum:"Original,Create subfolder,Don't create subfolder"`
	Files                       []types.TorrentFileFlatTreeNode `json:"files"`
	TorrentFile                 multipart.File                  `json:"torrentFile" required:"false"`
}
type DownloadTorrentReq struct {
	RawBody huma.MultipartFormFiles[DownloadTorrentData]
	Body    DownloadTorrentReqBody
}
type DownloadTorrentRes struct {
	Body types.Torrent
}

func (input *DownloadTorrentReq) Resolve(ctx huma.Context, prefix *huma.PathBuffer) []error {
	form := input.RawBody.Form
	requiredFields := []string{
		"savePath",
		"isIncompleteSavePathEnabled",
		"startTorrent",
		"addTopOfQueue",
		"downloadSequentially",
		"skipHashCheck",
		"contentLayout",
		"files",
	}
	var errors []error
	if form.File["torrentFile"] == nil && form.Value["magnet"] == nil {
		errors = append(errors, &huma.ErrorDetail{
			Location: prefix.String(),
			Message:  "either torrentFile or magnet is required",
			Value:    input,
		})
	}
	for _, requiredField := range requiredFields {
		if form.Value[requiredField] == nil {
			errors = append(errors, &huma.ErrorDetail{
				Location: prefix.String(),
				Message:  fmt.Sprintf("%s is required", requiredField),
				Value:    input,
			})
		}
	}
	if form.Value["magnet"] != nil {
		// Validate magnet
		if _, err := metainfo.ParseMagnetUri(input.RawBody.Form.Value["magnet"][0]); err != nil {
			errors = append(errors, &huma.ErrorDetail{
				Location: prefix.String(),
				Message:  "invalid magnet",
				Value:    input,
			})
		}
	}
	return errors
}

func (handler *TorrentHandler) DownloadTorrent(ctx context.Context, input *DownloadTorrentReq) (*DownloadTorrentRes, error) {
	res := &DownloadTorrentRes{}

	var torrent *gotorrent.Torrent
	var torrentSpec *gotorrent.TorrentSpec
	var err error

	if input.RawBody.Form.File["torrentFile"] != nil {
		// fileData := input.RawBody.Data()
		torrentFile, err := input.RawBody.Form.File["torrentFile"][0].Open()
		if err != nil {
			return nil, err
		}
		// Load the torrent file
		torrentMeta, err := metainfo.Load(torrentFile)
		if err != nil {
			return nil, err
		}
		torrentSpec = gotorrent.TorrentSpecFromMetaInfo(torrentMeta)
	} else {
		// Load from a magnet link
		torrentSpec, err = gotorrent.TorrentSpecFromMagnetUri(input.RawBody.Form.Value["magnet"][0])
		if err != nil {
			return nil, err
		}
	}

	// Register Torrent To DB
	dbTorrent, err := handler.Engine.RegisterTorrent(torrentSpec.InfoHash.String(), torrentSpec.DisplayName, input.RawBody.Form.Value["savePath"][0], torrentSpec.Trackers)
	if err != nil {
		return nil, err
	}
	// ADD TORRENT TO CLIENT
	torrent, err = handler.Engine.AddTorrent(dbTorrent.Infohash, dbTorrent.Trackers, dbTorrent.SavePath, input.RawBody.Form.Value["skipHashCheck"][0] != "true")
	if err != nil {
		return nil, err
	}

	// Convert form data to flat file tree
	flatFileTree := []types.TorrentFileFlatTreeNode{}
	err = json.Unmarshal([]byte(input.RawBody.Form.Value["files"][0]), &flatFileTree)
	if err != nil {
		return nil, err
	}
	// Register torrent files
	handler.Engine.RegisterFiles(torrent.InfoHash(), &flatFileTree)

	if input.RawBody.Form.Value["startTorrent"][0] == "true" {
		torrent, err = handler.Engine.StartTorrent(torrent)
		if err != nil {
			return nil, err
		}
		dbTorrent.Status = types.TorrentStatusStringMap[types.TorrentStatusDownloading]
	} else {
		dbTorrent.Status = types.TorrentStatusStringMap[types.TorrentStatusPaused]
	}

	dbTorrent.TotalSize = torrent.Length()
	torrentMetaInfo := torrent.Metainfo()
	magnetLink, err := torrentMetaInfo.MagnetV2()
	if err != nil {
		return nil, err
	}
	dbTorrent.Magnet = magnetLink.String()

	handler.Db.UpdateTorrent(dbTorrent)

	res.Body = *dbTorrent

	return res, nil
}

type GetMetaWithFileData struct {
	TorrentFile multipart.File `form-data:"torrentFile" content-type:"application/x-bittorrent" required:"true"`
}
type GetMetaWithFileReq struct {
	RawBody huma.MultipartFormFiles[DownloadTorrentData]
}

type GetMetaWithFileRes struct {
	Body types.TorrentMeta
}

func (input *GetMetaWithFileReq) Resolve(ctx huma.Context, prefix *huma.PathBuffer) []error {
	torrentFiles := input.RawBody.Form.File["torrentFile"]

	// Form validation
	if len(torrentFiles) == 0 {
		return []error{
			&huma.ErrorDetail{
				Location: prefix.String(),
				Message:  "no torrent file provided",
				Value:    input,
			},
		}
	}
	if len(torrentFiles) > 1 {
		return []error{
			&huma.ErrorDetail{
				Location: prefix.String(),
				Message:  "only one torrent file can be provided",
				Value:    input,
			},
		}
	}
	return nil
}

func (handler *TorrentHandler) GetMetaWithFile(ctx context.Context, input *GetMetaWithFileReq) (*GetMetaWithFileRes, error) {
	// TODO(fatih): In the future, we should support multiple torrents
	res := &GetMetaWithFileRes{}
	torrentFiles := input.RawBody.Form.File["torrentFile"]

	var info metainfo.Info
	var infohash string
	var magnetLink string

	// Load the torrent file
	torrentFile, err := torrentFiles[0].Open()
	if err != nil {
		return nil, err
	}
	defer torrentFile.Close()

	torrentMeta, err := metainfo.Load(torrentFile)
	if err != nil {
		return nil, err
	}
	info, err = torrentMeta.UnmarshalInfo()
	if err != nil {
		return nil, err
	}
	infohash = torrentMeta.HashInfoBytes().String()
	magnet, err := torrentMeta.MagnetV2()
	if err != nil {
		return nil, err
	}
	magnetLink = magnet.String()
	if err != nil {
		return nil, err
	}

	fileTree := handler.Engine.CreateFileTreeFromMeta(info)

	res.Body = types.TorrentMeta{
		TotalSize: info.TotalLength(),
		Files:     fileTree,
		Name:      info.Name,
		Infohash:  infohash,
		Magnet:    magnetLink,
	}
	return res, nil
}

type GetMetaWithMagnetReq struct {
	Body struct {
		Magnet string `json:"magnet" minLength:"1"`
	}
}

type GetMetaWithMagnetRes struct {
	Body types.TorrentMeta
}

func (handler *TorrentHandler) GetMetaWithMagnet(ctx context.Context, input *GetMetaWithMagnetReq) (*GetMetaWithMagnetRes, error) {
	// TODO(fatih): In the future, we should support multiple torrents
	res := &GetMetaWithMagnetRes{}

	var info metainfo.Info
	var infohash string

	magnet := input.Body.Magnet
	if _, err := metainfo.ParseMagnetUri(magnet); err != nil {
		ctx.Err()
		return nil, huma.Error400BadRequest("invalid magnet")
	}
	// Load from a magnet link

	torrent, err := handler.Engine.Client.AddMagnet(magnet)
	if err != nil {
		return nil, err
	}

	<-torrent.GotInfo()

	info = *torrent.Info()
	infohash = torrent.InfoHash().String()

	fileTree := handler.Engine.CreateFileTreeFromMeta(info)

	res.Body = types.TorrentMeta{
		TotalSize: info.TotalLength(),
		Files:     fileTree,
		Name:      info.Name,
		Infohash:  infohash,
		Magnet:    magnet,
	}

	torrent.Drop()
	return res, nil
}
