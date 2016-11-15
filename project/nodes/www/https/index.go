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
		"uri": "/test",
		"route": Map{
			"name": "测试路由", "text": "测试路由", "coded": true,
			//"args": noggo.Data.Fields("main", "test"),
			"data": Map{
				"msg": Map{ "type": "string", "must": true, "name": "消息", "text": "消息" },
			},
			"state": noggo.Const.StateStrings("ok", "no"),
			"action": func(ctx *noggo.HttpContext) {

				noggo.Event.Publish("test.a", Map{ "msg": "消息11111" })
				noggo.Event.Publish("test.a", Map{ "msg": "消息222" })
				noggo.Event.Publish("test.a", Map{ "msg": "消息3333" })
				noggo.Event.Publish("test.a", Map{ "msg": "消息4444" })

				noggo.Event.Publish("test.b", Map{ "id": 1 })
				noggo.Event.Publish("test.b", Map{ "id": 2 })
				noggo.Event.Publish("test.b", Map{ "id": 3 })
				noggo.Event.Publish("test.b", Map{ "id": 4 })


				ctx.Json(Map{ "msg": "hahaha", "url": ctx.Url.Route("test") })
			},
		},
	})

}