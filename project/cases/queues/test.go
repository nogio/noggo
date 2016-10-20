package quques

import (
	"github.com/nogio/noggo"
	. "github.com/nogio/noggo/base"
)

func init() {

	noggo.Queue.Route("test", Map{
		"line": 2,
		"route": Map{
			"name": "测试队列", "text": "测试队列",
			"action": func(ctx *noggo.QueueContext) {
				noggo.Logger.Debug(ctx.Node.Name, "队列", "test", "开始了")
				ctx.Finish()
			},
		},
	})



}