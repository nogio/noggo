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
				ctx.Data["msg"] = "这个消息来自路由"
				ctx.Data["items"] = []Map{
					Map{ "id": 1, "name": "名字1" },
					Map{ "id": 2, "name": "名字2" },
					Map{ "id": 3, "name": "名字3" },
				}
				ctx.View("index/index")
			},
		},
	})

}