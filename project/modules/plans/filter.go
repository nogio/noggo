package plans

import (
	"github.com/nogio/noggo"
)

func init() {

	//请求拦截器
	noggo.Plan.RequestFilter("request", func(ctx *noggo.PlanContext) {
		ctx.Next()
	})
	//执行拦截器
	noggo.Plan.ExecuteFilter("execute", func(ctx *noggo.PlanContext) {
		ctx.Next()
	})
	//响应拦截器
	noggo.Plan.ResponseFilter("response", func(ctx *noggo.PlanContext) {
		ctx.Next()
	})
}