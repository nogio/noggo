package data_sqlite

import (
	"github.com/nogio/noggo/driver"
	_ "github.com/mattn/go-sqlite3"   //此包自动注册名为sqlite的sql驱动
)

const (
	SQLDRIVER = "sqlite"
)

//返回驱动
func Driver() (driver.DataDriver) {
	return &SqliteDriver{}
}