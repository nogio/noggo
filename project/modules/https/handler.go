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
	noggo.Http.FailedHandler("error", func(ctx *noggo.HttpContext) {
		ctx.Text("http error")
	})
	//拒绝处理
	noggo.Http.DeniedHandler("denied", func(ctx *noggo.HttpContext) {
		ctx.Text("http denied")
	})
}