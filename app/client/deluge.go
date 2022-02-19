package client

var (
	delugeClient = newDelugeClient()
)

type deluge struct {
	BoxConfig
}

func (t *deluge) Login() error {
	panic("System Error")
	return nil
}

func (t *deluge) GetTorrentsInfo(hash string) (*TorrentItem, error) {
	panic("System Error")
	return nil, nil
}

func (t *deluge) GetTorrentsList() (TorrentList, error) {
	panic("System Error")
	return nil, nil
}

func (t *deluge) AddTorrent(url string) error {
	panic("System Error")
	return nil
}

func (t *deluge) DeleteTorrents(hash string) error {
	panic("System Error")
	return nil
}

func (t *deluge) PauseTorrents(hash string) error {
	panic("System Error")
	return nil
}

func (t *deluge) ResumeTorrents(hash string) error {
	panic("System Error")
	return nil
}

func (t *deluge) SetConfig(config *BoxConfig) error {
	panic("System Error")
	return nil
}

func (t *deluge) ForceStart(hash string) error {
	panic("System Error")
	return nil
}

func (t *deluge) GetBoxInfo() *BoxInfo {
	panic("System Error")
	return nil
}

var _ clientDistributeerService = (*deluge)(nil)

func newDelugeClient() clientDistributeerService {
	return &deluge{}
}
