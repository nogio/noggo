package https

import (
	"github.com/nogio/noggo"
	. "github.com/nogio/noggo/base"
)

func init() {

	noggo.Http.Route("session", Map{
		"uri": "/session",
		"route": Map{
			"name": "会话", "text": "会话",
			"action": func(ctx *noggo.HttpContext) {
				ctx.Json(ctx.Session)
			},
		},
	})


	noggo.Http.Route("session.test", Map{
		"uri": "/session/test",
		"route": Map{
			"name": "测试会话", "text": "测试会话", "coded": true,
			"args": Map{
				"msg": Map{ "type": "string", "must": true, "name": "消息", "text": "消息" },
			},
			"action": func(ctx *noggo.HttpContext) {
				ctx.Session["msg"] = ctx.Args["msg"];
				ctx.Text("ok")
			},
		},
	})

}