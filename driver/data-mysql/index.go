package data_mysql

import (
	"github.com/nogio/noggo/driver"
	_ "github.com/go-sql-driver/mysql"   //此包自动注册名为mysql的sql驱动
)

const (
	SQLDRIVER = "mysql"
)

//返回驱动
func Driver() (driver.DataDriver) {
	return &MysqlDriver{}
}