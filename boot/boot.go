package boot

import (
	"github.com/BurntSushi/toml"
	"log"
	"os"
	"ptgo/app/bot"
	"ptgo/app/client"
	"ptgo/app/g"
)

func init() {
	var config tomlConfig
	str, _ := os.Getwd()
	_, err := toml.DecodeFile(str+"/config/config.toml", &config)
	if err != nil {
		log.Println("init config fail:", err.Error())
		return
	} else if config.Refresh == 0 {
		config.Refresh = 60
	}
	// 初始化数据库连接池
	if err := g.DB.SetConfig(str + "/config/data.db"); err != nil {
		log.Println("init DB fail:", err.Error())
		return
	}
	log.Println("数据库初始化成功")
	// 初始化pt-go bot
	if err := bot.Telegram.SetTGBotConfig(&bot.TGBotConfig{
		UserID:     config.UserID,
		TGBotToken: config.OperationBotToken,
		Debug:      config.Debug,
	}); err != nil {
		log.Println("初始化bot失败:", err.Error())
	} else {
		if err := bot.Telegram.Run(); err != nil {
			log.Println("初始化bot失败:", err.Error())
		}
	}
	// 初始化rss推送功能
	if config.PTConfig != nil {
		if err := bot.Channel.SetChannelConfig(config.PTRss, config.ChannelID, config.OperationBotName); err != nil {
			log.Println("初始化推送失败,", err.Error())
		} else {
			go bot.Channel.Refresh(config.Refresh)
			log.Println("rss推送初始化成功")
		}
	}
	// 初始化PT客户端
	if config.TelegramConfig != nil {
		if err := client.Box.SetConfig(&client.BoxConfig{
			BoxName:     config.BoxName,
			BoxUrl:      config.BoxUrl,
			BoxUserName: config.BoxUserName,
			BoxPassWord: config.BoxPassWord,
		}); err != nil {
			log.Println("初始化box连接失败,", err.Error())
		} else {
			log.Println("初始化box成功")
		}
	} else {
		log.Println("无box config，跳过")
	}
	log.Println("初始化完成")
}

type tomlConfig struct {
	*PTConfig       `toml:"box"`
	*TelegramConfig `toml:"telegram"`
}

type PTConfig struct {
	BoxName     string
	BoxUrl      string
	BoxUserName string
	BoxPassWord string
	PTRss       string
	Refresh     int
}

type TelegramConfig struct {
	UserID            int64
	ChannelID         string
	OperationBotToken string
	OperationBotName  string
	Debug             bool
}
