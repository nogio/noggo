package https

import (
	"github.com/nogio/noggo"
)

func init() {

	//请求拦截器
	noggo.Plan.RequestFilter("request", func(ctx *noggo.PlanContext) {
		noggo.Logger.Debug(ctx.Node.Name, "计划拦截器开始")
		ctx.Next()
	})


}