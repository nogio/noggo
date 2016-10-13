package https

import (
	"github.com/nogio/noggo"
	"time"
)

func init() {

	//请求日志拦截器
	noggo.Http.RequestFilter("logger", func(ctx *noggo.HttpContext) {

		//这里是请求开始前

		begin := time.Now()
		ctx.Next()
		noggo.Logger.Info(ctx.Ip(), ctx.Id, ctx.Method, ctx.Path, ctx.Code, time.Now().Sub(begin))

		//这里是请求结束前
		//这里已经响应客户端
	})



	//执行拦截器
	noggo.Http.ExecuteFilter("execute", func(ctx *noggo.HttpContext) {
		//这里是路由action之前
		ctx.Next()
		//这里是路由action之后
	})



	//响应拦截器
	noggo.Http.ResponseFilter("response", func(ctx *noggo.HttpContext) {
		//在这里响应以前，您可以对响应结果进行修改，比如压缩body
		//ctx.Body 就可以操作啥啥啥啥的
		ctx.Next()
		//处理完记得一定要ctx.Next()
		//否则不会调用响应函数， 就不会发送回客户端了
	})


}