package https

import (
	"github.com/nogio/noggo"
	. "github.com/nogio/noggo/base"
)

func init() {

	//找不到处理，404
	noggo.Http.FoundHandler("found", func(ctx *noggo.HttpContext) {
		ctx.View("found")
	})

	//错误处理，500
	noggo.Http.ErrorHandler("error", func(ctx *noggo.HttpContext) {
		ctx.View("error")
	})

	//失败处理，参数处理失败，等情况时
	noggo.Http.FailedHandler("failed", func(ctx *noggo.HttpContext) {
		ctx.Alert(ctx.Wrong.Text)
	})


	//拒绝处理，未登录拒绝访问时，
	noggo.Http.DeniedHandler("denied", func(ctx *noggo.HttpContext) {
		//跳转登录页
		ctx.Route("login",nil, Map{ "back": true })
	})
}