package http_logger

import (
	"github.com/nogio/noggo"
	"time"
)

func Logger() (noggo.HttpFunc) {
	return func(ctx *noggo.HttpContext) {
		begin := time.Now()
		ctx.Next()
		end := time.Now()
		noggo.Logger.Info(ctx.Ip(), ctx.Id, ctx.Method, ctx.Path, ctx.Code, end.Sub(begin))
	}
}