package triggers

import (
	"github.com/nogio/noggo"
)

func init() {

	//找不到处理
	noggo.Trigger.FoundHandler("found", func(ctx *noggo.TriggerContext) {
		ctx.Finish()
	})

	//错误处理
	noggo.Trigger.ErrorHandler("error", func(ctx *noggo.TriggerContext) {
		ctx.Finish()
	})
}