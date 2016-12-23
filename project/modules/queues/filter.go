package quques

import (
    "github.com/nogio/noggo"
)

func init() {

    //请求拦截器
    noggo.Queue.RequestFilter("request", func(ctx *noggo.QueueContext) {
        ctx.Next()
    })
    //执行拦截器
    noggo.Queue.ExecuteFilter("execute", func(ctx *noggo.QueueContext) {
        ctx.Next()
    })
    //响应拦截器
    noggo.Queue.ResponseFilter("response", func(ctx *noggo.QueueContext) {
        ctx.Next()
    })
}