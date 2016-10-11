package plans

import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo"
)

func init() {

	noggo.Task.Route("test", Map{
		"route": Map{
			"name": "测试任务", "text": "测试任务",
			"action": func(ctx *noggo.TaskContext) {
				noggo.Logger.Debug("task", ctx.Name, ctx.Id, "is work")
				ctx.Finish()
			},
		},
	})

}
