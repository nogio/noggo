package data_oracle

import (
	_ "github.com/nogio/noggo/depend/go-oci8"   //此包自动注册名为oracle的sql驱动
	"github.com/nogio/noggo"
)


const (
	SQLDRIVER = "oci8"

	TriggerCreate = "data.create"
	TriggerChange = "data.change"
	TriggerRemove = "data.remove"
	TriggerRecover = "data.recover"
)


//返回驱动
func Driver() (noggo.DataDriver) {
	return &OracleDriver{}
}