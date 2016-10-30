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

		noggo.Logger.Debug(ctx.Name, "event.error", ctx.Wrong.Text)

		ctx.Finish()
	})
}