package http_default

import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo/driver"
	"net/http"
)


type (
	//驱动
	DefaultHttpDriver struct {
	}
	//连接
	DefaultHttpConnect struct {
		config Map
		execute driver.HttpAcceptCall    //func(res http.ResponseWriter, req *http.Request)
		server *http.Server
	}
)


//返回驱动
func Driver() *DefaultHttpDriver {
	return &DefaultHttpDriver{}
}





//连接
func (driver *DefaultHttpDriver) Connect(config Map) (driver.HttpConnect) {
	//新建连接
	return &DefaultHttpConnect{
		config: config,
	}
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






//监听
func (connect *DefaultHttpConnect) Accept(execute driver.HttpAcceptCall) error {
	connect.execute = execute
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
	if connect.execute == nil {
		panic("未监听http请求")
	}
	connect.execute(res, req)
}

