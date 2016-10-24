package https

import (
	"github.com/nogio/noggo"
)

func init() {

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
		ctx.Text("http failed " + ctx.Wrong.Error())
	})

	//拒绝处理
	noggo.Http.DeniedHandler("denied", func(ctx *noggo.HttpContext) {
		ctx.Text("http denied")
	})
}