package drivers


import (
	"github.com/nogio/noggo"
	"github.com/nogio/noggo/driver/router-noggo"
)


func init() {
	//注册默认路由器驱动
	noggo.Router.Register("noggo", router_noggo.Driver())
}
