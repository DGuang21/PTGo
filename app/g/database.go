package g

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var (
	DB = NewGDB()
)

type gdb struct {
	*sql.DB
}

func NewGDB() *gdb {
	return &gdb{}
}

func (t *gdb) SetConfig(path string) error {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return err
	}
	t.DB = db
	return nil
}
