package https

import (
    "github.com/nogio/noggo"
)

func init() {

    //请求拦截器
    noggo.Http.RequestFilter("request", func(ctx *noggo.HttpContext) {
        ctx.Next()
    })
    //执行拦截器
    noggo.Http.ExecuteFilter("execute", func(ctx *noggo.HttpContext) {
        ctx.Next()
    })
    //响应拦截器
    noggo.Http.ResponseFilter("response", func(ctx *noggo.HttpContext) {
        ctx.Next()
    })
}