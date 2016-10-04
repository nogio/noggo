package https

import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo"
)

func init() {


	noggo.Http.Route("index", Map{
		"uri": "/",
		"route": Map{
			"name": "测试HTTP", "text": "测试HTTP",
			"action": func(ctx *noggo.HttpContext) {
				ctx.Text("index")
			},
		},
	})


	noggo.Http.Route("view", Map{
		"uri": "/view/{id}",
		"route": Map{
			"name": "测试HTTP", "text": "测试HTTP",
			"action": func(ctx *noggo.HttpContext) {
				ctx.Json(ctx.Param)
			},
		},
	})


	noggo.Http.Route("json", Map{
		"uri": "/json",
		"route": Map{
			"get": Map{
				"name": "测试JSON", "text": "测试JSON",
				"action": func(ctx *noggo.HttpContext) {
					ctx.Text("here is get json")
				},
			},
			"post": Map{
				"name": "测试JSON", "text": "测试JSON",
				"action": func(ctx *noggo.HttpContext) {
					ctx.Json(Map{ "msg": "哈哈哈睦的是不是要工要" })
				},
			},
		},
	})

}
