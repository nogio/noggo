package https

import (
	"github.com/nogio/noggo"
	. "github.com/nogio/noggo/base"
)

func init() {

	noggo.Http.Route("cache", Map{
		"uri": "/cache",
		"route": Map{
			"name": "缓存", "text": "缓存",
			"action": func(ctx *noggo.HttpContext) {
				cb := noggo.Cache.Base("main"); defer cb.Close();
				val,_ := cb.Get("msg")

				ctx.Text(val)
			},
		},
	})


	noggo.Http.Route("cache.test", Map{
		"uri": "/cache/test",
		"route": Map{
			"name": "测试会话", "text": "测试会话", "coded": true,
			"args": Map{
				"msg": Map{ "type": "string", "must": true, "name": "消息", "text": "消息" },
			},
			"action": func(ctx *noggo.HttpContext) {

				cb := noggo.Cache.Base("main"); defer cb.Close();
				cb.Set("msg", ctx.Args["msg"])

				ctx.Text("ok")
			},
		},
	})

}