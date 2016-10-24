package events

import (
	"github.com/nogio/noggo"
	. "github.com/nogio/noggo/base"
)

func init() {

	noggo.Event.Route("test", Map{
		"route": Map{
			"name": "测试事件", "text": "测试事件",
			"action": func(ctx *noggo.EventContext) {
				noggo.Logger.Debug(ctx.Node.Name, "event", "test", "开始了")
				ctx.Finish()
			},
		},
	})



}