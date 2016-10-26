package http_logger


import (
	"github.com/nogio/noggo"
	"time"
)



//返回中间件
func Middler() (noggo.HttpFunc) {
	return func(ctx *noggo.HttpContext) {
		begin := time.Now()
		ctx.Next()
		noggo.Logger.Info(ctx.Ip(), ctx.Id, ctx.Method, ctx.Node.Name, ctx.Path, ctx.Code, time.Now().Sub(begin))
	}
}

