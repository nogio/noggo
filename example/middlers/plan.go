package middlers

import (
	"github.com/nogio/noggo"
	//"github.com/nogio/noggo/middler/plan-logger"
)

func init() {

	noggo.Plan.RequestFilter("asdfasdf", func(ctx *noggo.PlanContext) {
		noggo.Logger.Debug("plan.RequestFilter.begin")
		ctx.Next()
		noggo.Logger.Debug("plan.RequestFilter.end")
	})

	noggo.Plan.ExecuteFilter("asdfasdf", func(ctx *noggo.PlanContext) {
		noggo.Logger.Debug("plan.ExecuteFilter.begin")
		ctx.Next()
		noggo.Logger.Debug("plan.ExecuteFilter.end")
	})

	noggo.Plan.ResponseFilter("asdfasdf", func(ctx *noggo.PlanContext) {
		noggo.Logger.Debug("plan.ResponseFilter.begin")
		ctx.Next()
		noggo.Logger.Debug("plan.ResponseFilter.end")
	})

}
