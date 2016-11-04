package triggers

import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo"
)

func init() {

	noggo.Trigger.Route("test", Map{
		"route": Map{
			"name": "测试触发器", "text": "测试触发器说明",
			"args": Map{
				"id": Map{
					"type": "int", "must": true, "name": "编号", "text": "编号",
				},
			},
			"action": func(ctx *noggo.TriggerContext) {
				noggo.Logger.Debug("测试触发器工作了", ctx.Args)
				ctx.Finish();
			},
		},
	})

}
