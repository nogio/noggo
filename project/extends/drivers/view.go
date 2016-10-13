package drivers


import (
	"github.com/nogio/noggo"
	"github.com/nogio/noggo/driver/view-default"
)


func init() {
	//注册默认View驱动
	noggo.View.Driver("default", view_default.Driver())
}
