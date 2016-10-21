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

				//noggo.Event.Publish("test")
				for i:=0;i<10;i++ {
					noggo.Queue.Publish("test")
				}

				ctx.Data["msg"] = "消息来自路由"
				ctx.View("index")
			},
		},
	})


	noggo.Http.Route("test", Map{
		"uri": "/test",
		"route": Map{
			"name": "test", "text": "test",
			"action": func(ctx *noggo.HttpContext) {
				ctx.Return(Map{
					"msg": "hahaha",
				})
			},
		},
	})

}