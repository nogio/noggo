package middlers

import (
	"github.com/nogio/noggo"
	//"github.com/nogio/noggo/middler/plan-logger"
)

func init() {

	noggo.Trigger.RequestFilter("asdfasdf", func(ctx *noggo.TriggerContext) {
		noggo.Logger.Debug("Trigger.RequestFilter.begin")
		ctx.Next()
		noggo.Logger.Debug("Trigger.RequestFilter.end")
	})

	noggo.Trigger.ExecuteFilter("asdfasdf", func(ctx *noggo.TriggerContext) {
		noggo.Logger.Debug("Trigger.ExecuteFilter.begin")
		ctx.Next()
		noggo.Logger.Debug("Trigger.ExecuteFilter.end")
	})

	noggo.Trigger.ResponseFilter("asdfasdf", func(ctx *noggo.TriggerContext) {
		noggo.Logger.Debug("Trigger.ResponseFilter.begin")
		ctx.Next()
		noggo.Logger.Debug("Trigger.ResponseFilter.end")
	})


	noggo.Trigger.FoundHandler("asdfasdf", func(ctx *noggo.TriggerContext) {
		noggo.Logger.Debug("Trigger.FoundHandler.begin")
		ctx.Next()
		noggo.Logger.Debug("Trigger.FoundHandler.end")
	})


}
