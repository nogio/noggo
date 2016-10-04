package router_default


import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo"
	"sync"
)



//路由器定义
type (
	//驱动
	DefaultRouter struct {
	}
	DefaultConnect struct {
		config Map
		routes map[string]string
		routesMutex sync.Mutex
	}
)


//返回驱动
func Driver() *DefaultRouter {
	return &DefaultRouter{}
}




//打开路由器
func (router *DefaultRouter) Connect(config Map) (noggo.RouterConnect) {
	//新建路由器连接
	return &DefaultConnect{
		config: config, routes: map[string]string{},
	}
}














//打开连接
func (router *DefaultConnect) Open() error {
	return nil
}



//关闭连接
func (router *DefaultConnect) Close() {

}




//注册路由
func (router *DefaultConnect) Route(name, uri string) {
	router.routes[name] = uri
}



//解析路由
func (router *DefaultConnect) Parse(host,path string) (*noggo.RouterResult){
	return nil
}
