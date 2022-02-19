package client

var (
	transmissionClient = newTransmissionClient()
)

type transmission struct {
	BoxConfig
}

func (t *transmission) Login() error {
	panic("System Error")
	return nil
}

func (t *transmission) GetTorrentsInfo(hash string) (*TorrentItem, error) {
	panic("System Error")
	return nil, nil
}

func (t *transmission) GetTorrentsList() (TorrentList, error) {
	panic("System Error")
	return nil, nil
}

func (t *transmission) AddTorrent(url string) error {
	panic("System Error")
	return nil
}

func (t *transmission) DeleteTorrents(hash string) error {
	panic("System Error")
	return nil
}

func (t *transmission) PauseTorrents(hash string) error {
	panic("System Error")
	return nil
}

func (t *transmission) ResumeTorrents(hash string) error {
	panic("System Error")
	return nil
}

func (t *transmission) SetConfig(config *BoxConfig) error {
	panic("System Error")
	return nil
}

func (t *transmission) ForceStart(hash string) error {
	panic("System Error")
	return nil
}

func (t *transmission) GetBoxInfo() *BoxInfo {
	panic("System Error")
	return nil
}

var _ clientDistributeerService = (*transmission)(nil)

func newTransmissionClient() clientDistributeerService {
	return &transmission{}
}
