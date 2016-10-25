package plans

import (
	"github.com/nogio/noggo"
)

func init() {

	//找不到处理
	noggo.Plan.FoundHandler("found", func(ctx *noggo.PlanContext) {

	})

	//错误处理
	noggo.Plan.ErrorHandler("error", func(ctx *noggo.PlanContext) {

	})

	//失败处理
	noggo.Plan.FailedHandler("failed", func(ctx *noggo.PlanContext) {

	})

	//拒绝处理
	noggo.Plan.DeniedHandler("denied", func(ctx *noggo.PlanContext) {

	})
}