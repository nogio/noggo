package https

import (
	"github.com/nogio/noggo"
)

func init() {

	//请求日志拦截器
	noggo.Http.RequestFilter("request", func(ctx *noggo.HttpContext) {
		//请求开始前
		ctx.Next()
		//请求结束
	})

	//执行拦截器
	noggo.Http.ExecuteFilter("execute", func(ctx *noggo.HttpContext) {
		//这里是路由action之前
		ctx.Next()
		//这里是路由action之后
	})

	//响应拦截器
	noggo.Http.ResponseFilter("response", func(ctx *noggo.HttpContext) {
		//给客户端响应前
		ctx.Next()
		//响应之后
	})

}