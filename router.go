package noggo

import (
	. "github.com/nogio/noggo/base"
	"sync"
)

/*
	router 路由器模块
*/


type (
	//路由器模块
	routerModule struct {
		drivers map[string]*SessionDriver
		driversLock sync.RWMutex

		driver *SessionDriver
	}

	//路由器结果
	RouterResult struct {
		Name	string
		Uri		string
		Params	Map
	}

	//路由器驱动
	RouterDriver interface {
		//注册路由
		Route(name, uri string)
		//解析路由
		Parse(uri string) *RouterResult
	}

)








//注册路由器驱动
func (router *routerModule) Register(name string, driver *SessionDriver) {
	router.driversLock.Lock()
	defer router.driversLock.Unlock()

	if driver == nil {
		panic("router: Register driver is nil")
	}
	if _, ok := router.drivers[name]; ok {
		panic("router: Registered driver " + name)
	}

	router.drivers[name] = driver
}