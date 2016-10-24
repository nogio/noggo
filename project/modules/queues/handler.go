package quques

import (
	"github.com/nogio/noggo"
)

func init() {

	//找不到处理
	noggo.Queue.FoundHandler("found", func(ctx *noggo.QueueContext) {
		ctx.Finish()
	})

	//错误处理
	noggo.Queue.ErrorHandler("error", func(ctx *noggo.QueueContext) {
		ctx.Finish()
	})

	//失败处理
	noggo.Queue.FailedHandler("failed", func(ctx *noggo.QueueContext) {
		ctx.Finish()
	})

	//拒绝处理
	noggo.Queue.DeniedHandler("denied", func(ctx *noggo.QueueContext) {
		ctx.Finish()
	})
}