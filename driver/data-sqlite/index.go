package data_sqlite

import (
	"github.com/nogio/noggo"
	_ "github.com/nogio/noggo/depend/go-sqlite3"   //此包自动注册名为sqlite的sql驱动
)

const (
	SQLDRIVER = "sqlite3"
)

//返回驱动
func Driver() (noggo.DataDriver) {
	return &SqliteDriver{}
}