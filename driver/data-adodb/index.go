package data_adodb

import (
	"github.com/nogio/noggo/driver"
	_ "github.com/mattn/go-adodb"   //此包自动注册名为mysql的sql驱动
)

const (
	SQLDRIVER = "adodb"
)

//返回驱动
func Driver() (driver.DataDriver) {
	return &AdodbDriver{}
}