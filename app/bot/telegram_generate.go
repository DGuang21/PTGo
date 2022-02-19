package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"ptgo/app/client"
	"ptgo/app/constant"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func (t *tgbot) generateMessageText(item *client.TorrentItem) string {
	text := constant.BoxTemplate
	elem := reflect.ValueOf(item).Elem()
	relType := elem.Type()
	for i := 0; i < relType.NumField(); i++ {
		tag := relType.Field(i).Tag.Get("json")
		value := ""
		switch elem.Field(i).Kind().String() {
		case "int64":
			value = strconv.FormatInt(elem.Field(i).Int(), 10)
		case "int":
			value = strconv.FormatInt(elem.Field(i).Int(), 10)
		case "string":
			value = elem.Field(i).String()
		case "bool":
			value = strconv.FormatBool(elem.Field(i).Bool())
		case "float64":
			value = fmt.Sprintf("%.2f", elem.Field(i).Float())
		}
		switch relType.Field(i).Tag.Get("json") {
		case "state":
			value = client.StateLanguage["cn"]["QBittorrent"][value]
		}
		if t.needByte2Other(tag) {
			value = t.byte2Other(value)
		}
		text = strings.ReplaceAll(text, "{"+tag+"}", value)
	}
	text = strings.ReplaceAll(text, "{time}", time.Now().Format("2006-01-02 15:04:05"))
	return text
}

func (t *tgbot) generateButton(item *client.TorrentItem) tgbotapi.InlineKeyboardMarkup {
	buttons := [][]tgbotapi.InlineKeyboardButton{}
	stateButtons := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⏸️暂停任务", "paused-"+item.Hash),
		tgbotapi.NewInlineKeyboardButtonData("⏯️强制继续", "force-"+item.Hash),
		tgbotapi.NewInlineKeyboardButtonData("🔥 强制删除", "delete-"+item.Hash),
	)
	// 状态控制按钮
	buttons = append(buttons, stateButtons)
	// 附加功能按钮
	buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("✨ 移动GD", "mupload-"+item.Hash),
		tgbotapi.NewInlineKeyboardButtonData("🌟 复制GD", "cupload-"+item.Hash),
		tgbotapi.NewInlineKeyboardButtonData("💫 一键整理", "sortout-"+item.Hash),
	))
	// 广告按钮
	buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL("🦆 凯皇nb", "https://cn.pornhub.com/"),
		tgbotapi.NewInlineKeyboardButtonURL("🐻 炸鸡吧nb", "https://cn.pornhub.com/"),
	))
	buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL("🍊 动漫年鉴", "https://cn.pornhub.com/"),
		tgbotapi.NewInlineKeyboardButtonURL("🐰 支持作者", "https://github.com/DGuang21/PTGo"),
	))
	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}

func (t *tgbot) generateSystemInfo() string {
	text := constant.BoxSystemInfo
	info := client.Box.GetBoxInfo()
	// 拼接box信息
	elem := reflect.ValueOf(info).Elem()
	relType := elem.Type()
	for i := 0; i < relType.NumField(); i++ {
		value := ""
		switch elem.Field(i).Kind().String() {
		case "int64":
			value = strconv.FormatInt(elem.Field(i).Int(), 10)
		case "int":
			value = strconv.FormatInt(elem.Field(i).Int(), 10)
		case "string":
			value = elem.Field(i).String()
		case "bool":
			value = strconv.FormatBool(elem.Field(i).Bool())
		case "float64":
			value = strconv.FormatFloat(elem.Field(i).Float(), 'f', -1, 64)
		}
		tag := relType.Field(i).Tag.Get("json")
		if t.needByte2Other(tag) {
			value = t.byte2Other(value)
		}
		text = strings.ReplaceAll(text, "{"+tag+"}", value)
	}
	text = strings.ReplaceAll(text, "{time}", time.Now().Format("2006-01-02 15:04:05"))
	// 拼接系统信息
	for k, v := range systemInfo {
		if strings.Contains(text, k) {
			text = strings.ReplaceAll(text, k, v())
		}
	}
	return text
}

func (t *tgbot) byte2Other(original interface{}) string {
	switch original.(type) {
	case string:
		i, _ := strconv.ParseFloat(original.(string), 64)
		return t.byte2Others(i)
	case int:
		return t.byte2Others(float64(original.(int)))
	case int64:
		return t.byte2Others(float64(original.(int64)))
	case float64:
		return t.byte2Other(original.(float64))
	}
	return "Unkonw Type"
}

func (t *tgbot) byte2Others(b float64) string {
	if b < 1024 {
		return fmt.Sprintf("%.2f", b) + "Byte"
	}
	if b /= 1024; b < 1024 {
		return fmt.Sprintf("%.2f", b) + "KB"
	}
	if b /= 1024; b < 1024 {
		return fmt.Sprintf("%.2f", b) + "MB"
	}
	if b /= 1024; b < 1024 {
		return fmt.Sprintf("%.2f", b) + "GB"
	}
	if b /= 1024; b < 1024 {
		return fmt.Sprintf("%.2f", b) + "TB"
	}
	if b /= 1024; b < 1024 {
		return fmt.Sprintf("%.2f", b) + "PB"
	}
	return fmt.Sprintf("%.2f", b) + "EB"
}

func (t *tgbot) needByte2Other(a string) bool {
	_, ok := needByte2OtherMap[a]
	return ok
}

var (
	needByte2OtherMap = map[string]bool{
		"amount_left":   true,
		"completed":     true,
		"dl_limit":      true,
		"dlspeed":       true,
		"downloaded":    true,
		"progress":      true,
		"size":          true,
		"total_size":    true,
		"up_limit":      true,
		"uploaded":      true,
		"upspeed":       true,
		"dl_info_speed": true,
		"dl_info_data":  true,
		"up_info_speed": true,
		"up_info_data":  true,
		"dl_rate_limit": true,
		"up_rate_limit": true,
	}
)
