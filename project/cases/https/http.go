package https

import (
	"github.com/nogio/noggo"
	"github.com/nogio/noggo/middler"
)

func init() {

	//HTTP驱动自带，不用注册，如果使用三方驱动，需要注册

	//中间件
	noggo.Http.Middler("logger", middler.HttpLogger())
	noggo.Http.Middler("static", middler.HttpStatic())


	//找不到处理
	noggo.Http.FoundHandler("found", func(ctx *noggo.HttpContext) {
		ctx.View("404")
	})

	//错误处理
	noggo.Http.ErrorHandler("error", func(ctx *noggo.HttpContext) {
		ctx.Text("http error")
	})

	//失败处理
	noggo.Http.FailedHandler("failed", func(ctx *noggo.HttpContext) {
		ctx.Text("http failed")
	})

	//拒绝处理
	noggo.Http.DeniedHandler("denied", func(ctx *noggo.HttpContext) {
		ctx.Text("http denied")
	})
}