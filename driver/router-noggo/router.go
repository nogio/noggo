package router_noggo


import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo"
	"sync"
)



//路由器定义
type Router struct {
	mutex sync.Mutex
	routes map[string]string
}

func NewRouter() *Router {
	return &Router{
		map[string]string{},
	}
}








//注册路由
func (router *Router) Route(name, uri string) {
	router.routes[name] = uri
}



//解析路由
func (router *Router) Parse(name, uri string) (*noggo.RouterResult){
	return nil
}
