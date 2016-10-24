package events

import (
	"github.com/nogio/noggo"
)

func init() {

	//找不到处理
	noggo.Event.FoundHandler("found", func(ctx *noggo.EventContext) {
		ctx.Finish()
	})

	//错误处理
	noggo.Event.ErrorHandler("error", func(ctx *noggo.EventContext) {
		ctx.Finish()
	})

	//失败处理
	noggo.Event.FailedHandler("failed", func(ctx *noggo.EventContext) {
		ctx.Finish()
	})

	//拒绝处理
	noggo.Event.DeniedHandler("denied", func(ctx *noggo.EventContext) {
		ctx.Finish()
	})
}