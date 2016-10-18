package https

import (
	"github.com/nogio/noggo"
	. "github.com/nogio/noggo/base"
)

func init() {

	noggo.Http.Route("index", Map{
		"uri": "/",
		"route": Map{
			"name": "首页", "text": "首页",
			"action": func(ctx *noggo.HttpContext) {
				ctx.Data["msg"] = "消息来自路由"
				ctx.View("index")
			},
		},
	})

}