package data_pgsql

import (
	"github.com/nogio/noggo/driver"
	_ "github.com/lib/pq"   //此包自动注册名为postgres的sql驱动
)


//返回驱动
func Driver() (driver.DataDriver) {
	return &PgsqlDriver{}
}