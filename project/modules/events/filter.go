package events

import (
	"github.com/nogio/noggo"
)

func init() {

	//请求拦截器
	noggo.Event.RequestFilter("request", func(ctx *noggo.EventContext) {
		ctx.Next()
	})
	//执行拦截器
	noggo.Event.ExecuteFilter("execute", func(ctx *noggo.EventContext) {
		ctx.Next()
	})
	//响应拦截器
	noggo.Event.ResponseFilter("response", func(ctx *noggo.EventContext) {
		ctx.Next()
	})
}