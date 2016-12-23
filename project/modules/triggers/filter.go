package triggers

import (
    "github.com/nogio/noggo"
)

func init() {

    //请求拦截器
    noggo.Trigger.RequestFilter("request", func(ctx *noggo.TriggerContext) {
        ctx.Next()
    })
    //执行拦截器
    noggo.Trigger.ExecuteFilter("execute", func(ctx *noggo.TriggerContext) {
        ctx.Next()
    })
    //响应拦截器
    noggo.Trigger.ResponseFilter("response", func(ctx *noggo.TriggerContext) {
        ctx.Next()
    })
}