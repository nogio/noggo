package plans

import (
	"github.com/nogio/noggo"
)

func init() {

	//找不到处理
	noggo.Plan.FoundHandler("found", func(ctx *noggo.PlanContext) {
		ctx.Finish()
	})

	//错误处理
	noggo.Plan.ErrorHandler("error", func(ctx *noggo.PlanContext) {
		ctx.Finish()
	})
}