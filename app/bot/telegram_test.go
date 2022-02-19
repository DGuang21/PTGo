package bot

import "testing"

func TestNewTGBot(t *testing.T) {
	t.Log(Telegram.getCpuInfo())
	t.Log(Telegram.getRamInfo())
	t.Log(Telegram.getIOInfo())
}
