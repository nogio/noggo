package tasks

import (
	"github.com/nogio/noggo"
)

func init() {

	//找不到处理
	noggo.Task.FoundHandler("found", func(ctx *noggo.TaskContext) {
		ctx.Finish()
	})

	//错误处理
	noggo.Task.ErrorHandler("error", func(ctx *noggo.TaskContext) {
		ctx.Finish()
	})
}