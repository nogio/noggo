package triggers

import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo"
)

func init() {

	noggo.Trigger.Route("data.create", Map{
		"route": Map{
			"name": "数据创建触发器", "text": "数据创建触发器",
			"args": Map{
				"entity": Map{
					"type": "json", "must": true, "name": "对象", "text": "对象",
				},
			},
			"action": func(ctx *noggo.TriggerContext) {
				noggo.Logger.Info("触发器", "data.create", ctx.Args)
				ctx.Finish();
			},
		},
	})

	noggo.Trigger.Route("data.change", Map{
		"route": Map{
			"name": "数据修改触发器", "text": "数据修改触发器",
			"args": Map{
				"before": Map{
					"type": "json", "must": true, "name": "旧对象", "text": "旧对象",
				},
				"after": Map{
					"type": "json", "must": true, "name": "新对象", "text": "新对象",
				},
			},
			"action": func(ctx *noggo.TriggerContext) {
				noggo.Logger.Info("触发器", "data.change", ctx.Args)
				ctx.Finish();
			},
		},
	})

	noggo.Trigger.Route("data.remove", Map{
		"route": Map{
			"name": "数据删除触发器", "text": "数据删除触发器",
			"args": Map{
				"entity": Map{
					"type": "json", "must": true, "name": "对象", "text": "对象",
				},
			},
			"action": func(ctx *noggo.TriggerContext) {
				noggo.Logger.Info("触发器", "data.remove", ctx.Args)
				ctx.Finish();
			},
		},
	})

	noggo.Trigger.Route("data.recover", Map{
		"route": Map{
			"name": "数据恢复触发器", "text": "数据恢复触发器",
			"args": Map{
				"entity": Map{
					"type": "json", "must": true, "name": "对象", "text": "对象",
				},
			},
			"action": func(ctx *noggo.TriggerContext) {
				noggo.Logger.Info("触发器", "data.recover", ctx.Args)
				ctx.Finish();
			},
		},
	})
}
