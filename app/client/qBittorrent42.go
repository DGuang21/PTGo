package client

var (
	qBittorrent42Client = newQBittorrent42Client()
)

type qBittorrent42 struct {
	BoxConfig
}

func (t *qBittorrent42) Login() error {
	panic("System Error")
	return nil
}

func (t *qBittorrent42) GetTorrentsInfo(hash string) (*TorrentItem, error) {
	panic("System Error")
	return nil, nil
}

func (t *qBittorrent42) GetTorrentsList() (TorrentList, error) {
	panic("System Error")
	return nil, nil
}

func (t *qBittorrent42) AddTorrent(url string) error {
	panic("System Error")
	return nil
}

func (t *qBittorrent42) DeleteTorrents(hash string) error {
	panic("System Error")
	return nil
}

func (t *qBittorrent42) PauseTorrents(hash string) error {
	panic("System Error")
	return nil
}

func (t *qBittorrent42) ResumeTorrents(hash string) error {
	panic("System Error")
	return nil
}

func (t *qBittorrent42) SetConfig(config *BoxConfig) error {
	panic("System Error")
	return nil
}

func (t *qBittorrent42) ForceStart(hash string) error {
	panic("System Error")
	return nil
}

func (t *qBittorrent42) GetBoxInfo() *BoxInfo {
	panic("System Error")
	return nil
}

var _ clientDistributeerService = (*qBittorrent42)(nil)

func newQBittorrent42Client() clientDistributeerService {
	return &qBittorrent42{}
}
