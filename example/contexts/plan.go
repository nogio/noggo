package contexts

import (
	"github.com/nogio/noggo"
)


type (
	DefaultPlanContext struct {}
)


//请求拦截器，在请求一开始执行
func (context *DefaultPlanContext) RequestFilter(ctx *noggo.PlanCtx){
	noggo.Logger.Debug("plan.request.filter")
	ctx.Next()
}

//响应拦截器，在响应开始前执行
func (context *DefaultPlanContext) ResponseFilter(ctx *noggo.PlanCtx){
	noggo.Logger.Debug("plan.response.filter")
	ctx.Next()
}

//执行拦截器，在执行action前执行
func (context *DefaultPlanContext) ExecuteFilter(ctx *noggo.PlanCtx){
	noggo.Logger.Debug("plan.execute.filter")
	ctx.Next()
}




//404处理器，找不到请求时执行
func (context *DefaultPlanContext) FoundHandler(ctx *noggo.PlanCtx){
	ctx.Next()
}


//错误处理器，发生错误时执行
func (context *DefaultPlanContext) ErrorHandler(ctx *noggo.PlanCtx){
	ctx.Next()
}


//失败处理器，失败时执行，如参数解析失败
func (context *DefaultPlanContext) FailedHandler(ctx *noggo.PlanCtx){
	ctx.Next()
}


//拒绝处理器，拒绝时执行，主要用于Sign签名认证
func (context *DefaultPlanContext) DeniedHandler(ctx *noggo.PlanCtx){
	ctx.Next()
}





func init() {
	//注册默认上下文
	noggo.Plan.Context("default", &DefaultPlanContext{})
}
