package bot

import (
	"crypto/md5"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"golang.org/x/net/html"
	"log"
	"ptgo/app/constant"
	"ptgo/app/g"
	"ptgo/app/tool"
	"ptgo/app/vo"
	"strconv"
	"strings"
	"time"
)

var (
	Channel = NewChannel()
)

type channel struct {
	RssLink          []string
	ChannelID        string
	OperationBotName string
	PushBot          *tgbotapi.BotAPI
}

func NewChannel() *channel {
	return &channel{}
}

func (c *channel) SetChannelConfig(link string, channelID string, botName string) error {
	switch {
	case link == "" || channelID == "" || botName == "":
		return errors.New("unknown Link Or Token")
	default:
		c.ChannelID = channelID
		c.RssLink = strings.Split(link, ",")
		c.OperationBotName = strings.ReplaceAll(botName, "@", "")
		// init push bot
		if Telegram.TGBot == nil {
			return errors.New("请填写TG-bot Token")
		}
		c.PushBot = Telegram.TGBot
		return nil
	}
}

func (c *channel) Refresh(t int) {
	ti := time.NewTicker(time.Second * time.Duration(t))
	for {
		select {
		case <-ti.C:
			for _, v := range c.RssLink {
				c.refresh(v)
			}
		}
	}
}

func (c *channel) refreshChannelItem(rssData []vo.RssItem,index int) error {
	if len(rssData) < index + 1 {
		return nil
	}
	ptHash := fmt.Sprintf("%x", md5.Sum([]byte(rssData[index].Enclosure.URL)))
	constant.DBLock.RLock()
	row, err := g.DB.Query(fmt.Sprintf("select * from main.history where hash='%v'", ptHash))
	constant.DBLock.RUnlock()
	if err != nil {
		log.Println("DB Error:", err.Error())
		return c.refreshChannelItem(rssData,index + 1)
	}
	if !row.Next() {
		// 解析PT中的html获取第一张图片
		node, err := goquery.NewDocumentFromReader(strings.NewReader(rssData[index].Description))
		if err != nil {
			log.Println("Parse Fail,", err.Error())
			return c.refreshChannelItem(rssData,index + 1)
		}
		// default img
		imgLink := ""
		if node != nil {
			imgs := node.Find("img")
			imgLink = c.findFirstImg(imgs.Nodes,0,rssData[index].Link)
		}
		err = c.sendChannel(imgLink,ptHash,rssData[index])
		if err != nil {
			return c.refreshChannelItem(rssData,index + 1)
		}
		constant.DBLock.Lock()
		tx, err := g.DB.Begin()
		if err != nil {
			log.Println("Data Error:", err.Error())
			return c.refreshChannelItem(rssData,index + 1)
		}
		_, err = tx.Exec(fmt.Sprintf("insert into main.history (hash,url) VALUES ('%v','%v')", ptHash, rssData[index].Enclosure.URL))
		if err != nil {
			log.Println("Data Error:", err.Error())
			tx.Rollback()
			constant.DBLock.Unlock()
			return c.refreshChannelItem(rssData,index + 1)
		}
		tx.Commit()
		constant.DBLock.Unlock()
	}
	row.Close()
	return c.refreshChannelItem(rssData,index + 1)
}

func (c *channel) findFirstImg(htmlData []*html.Node,index int,prefix string) string {
	if len(htmlData) < index + 1 {
		return constant.DefaultIMG
	}
	if imgLink := c.findHtmlAttr(htmlData[index].Attr,0,prefix); len(imgLink) != 0 {
		return imgLink
	}
	return c.findFirstImg(htmlData,index+1,prefix)
}

func (c *channel) findHtmlAttr(data []html.Attribute,index int,prefix string) string {
	if len(data) < index + 1 {
		return constant.DefaultIMG
	}
	link := c.parseImg(data[index].Val, prefix)
	if link != "" {
		return link
	}
	return c.findHtmlAttr(data,index + 1,prefix)
}

func (c *channel) sendChannel(imgLink,ptHash string,v vo.RssItem) error {
	msg := tgbotapi.NewPhotoShare(0, "")
	msg.ChannelUsername = c.ChannelID
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("👀 前去查看", v.Link),
			tgbotapi.NewInlineKeyboardButtonURL("⏬ 一键下载", "http://t.me/"+c.OperationBotName+"?start=down-"+ptHash),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("🤓 凯皇牛逼", "https://cn.pornhub.com/"),
			tgbotapi.NewInlineKeyboardButtonURL("🍆 炸鸡吧牛逼", "https://t.me/EmbyPublic"),
		),
	)
	msg.FileID = imgLink
	text := constant.PushTemplate
	text = strings.ReplaceAll(text, "{_name}", v.Title)
	size, _ := strconv.Atoi(v.Enclosure.Length)
	text = strings.ReplaceAll(text, "{_size}", strconv.Itoa(size/1024/1024/1024))
	text = strings.ReplaceAll(text, "{_from}", strings.ReplaceAll(v.Title, "Torrents", ""))
	msg.Caption = text
	msg.ParseMode = "Markdown"
	_, err = c.PushBot.Send(msg)
	if err != nil {
		return errors.New("sendMessage Fail:" +  err.Error())
	}
	return nil
}

// 开始刷新
func (c *channel) refresh(rssUrl string) {
	resp, err := tool.Ghttp.Get(rssUrl, nil)
	if err != nil {
		log.Println("Get RSS Info Fail,", err.Error())
		return
	}
	var info vo.Rss
	err = xml.Unmarshal([]byte(resp["response"]), &info)
	if err != nil {
		log.Println("Unmarshal RSS Fail,", err.Error())
		return
	}
	err = c.refreshChannelItem(info.Channel.Item,0)
	if err != nil {
		log.Println("refresh RSS Fail,", err.Error())
		return
	}
}

// img 识别
func (c *channel) parseImg(img string, link string) string {
	if strings.Index(img, "http") != -1 {
		if strings.Index(img, "i.loli.net") == -1 {
			// 有图片直链并且不为sm.ms
			return img
		}
	}
	if strings.Index(img, "attachments") != -1 {
		return link + "/" + img
	}
	if strings.Index(img, "/attachments") != -1 {
		return link + img
	}
	return ""
}
