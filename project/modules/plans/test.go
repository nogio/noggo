package https

import (
	"github.com/nogio/noggo"
	. "github.com/nogio/noggo/base"
	"time"
)

func init() {

	noggo.Plan.Route("test", Map{
		"time": "*/5 * * * * *",
		"route": Map{
			"name": "测试计划", "text": "测试计划",
			"action": func(ctx *noggo.PlanContext) {
				noggo.Logger.Debug(ctx.Node.Name, "测试计划开始了")

				//2秒后执行任务
				noggo.Task.Touch("test", time.Second*2)

				ctx.Finish()
			},
		},
	})


}