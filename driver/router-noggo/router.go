package router_noggo


import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo/driver"
	"sync"
)



//路由器定义
type (
	//驱动
	NoggoRouter struct {
	}
	NoggoConnect struct {
		config Map
		routes map[string]string
		routesMutex sync.Mutex
	}
)


//返回驱动
func Driver() *NoggoRouter {
	return &NoggoRouter{}
}




//打开路由器
func (router *NoggoRouter) Connect(config Map) (driver.RouterConnect) {
	//新建路由器连接
	return &NoggoConnect{
		config: config, routes: map[string]string{},
	}
}














//打开连接
func (router *NoggoConnect) Open() {

}



//关闭连接
func (router *NoggoConnect) Close() {

}




//注册路由
func (router *NoggoConnect) Route(name, uri string) {
	router.routes[name] = uri
}



//解析路由
func (router *NoggoConnect) Parse(uri string) (*driver.RouterResult){
	return nil
}
