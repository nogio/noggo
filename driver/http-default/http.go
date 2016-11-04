package http_default

import (
	. "github.com/nogio/noggo/base"
	"net/http"
	"github.com/nogio/noggo"
)


type (
	//驱动
	DefaultHttpDriver struct {
	}
	//连接
	DefaultHttpConnect struct {
		config Map
		handler noggo.HttpHandler    //func(req *http.Request, res http.ResponseWriter)
		server *http.Server
	}
)


//返回驱动
func Driver() (noggo.HttpDriver) {
	return &DefaultHttpDriver{}
}





//连接
func (drv *DefaultHttpDriver) Connect(config Map) (noggo.HttpConnect,error) {
	return &DefaultHttpConnect{
		config: config,
	},nil
}












//打开连接
func (connect *DefaultHttpConnect) Open() error {
	connect.server = &http.Server{ Handler: connect }
	return nil
}
//关闭连接
func (connect *DefaultHttpConnect) Close() error {
	return nil
}






//Start 应该不要阻塞线程
func (connect *DefaultHttpConnect) Accept(handler noggo.HttpHandler) error {
	connect.handler = handler
	return nil
}


//Start 应该不要阻塞线程
func (connect *DefaultHttpConnect) Start(addr string) error {
	if connect.server == nil {
		panic("请先初始化http server")
	}

	connect.server.Addr = addr
	go connect.server.ListenAndServe()
	return nil
}
func (connect *DefaultHttpConnect) StartTLS(addr string, certFile, keyFile string) error {
	if connect.server == nil {
		panic("请先初始化http server")
	}

	connect.server.Addr = addr
	go connect.server.ListenAndServeTLS(certFile, keyFile)
	return nil
}








//servehttp
func (connect *DefaultHttpConnect) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if connect.handler == nil {
		panic("未监听http请求")
	}
	connect.handler(req, res)
}

