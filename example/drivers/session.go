package drivers


import (
	"github.com/nogio/noggo"
	"github.com/nogio/noggo/driver/session-memory"
)


func init() {
	//注册默认会话驱动
	noggo.Session.Register("memory", session_memory.Driver())
}
