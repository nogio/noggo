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
}