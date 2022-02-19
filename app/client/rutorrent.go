package client

var (
	rutorrentClient = newRutorrentClient()
)

type rutorrent struct {
	BoxConfig
}

func (t *rutorrent) Login() error {
	panic("System Error")
	return nil
}

func (t *rutorrent) GetTorrentsInfo(hash string) (*TorrentItem, error) {
	panic("System Error")
	return nil, nil
}

func (t *rutorrent) GetTorrentsList() (TorrentList, error) {
	panic("System Error")
	return nil, nil
}

func (t *rutorrent) AddTorrent(url string) error {
	panic("System Error")
	return nil
}

func (t *rutorrent) DeleteTorrents(hash string) error {
	panic("System Error")
	return nil
}

func (t *rutorrent) PauseTorrents(hash string) error {
	panic("System Error")
	return nil
}

func (t *rutorrent) ResumeTorrents(hash string) error {
	panic("System Error")
	return nil
}

func (t *rutorrent) SetConfig(config *BoxConfig) error {
	panic("System Error")
	return nil
}

func (t *rutorrent) ForceStart(hash string) error {
	panic("System Error")
	return nil
}

func (t *rutorrent) GetBoxInfo() *BoxInfo {
	panic("System Error")
	return nil
}

var _ clientDistributeerService = (*qBittorrent)(nil)

func newRutorrentClient() clientDistributeerService {
	return &rutorrent{}
}
