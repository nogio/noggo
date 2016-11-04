package data_mysql

import (
	"github.com/nogio/noggo"
	_ "github.com/nogio/noggo/depend/mysql"   //此包自动注册名为mysql的sql驱动
)

const (
	SQLDRIVER = "mysql"
)

//返回驱动
func Driver() (noggo.DataDriver) {
	return &MysqlDriver{}
}