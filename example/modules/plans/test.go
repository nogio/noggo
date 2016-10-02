package plans

import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo"
)

func init() {

	noggo.Plan.Route("test", Map{
		"time": "*/5 * * * * *",
		"route": Map{
			"name": "测试计划", "text": "测试计划",
			"action": func(ctx *noggo.PlanContext) {
				noggo.Logger.Debug("测试计划工作了")
				noggo.Trigger.Touch("test.abcd.1234", Map{ "msg": "触发在计划中" })
				ctx.Finish();
			},
		},
	})

}
