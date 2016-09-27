package contexts

import (
	"github.com/nogio/noggo"
)


type (
	DefaultTriggerContext struct {}
)





//请求拦截器，在请求一开始执行
func (context *DefaultTriggerContext) RequestFilter(ctx *noggo.TriggerCtx){
	noggo.Logger.Debug("trigger.request.filter")
	ctx.Next()
}

//响应拦截器，在响应开始前执行
func (context *DefaultTriggerContext) ResponseFilter(ctx *noggo.TriggerCtx){
	noggo.Logger.Debug("trigger.response.filter")
	ctx.Next()
}

//执行拦截器，在执行action前执行
func (context *DefaultTriggerContext) ExecuteFilter(ctx *noggo.TriggerCtx){
	noggo.Logger.Debug("trigger.execute.filter")
	ctx.Next()
}




//404处理器，找不到请求时执行
func (context *DefaultTriggerContext) FoundHandler(ctx *noggo.TriggerCtx){
	ctx.Next()
}


//错误处理器，发生错误时执行
func (context *DefaultTriggerContext) ErrorHandler(ctx *noggo.TriggerCtx){
	ctx.Next()
}


//失败处理器，失败时执行，如参数解析失败
func (context *DefaultTriggerContext) FailedHandler(ctx *noggo.TriggerCtx){
	ctx.Next()
}


//拒绝处理器，拒绝时执行，主要用于Sign签名认证
func (context *DefaultTriggerContext) DeniedHandler(ctx *noggo.TriggerCtx){
	ctx.Next()
}





func init() {
	//注册默认上下文
	noggo.Trigger.Context("default", &DefaultTriggerContext{})
}
