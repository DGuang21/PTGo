package client

import (
	"errors"
)

var (
	Box           = newClientDistributeService()
	distributeMap = map[string]clientDistributeerService{
		"Deluge":        delugeClient,
		"QBittorrent":   qBittorrentClient,
		"QBittorrent42": qBittorrent42Client,
		"Rutorrent":     rutorrentClient,
		"Transmission":  transmissionClient,
		"Utorrent":      utorrentClient,
	}
)

func (t *clientDistributeService) Login() error {
	panic("System Error")
	return nil
}

func (t *clientDistributeService) GetTorrentsList() (TorrentList, error) {
	panic("System Error")
	return nil, nil
}

func (t *clientDistributeService) AddTorrent(url string) error {
	panic("System Error")
	return nil
}

func (t *clientDistributeService) DeleteTorrents(hash string) error {
	panic("System Error")
	return nil
}

func (t *clientDistributeService) PauseTorrents(hash string) error {
	panic("System Error")
	return nil
}

func (t *clientDistributeService) ResumeTorrents(hash string) error {
	panic("System Error")
	return nil
}

func (t *clientDistributeService) ForceStart(hash string) error {
	panic("System Error")
	return nil
}

func (t clientDistributeService) SetConfig(config *BoxConfig) error {
	if config == nil {
		return errors.New("no config")
	}
	realization, ok := distributeMap[config.BoxName]
	if !ok {
		return errors.New("Unkonw Client Name")
	}
	if err := realization.SetConfig(config); err != nil {
		return err
	}
	Box = realization
	return nil
}

type clientDistributeService struct {
	clientDistributeerService
}

type BoxConfig struct {
	BoxName     string
	BoxUrl      string
	BoxUserName string
	BoxPassWord string
	BoxCookie   string
}

var _ clientDistributeerService = (*clientDistributeService)(nil)

func newClientDistributeService() clientDistributeerService {
	return &clientDistributeService{}
}
