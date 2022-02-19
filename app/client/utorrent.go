package client

var (
	utorrentClient = newUtorrentClient()
)

type utorrent struct {
	BoxConfig
}

func (t *utorrent) Login() error {
	panic("System Error")
	return nil
}

func (t *utorrent) GetTorrentsInfo(hash string) (*TorrentItem, error) {
	panic("System Error")
	return nil, nil
}

func (t *utorrent) GetTorrentsList() (TorrentList, error) {
	panic("System Error")
	return nil, nil
}

func (t *utorrent) AddTorrent(url string) error {
	panic("System Error")
	return nil
}

func (t *utorrent) DeleteTorrents(hash string) error {
	panic("System Error")
	return nil
}

func (t *utorrent) PauseTorrents(hash string) error {
	panic("System Error")
	return nil
}

func (t *utorrent) ResumeTorrents(hash string) error {
	panic("System Error")
	return nil
}

func (t *utorrent) SetConfig(config *BoxConfig) error {
	panic("System Error")
	return nil
}

func (t *utorrent) ForceStart(hash string) error {
	panic("System Error")
	return nil
}

func (t *utorrent) GetBoxInfo() *BoxInfo {
	panic("System Error")
	return nil
}

var _ clientDistributeerService = (*utorrent)(nil)

func newUtorrentClient() clientDistributeerService {
	return &utorrent{}
}
