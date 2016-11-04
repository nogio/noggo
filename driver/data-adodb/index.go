package data_adodb

import (
	"github.com/nogio/noggo"
	_ "github.com/nogio/noggo/depend/go-adodb"   //此包自动注册名为mysql的sql驱动
)

const (
	SQLDRIVER = "adodb"
)

//返回驱动
func Driver() (noggo.DataDriver) {
	return &AdodbDriver{}
}