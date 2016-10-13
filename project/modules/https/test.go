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
				ctx.Text("这里是首页")
			},
		},
	})

	noggo.Http.Route("test", Map{
		"uri": "/test",
		"route": Map{
			"name": "测试路由", "text": "测试路由",
			"action": func(ctx *noggo.HttpContext) {
				ctx.Json(Map{ "msg": "hahaha" })
			},
		},
	})


	noggo.Http.Route("test.method", Map{
		"uri": "/test/method",
		"route": Map{
			"get": Map{
				"name": "get测试路由", "text": "测试路由",
				"action": func(ctx *noggo.HttpContext) {
					ctx.Json(Map{ "msg": "hahaha" })
				},
			},
			"post": Map{
				"name": "post测试路由", "text": "测试路由",
				"action": func(ctx *noggo.HttpContext) {
					ctx.Json(Map{ "msg": "hahaha" })
				},
			},
		},
	})

}