package noggo

import (
	//. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo/driver"
	"sync"
)

/*
	router 路由器模块
*/


type (
	//路由器模块
	routerModule struct {
		drivers map[string]driver.RouterDriver
		driversMutex sync.Mutex
	}
)








//注册路由器驱动
func (router *routerModule) Register(name string, driver driver.RouterDriver) {
	router.driversMutex.Lock()
	defer router.driversMutex.Unlock()

	if driver == nil {
		panic("router: Register driver is nil")
	}
	if _, ok := router.drivers[name]; ok {
		panic("router: Registered driver " + name)
	}

	router.drivers[name] = driver
}



//路由器初始货
func (session *routerModule) init() {

}
//路由器退出
func (session *routerModule) exit() {

}