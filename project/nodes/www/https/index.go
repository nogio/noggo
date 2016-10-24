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
				ctx.View("index")
			},
		},
	})


	noggo.Http.Route("test", Map{
		"uri": "/testhahahaa",
		"route": Map{
			"name": "测试路由", "text": "测试路由", "coded": true,
			//"args": noggo.Data.Fields("main", "test"),
			"data": Map{
				"msg": Map{ "type": "string", "must": true, "name": "消息", "text": "消息" },
			},
			"state": noggo.Sugar.States("ok", "no"),
			"action": func(ctx *noggo.HttpContext) {
				ctx.Json(Map{ "msg": "hahaha", "url": ctx.Url.Route("test") })
			},
		},
	})

}