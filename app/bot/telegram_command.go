package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"log"
	"math"
	"ptgo/app/client"
	"ptgo/app/constant"
	"ptgo/app/g"
	"strconv"
	"sync"
	"time"
)

var (
	deleteMap = map[string]int{}
	deleteLock = sync.RWMutex{}
)

func (t *tgbot) delete(hash string, callBack *tgbotapi.CallbackQuery) {
	deleteLock.Lock()
	deleteMap[hash]++
	if deleteMap[hash] != 5 {
		deleteLock.Unlock()
		t.TGBot.AnswerCallbackQuery(tgbotapi.NewCallback(callBack.ID, fmt.Sprintf("防止误触,请继续点击以确认[%v/5]", deleteMap[hash])))
		return
	}
	deleteMap[hash] = 0
	deleteLock.Unlock()
	err := client.Box.DeleteTorrents(hash)
	if err != nil {
		t.TGBot.AnswerCallbackQuery(tgbotapi.NewCallback(callBack.ID, "unkonw button."))
		t.TGBot.Send(tgbotapi.NewMessage(t.UserID, err.Error()))
		return
	}
	var id, msgID int
	var hash_, box string
	constant.DBLock.Lock()
	row, err := g.DB.Query(fmt.Sprintf("select * from main.client where hash='%v'", hash))
	if err != nil {
		log.Println(err.Error())
		return
	}
	for row.Next() {
		if err := row.Scan(&id, &hash_, &msgID, &box); err != nil {
			log.Println(err.Error())
			t.TGBot.AnswerCallbackQuery(tgbotapi.NewCallback(callBack.ID, "删除失败"))
			return
		}
	}
	row.Close()
	tx, err := g.DB.Begin()
	if err != nil {
		log.Println(err.Error())
		return
	}
	if _, err = tx.Exec(fmt.Sprintf("delete from main.client where hash='%v'", hash)); err != nil {
		log.Println(err.Error())
		if err := tx.Rollback(); err != nil {
			log.Println(err.Error())
			t.TGBot.AnswerCallbackQuery(tgbotapi.NewCallback(callBack.ID, "删除失败"))
		}
		return
	}
	if err := tx.Commit(); err != nil {
		log.Println(err.Error())
		return
	}
	constant.DBLock.Unlock()
	t.TGBot.AnswerCallbackQuery(tgbotapi.NewCallback(callBack.ID, "删除成功"))
	t.TGBot.DeleteMessage(tgbotapi.DeleteMessageConfig{
		ChatID:    t.UserID,
		MessageID: msgID,
	})
}

func (t *tgbot) pause(hash string, callBack *tgbotapi.CallbackQuery) {
	err := client.Box.PauseTorrents(hash)
	if err != nil {
		t.TGBot.AnswerCallbackQuery(tgbotapi.NewCallback(callBack.ID, "暂停失败"))
		log.Println(err.Error())
		return
	}
	t.TGBot.AnswerCallbackQuery(tgbotapi.NewCallback(callBack.ID, "暂停成功"))
}

func (t *tgbot) resume(hash string, callBack *tgbotapi.CallbackQuery) {
	err := client.Box.ResumeTorrents(hash)
	if err != nil {
		t.TGBot.AnswerCallbackQuery(tgbotapi.NewCallback(callBack.ID, "暂停失败"))
		log.Println(err.Error())
		return
	}
	t.TGBot.AnswerCallbackQuery(tgbotapi.NewCallback(callBack.ID, "恢复成功"))
}

func (t *tgbot) force(hash string, callBack *tgbotapi.CallbackQuery) {
	err := client.Box.ForceStart(hash)
	if err != nil {
		t.TGBot.AnswerCallbackQuery(tgbotapi.NewCallback(callBack.ID, "强制继续失败"))
		log.Println(err.Error())
		return
	}
	t.TGBot.AnswerCallbackQuery(tgbotapi.NewCallback(callBack.ID, "强制继续成功"))
}

func (t *tgbot) cupload(hash string, callBack *tgbotapi.CallbackQuery) {
	t.TGBot.AnswerCallbackQuery(tgbotapi.NewCallback(callBack.ID, "还没做呢"))
}

func (t *tgbot) dupload(hash string, callBack *tgbotapi.CallbackQuery) {
	t.TGBot.AnswerCallbackQuery(tgbotapi.NewCallback(callBack.ID, "还没做呢"))
}

func (t *tgbot) sortOut(hash string, callBack *tgbotapi.CallbackQuery) {
	t.TGBot.AnswerCallbackQuery(tgbotapi.NewCallback(callBack.ID, "还没做呢"))
}

func (t *tgbot) getCpuInfo() string {
	percent, _ := cpu.Percent(time.Second, false)
	return strconv.FormatFloat(math.Trunc(percent[0]*1e2+0.5)*1e-2, 'f', -1, 64) + "%"
}

func (t *tgbot) getRamInfo() string {
	memInfo, _ := mem.VirtualMemory()
	return strconv.FormatFloat(math.Trunc(memInfo.UsedPercent*1e2+0.5)*1e-2, 'f', -1, 64) + "%"
}

func (t *tgbot) getIOInfo() string {
	parts, err := disk.Partitions(true)
	if err != nil {
		return err.Error()
	}
	resp := ""
	for _, part := range parts {
		diskInfo, _ := disk.Usage(part.Mountpoint)
		resp += part.Device + strconv.FormatFloat(math.Trunc(diskInfo.UsedPercent*1e2+0.5)*1e-2, 'f', -1, 64) + "%" + "  "
	}
	return resp
}
