package middlers

import (
	"github.com/nogio/noggo"
	"github.com/nogio/noggo/middler/http-static"
	"github.com/nogio/noggo/middler/http-logger"
)


func init() {

	//请求日志中间件
	noggo.Http.Middler("logger", http_logger.Middler())

	//注册HTTP静态文件中间件
	noggo.Http.Middler("static", http_static.Middler())

}
