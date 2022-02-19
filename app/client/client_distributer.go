package client

type clientDistributeerService interface {
	Login() error

	GetTorrentsInfo(hash string) (*TorrentItem, error)

	GetTorrentsList() (TorrentList, error)

	AddTorrent(url string) error

	DeleteTorrents(hash string) error

	PauseTorrents(hash string) error

	ResumeTorrents(hash string) error

	SetConfig(config *BoxConfig) error

	ForceStart(hash string) error

	GetBoxInfo() *BoxInfo
}

var (
	StateLanguage = map[string]map[string]map[string]string{
		"cn": {"QBittorrent": {
			"error":              "错误",
			"missingFiles":       "文件丢失",
			"uploading":          "做种",
			"pausedUP":           "暂停做种",
			"queuedUP":           "等待上传",
			"stalledUP":          "正在做种，但没有建立连接",
			"checkingUP":         "检查中",
			"forcedUP":           "[F]做种",
			"allocating":         "分配磁盘",
			"downloading":        "下载中",
			"metaDL":             "获取元数据",
			"pausedDL":           "暂停下载",
			"queuedDL":           "等待下载",
			"stalledDL":          "正在下载，但没有建立连接",
			"checkingDL":         "下载检查",
			"forcedDL":           "[F]下载",
			"checkingResumeData": "检查中",
			"moving":             "移动中",
			"unknown":            "未知状态",
		}},
	} // 语言 - 软件
)

type BoxInfo struct {
	DownloadInfoSpeed int    `json:"dl_info_speed"`
	DownloadInfoData  int    `json:"dl_info_data"`
	UploadInfoSpeed   int    `json:"up_info_speed"`
	UploadInfoData    int    `json:"up_info_data"`
	DownloadRateLimit int    `json:"dl_rate_limit"`
	UploadRateLimit   int    `json:"up_rate_limit"`
	DTHNode           int    `json:"dht_nodes"`
	BoxName           string `json:"box_name"`
}

type TorrentList []*TorrentItem

type TorrentItem struct {
	AddedOn           int     `json:"added_on"`           // 将 torrent 添加到客户端的时间（Unix Epoch）
	AmountLeft        int64   `json:"amount_left"`        // 剩余大小（字节）
	AutoTmm           bool    `json:"auto_tmm"`           // 此 torrent 是否由 Automatic Torrent Management 管理
	Availability      float64 `json:"availability"`       // 当前百分比
	Category          string  `json:"category"`           //
	Completed         int64   `json:"completed"`          // 完成的传输数据量（字节）
	CompletionOn      int     `json:"completion_on"`      // Torrent 完成的时间（Unix Epoch）
	ContentPath       string  `json:"content_path"`       // torrent 内容的绝对路径（多文件 torrent 的根路径，单文件 torrent 的绝对文件路径）
	DlLimit           int     `json:"dl_limit"`           // Torrent 下载速度限制（字节/秒）
	Dlspeed           int     `json:"dlspeed"`            // Torrent 下载速度（字节/秒）
	Downloaded        int64   `json:"downloaded"`         // 已经下载大小
	DownloadedSession int64   `json:"downloaded_session"` // 此会话下载的数据量
	Eta               int     `json:"eta"`                //
	FLPiecePrio       bool    `json:"f_l_piece_prio"`     // 如果第一个最后一块被优先考虑，则为true
	ForceStart        bool    `json:"force_start"`        // 如果为此 torrent 启用了强制启动，则为true
	Hash              string  `json:"hash"`               //
	LastActivity      int     `json:"last_activity"`      // 上次活跃的时间（Unix Epoch）
	MagnetURI         string  `json:"magnet_uri"`         // 与此 torrent 对应的 Magnet URI
	MaxRatio          int     `json:"max_ratio"`          // 种子/上传停止种子前的最大共享比率
	MaxSeedingTime    int     `json:"max_seeding_time"`   // 停止种子种子前的最长种子时间（秒）
	Name              string  `json:"name"`               //
	NumComplete       int     `json:"num_complete"`       //
	NumIncomplete     int     `json:"num_incomplete"`     //
	NumLeechs         int     `json:"num_leechs"`         // 连接到的 leechers 的数量
	NumSeeds          int     `json:"num_seeds"`          // 连接到的种子数
	Priority          int     `json:"priority"`           // 速度优先。如果队列被禁用或 torrent 处于种子模式，则返回 -1
	Progress          float64 `json:"progress"`           // 进度
	Ratio             float64 `json:"ratio"`              // Torrent 共享比率
	RatioLimit        int     `json:"ratio_limit"`        //
	SavePath          string  `json:"save_path"`
	SeedingTime       int     `json:"seeding_time"`       // Torrent 完成用时（秒）
	SeedingTimeLimit  int     `json:"seeding_time_limit"` // max_seeding_time
	SeenComplete      int     `json:"seen_complete"`      // 上次 torrent 完成的时间
	SeqDl             bool    `json:"seq_dl"`             // 如果启用顺序下载，则为true
	Size              int64   `json:"size"`               //
	State             string  `json:"state"`              // 参见https://github.com/qbittorrent/qBittorrent/wiki/WebUI-API-(qBittorrent-4.1)#get-torrent-list
	SuperSeeding      bool    `json:"super_seeding"`      // 如果启用超级播种，则为true
	Tags              string  `json:"tags"`               // Torrent 的逗号连接标签列表
	TimeActive        int     `json:"time_active"`        // 总活动时间（秒）
	TotalSize         int64   `json:"total_size"`         // 此 torrent 中所有文件的总大小（字节）（包括未选择的文件）
	Tracker           string  `json:"tracker"`            // 第一个具有工作状态的tracker。如果没有tracker在工作，则返回空字符串。
	TrackersCount     int     `json:"trackers_count"`     //
	UpLimit           int     `json:"up_limit"`           // 上传限制
	Uploaded          int64   `json:"uploaded"`           // 累计上传
	UploadedSession   int64   `json:"uploaded_session"`   // 当前session累计上传
	Upspeed           int     `json:"upspeed"`            // 上传速度（字节/秒）
}
