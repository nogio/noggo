package drivers


import (
	"github.com/nogio/noggo"
	"github.com/nogio/noggo/driver/session-default"
)


func init() {
	//注册默认会话驱动
	noggo.Session.Driver("default", session_default.Driver())
}
