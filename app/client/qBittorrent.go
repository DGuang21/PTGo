package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"ptgo/app/tool"
)

var (
	qBittorrentClient = newQBittorrentClient()
)

type qBittorrent struct {
	*BoxConfig
}

func (t *qBittorrent) Login() error {
	header := map[string]string{"Content-type": "application/x-www-form-urlencoded; charset=UTF-8"}
	param := fmt.Sprintf("username=%v&password=%v", t.BoxUserName, t.BoxPassWord)
	data, err := tool.Ghttp.Post(t.BoxUrl+"/api/v2/auth/login", param, header)
	if err != nil {
		return err
	}
	if data["code"] != "200" {
		return errors.New(data["response"])
	}
	if data["cookie"] == "" {
		return errors.New(data["response"])
	}
	t.BoxCookie = data["cookie"]
	return nil
}

func (t *qBittorrent) GetTorrentsInfo(hash string) (*TorrentItem, error) {
	header := map[string]string{
		"Cookie": t.BoxCookie,
		"Accept": "",
	}
	data, err := tool.Ghttp.Get(t.BoxUrl+"/api/v2/torrents/properties?hash="+hash, header)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	if data["code"] == "403" || data["code"] == "401" {
		log.Printf("get info error:%v.retry login", err)
		if err := t.Login(); err != nil {
			log.Println(err)
			return nil, err
		}
		return t.GetTorrentsInfo(hash)
	}
	if data["code"] == "404" {
		return nil, errors.New("不存在的hash")
	}
	var info TorrentItem
	err = json.Unmarshal([]byte(data["response"]), &info)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return &info, nil
}

func (t *qBittorrent) GetTorrentsList() (TorrentList, error) {
	header := map[string]string{
		"Cookie": t.BoxCookie,
		"Accept": "application/json",
	}
	data, err := tool.Ghttp.Get(t.BoxUrl+"/api/v2/torrents/info", header)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	if data["code"] == "403" || data["code"] == "401" {
		log.Printf("get list error:%v.retry login", err)
		if err := t.Login(); err != nil {
			log.Println(err)
			return nil, err
		}
		return t.GetTorrentsList()
	}
	if len(data["response"]) == 0 {
		return nil, err
	}
	var list TorrentList
	err = json.Unmarshal([]byte(data["response"]), &list)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return list, nil
}

func (t *qBittorrent) AddTorrent(url string) error {
	header := map[string]string{
		"Cookie": t.BoxCookie,
		"Accept": "application/json",
	}
	param := fmt.Sprintf("urls:%v", url)
	resp, err := tool.Ghttp.Post(t.BoxUrl+"/api/v2/torrents/add", param, header)
	if err != nil {
		log.Println(err)
		return err
	}
	if resp["code"] == "403" || resp["code"] == "401" {
		log.Printf("get list error:%v.retry login", err)
		if err := t.Login(); err != nil {
			log.Println(err)
			return err
		}
		return t.AddTorrent(url)
	}
	if resp["response"] != "Ok." {
		log.Println(resp)
		return errors.New("unkonw error")
	}
	return nil
}

func (t *qBittorrent) DeleteTorrents(hash string) error {
	header := map[string]string{
		"Content-type": "application/x-www-form-urlencoded; charset=UTF-8",
		"Cookie":       t.BoxCookie,
	}
	param := fmt.Sprintf("hashes=%v&deleteFiles=true", hash)
	data, err := tool.Ghttp.Post(t.BoxUrl+"/api/v2/torrents/delete", param, header)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	if data["code"] == "403" || data["code"] == "401" {
		log.Printf("get list error:%v.retry login", err)
		if err := t.Login(); err != nil {
			log.Println(err)
			return err
		}
		return t.DeleteTorrents(hash)
	}
	return nil
}

// 暂停种子
func (t *qBittorrent) PauseTorrents(hash string) error {
	header := map[string]string{
		"Content-type": "application/x-www-form-urlencoded; charset=UTF-8",
		"Cookie":       t.BoxCookie,
	}
	param := fmt.Sprintf("hashes=%v", hash)
	resp, err := tool.Ghttp.Post(t.BoxUrl+"/api/v2/torrents/pause", param, header)
	if err != nil {
		log.Println(err)
		return err
	}
	if resp["code"] == "403" || resp["code"] == "401" {
		log.Printf("get list error:%v.retry login", err)
		if err := t.Login(); err != nil {
			log.Println(err)
			return err
		}
		return t.PauseTorrents(hash)
	}
	return nil
}

// 恢复种子下载
func (t *qBittorrent) ResumeTorrents(hash string) error {
	header := map[string]string{
		"Content-type": "application/x-www-form-urlencoded; charset=UTF-8",
		"Cookie":       t.BoxCookie,
	}
	param := fmt.Sprintf("hashes=%v", hash)
	resp, err := tool.Ghttp.Post(t.BoxUrl+"/api/v2/torrents/resume", param, header)
	if err != nil {
		log.Println(err)
		return err
	}
	if resp["code"] == "403" || resp["code"] == "401" {
		if err := t.Login(); err != nil {
			log.Println(err)
			return err
		}
		return t.ResumeTorrents(hash)
	}
	return nil
}

func (t *qBittorrent) ForceStart(hash string) error {
	header := map[string]string{
		"Content-type": "application/x-www-form-urlencoded; charset=UTF-8",
		"Cookie":       t.BoxCookie,
	}
	param := fmt.Sprintf("value=true&hashes=%v", hash)
	resp, err := tool.Ghttp.Post(t.BoxUrl+"/api/v2/torrents/setForceStart", param, header)
	if err != nil {
		log.Println(err)
		return err
	}
	if resp["code"] == "403" || resp["code"] == "401" {
		if err := t.Login(); err != nil {
			log.Println(err)
			return err
		}
		return t.ForceStart(hash)
	}
	return nil
}

func (t *qBittorrent) GetBoxInfo() *BoxInfo {
	return t.getBoxState()
}

func (t *qBittorrent) getBoxState() *BoxInfo {
	var state *BoxInfo
	header := map[string]string{
		"Content-type": "application/x-www-form-urlencoded; charset=UTF-8",
		"Cookie":       t.BoxCookie,
	}
	resp, err := tool.Ghttp.Post(t.BoxUrl+"/api/v2/transfer/info", "", header)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	if resp["code"] == "403" || resp["code"] == "401" {
		if err := t.Login(); err != nil {
			log.Println(err)
			return nil
		}
		return t.getBoxState()
	}
	if err := json.Unmarshal([]byte(resp["response"]), &state); err != nil {
		log.Println(err.Error())
		return nil
	}
	state.BoxName = t.getBoxVersion()
	return state
}

func (t *qBittorrent) getBoxVersion() string {
	header := map[string]string{
		"Content-type": "application/x-www-form-urlencoded; charset=UTF-8",
		"Cookie":       t.BoxCookie,
	}
	resp, err := tool.Ghttp.Post(t.BoxUrl+"/api/v2/app/version", "", header)
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	if resp["code"] == "403" || resp["code"] == "401" {
		if err := t.Login(); err != nil {
			log.Println(err)
			return ""
		}
		return t.getBoxVersion()
	}
	return t.BoxName + "  " + resp["response"]
}

func (t *qBittorrent) SetConfig(config *BoxConfig) error {
	t.BoxConfig = config
	return t.Login()
}

var _ clientDistributeerService = (*qBittorrent)(nil)

func newQBittorrentClient() clientDistributeerService {
	return &qBittorrent{}
}
