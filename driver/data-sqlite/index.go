package data_sqlite

import (
	"github.com/nogio/noggo/driver"
	_ "github.com/nogio/noggo/driver/data-sqlite/go-sqlite3"   //此包自动注册名为sqlite的sql驱动
)

const (
	SQLDRIVER = "sqlite3"
)

//返回驱动
func Driver() (driver.DataDriver) {
	return &SqliteDriver{}
}