package bot

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"ptgo/app/client"
	"ptgo/app/constant"
	"ptgo/app/g"
	"strings"
	"time"
)

var (
	Telegram     = NewTGBot()
	err          error
	callBackFunc = map[string]func(hash string, callBack *tgbotapi.CallbackQuery){
		"delete":  Telegram.delete,
		"paused":  Telegram.pause,
		"force":   Telegram.force,
		"cupload": Telegram.cupload,
		"dupload": Telegram.dupload,
		"sortout": Telegram.sortOut,
	}
	systemInfo = map[string]func() string{
		"{CPU}": Telegram.getCpuInfo,
		"{RAM}": Telegram.getRamInfo,
		"{IO}":  Telegram.getIOInfo,
	}
)

type TGBotConfig struct {
	UserID     int64
	TGBotToken string
	Debug      bool
}

type tgbot struct {
	TGBot *tgbotapi.BotAPI
	*TGBotConfig
}

func NewTGBot() *tgbot {
	return &tgbot{}
}

func (t *tgbot) SetTGBotConfig(config *TGBotConfig) error {
	switch {
	case config.UserID == 0:
		return errors.New("need UserID")
	case config.TGBotToken == "":
		return errors.New("need TGBotToken")
	default:
		t.TGBotConfig = config
		t.TGBot, err = tgbotapi.NewBotAPI(t.TGBotToken)
		return err
	}
}

func (t *tgbot) Run() error {
	t.TGBot.Debug = t.Debug
	log.Printf("telegram bot初始化成功:@%s", t.TGBot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := t.TGBot.GetUpdatesChan(u)
	if err != nil {
		return err
	}
	go func() {
		for update := range updates {
			if update.Message != nil {
				t.reply(update.Message)
			} else if update.CallbackQuery != nil {
				t.replyCallBack(update.CallbackQuery)
			}

		}
	}()
	go func() {
		ti := time.NewTicker(time.Second * 30)
		for {
			select {
			case <-ti.C:
				t.refreshBox()
			}
		}
	}()
	return nil
}

func (t *tgbot) refreshBox() {
	list, err := client.Box.GetTorrentsList()
	if err != nil {
		log.Println(err.Error())
		return
	}
	t.updateState()
	for _, item := range list {
		// 查询是否有记录
		constant.DBLock.Lock()
		row, err := g.DB.Query(fmt.Sprintf("select * from main.client where hash='%v'", item.Hash))
		constant.DBLock.Unlock()
		if err != nil {
			log.Println("sqlite error:", err.Error())
			return
		}
		// 有记录
		if row.Next() {
			var id, msgID int
			var hash, box string
			if err := row.Scan(&id, &hash, &msgID, &box); err != nil {
				log.Println("sqlite error:", err.Error())
				return
			}
			t.updateTask(item, msgID)
			continue
		}
		row.Close()
		// 无记录
		t.newTask(item)
	}
}

func (t *tgbot) reply(update *tgbotapi.Message) {
	if int64(update.From.ID) != t.TGBotConfig.UserID {
		t.TGBot.Send(tgbotapi.NewMessage(int64(update.From.ID), "fuck you"))
		return
	}
	fmt.Printf("%+v\n", update)
}

func (t *tgbot) newTask(item *client.TorrentItem) {
	msg := tgbotapi.NewMessage(t.UserID, t.generateMessageText(item))
	msg.ReplyMarkup = t.generateButton(item)
	result, err := t.TGBot.Send(msg)
	if err != nil {
		log.Println("sendMessage Fail:", err.Error())
		return
	}
	constant.DBLock.Lock()
	tx, err := g.DB.Begin()
	if err != nil {
		log.Println("sqlite error:", err.Error())
		return
	}
	if _, err = tx.Exec(fmt.Sprintf("insert into main.client (hash,msg_id,box) VALUES ('%v','%v','%v')", item.Hash, result.MessageID, "default")); err != nil {
		log.Println("sqlite error:", err.Error())
		constant.DBLock.Unlock()
		if err := tx.Rollback(); err != nil {
			return
		}
		return
	}
	if err := tx.Commit(); err != nil {
		constant.DBLock.Unlock()
		return
	}
	constant.DBLock.Unlock()
}

func (t *tgbot) updateTask(item *client.TorrentItem, id int) {
	text := t.generateMessageText(item)
	msg := tgbotapi.NewEditMessageText(t.UserID, id, text)
	button := t.generateButton(item)
	msg.ReplyMarkup = &button
	if _, err := t.TGBot.Send(msg); err != nil {
		log.Println("editMessage Fail:", err.Error())
		return
	}
}

func (t *tgbot) replyCallBack(callBack *tgbotapi.CallbackQuery) {
	info := strings.Split(callBack.Data, "-")
	if len(info) != 2 {
		t.TGBot.AnswerCallbackQuery(tgbotapi.NewCallback(callBack.ID, "unkonw button."))
		return
	}
	command, hash := info[0], info[1]
	method, ok := callBackFunc[command]
	if !ok {
		t.TGBot.AnswerCallbackQuery(tgbotapi.NewCallback(callBack.ID, "unkonw button."))
		return
	}
	go method(hash, callBack)
}

func (t *tgbot) updateState() {
	text := t.generateSystemInfo()
	constant.DBLock.Lock()
	row, err := g.DB.Query(fmt.Sprintf("select * from main.client where hash='%v'", "state"))
	constant.DBLock.Unlock()
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer row.Close()
	// 有记录，直接update
	if row.Next() {
		var id, msgID int
		var hash, box string
		if err := row.Scan(&id, &hash, &msgID, &box); err != nil {
			log.Println(err.Error())
			return
		}
		if _, err := t.TGBot.Send(tgbotapi.NewEditMessageText(t.UserID, msgID, text)); err != nil {
			log.Println(err.Error())
		}
		return
	}
	// 无记录
	msg := tgbotapi.NewMessage(t.UserID, text)
	result, err := t.TGBot.Send(msg)
	if err != nil {
		log.Println(err.Error())
		return
	}
	if _, err := t.TGBot.PinChatMessage(tgbotapi.PinChatMessageConfig{
		ChatID:              t.UserID,
		MessageID:           result.MessageID,
		DisableNotification: true,
	}); err != nil {
		log.Println(err.Error())
		return
	}
	// 保存
	constant.DBLock.Lock()
	tx, err := g.DB.Begin()
	if err != nil {
		log.Println("sqlite error:", err.Error())
		return
	}
	if _, err = tx.Exec(fmt.Sprintf("insert into main.client (hash,msg_id,box) VALUES ('%v','%v','%v')", "state", result.MessageID, "default")); err != nil {
		log.Println("sqlite error:", err.Error())
		constant.DBLock.Unlock()
		if err := tx.Rollback(); err != nil {
			return
		}
		return
	}
	if err := tx.Commit(); err != nil {
		constant.DBLock.Unlock()
		return
	}
	constant.DBLock.Unlock()
}
