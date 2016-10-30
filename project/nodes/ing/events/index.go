package events

import (
	"github.com/nogio/noggo"
	. "github.com/nogio/noggo/base"
)

func init() {

	noggo.Event.Route("test.a", Map{
		"route": Map{
			"name": "测试事件", "text": "测试事件",
			"args": Map{
				"msg": Map{
					"type": "string", "must": true, "name": "消息", "text": "消息",
				},
			},
			"action": func(ctx *noggo.PlanContext) {
				noggo.Logger.Debug(ctx.Node.Name, "test.a", ctx.Args)
				ctx.Finish()
			},
		},
	})
	noggo.Event.Route("test.b", Map{
		"route": Map{
			"name": "测试事件", "text": "测试事件",
			"args": Map{
				"id": Map{
					"type": "int", "must": true, "name": "编号", "text": "编号",
				},
			},
			"action": func(ctx *noggo.PlanContext) {
				noggo.Logger.Debug(ctx.Node.Name, "test.b", ctx.Args)
				ctx.Finish()
			},
		},
	})

}