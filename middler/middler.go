package middler

import (
	"github.com/nogio/noggo"
	"github.com/nogio/noggo/middler/http-logger"
	"github.com/nogio/noggo/middler/http-static"
)

//请求日志中间件
func HttpLogger() (noggo.HttpFunc) {
	return http_logger.Middler()
}

//静态文件中间件
func HttpStatic(paths ...string) (noggo.HttpFunc) {
	return http_static.Middler(paths...)
}
