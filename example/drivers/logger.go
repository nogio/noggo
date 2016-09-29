package drivers


import (
	"github.com/nogio/noggo"
	"github.com/nogio/noggo/driver/logger-default"
)


func init() {
	//注册默认日志驱动
	noggo.Logger.Register("default", logger_default.Driver())
}
