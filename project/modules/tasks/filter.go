package tasks

import (
    "github.com/nogio/noggo"
)

func init() {

    //请求拦截器
    noggo.Task.RequestFilter("request", func(ctx *noggo.TaskContext) {
        ctx.Next()
    })
    //执行拦截器
    noggo.Task.ExecuteFilter("execute", func(ctx *noggo.TaskContext) {
        ctx.Next()
    })
    //响应拦截器
    noggo.Task.ResponseFilter("response", func(ctx *noggo.TaskContext) {
        ctx.Next()
    })
}