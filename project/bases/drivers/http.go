package drivers

import (
	"github.com/nogio/noggo"
	"github.com/nogio/noggo/driver/http-default"
)

func init() {
	//注册默认HTTP驱动
	noggo.Http.Driver("default", http_default.Driver())
}
