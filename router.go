package noggo

import (
	. "github.com/nogio/noggo/base"
	"sync"
)

/*
	router 路由器模块
*/


type (
	//路由器结果
	RouterResult struct {
		Name	string
		Uri		string
		Param	Map
	}

	//路由器驱动
	RouterDriver interface {
		Connect(config Map) (RouterConnect)
	}
	//路由器连接
	RouterConnect interface {
		//打开连接
		Open() error
		//关闭连接
		Close()
		//注册路由
		Route(name, uri string)
		//解析路由
		Parse(host, path string) *RouterResult
	}
	//路由器模块
	routerModule struct {
		drivers map[string]RouterDriver
		driversMutex sync.Mutex

		//这个是默认的路由器
		routerConfig *routerConfig
		routerConnect RouterConnect
	}
)






//注册路由器驱动
func (router *routerModule) Driver(name string, driver RouterDriver) {
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



//连接驱动
func (router *routerModule) connect(config *routerConfig) (RouterConnect) {
	if routerDriver,ok := router.drivers[config.Driver]; ok {
		return routerDriver.Connect(config.Config)
	}
	return nil
}

//路由器初始化
func (router *routerModule) init() {

	//默认配置
	router.routerConfig = Config.Router
	router.routerConnect = Router.connect(router.routerConfig)

	err := router.routerConnect.Open()
	if err != nil {
		panic("打开路由器连接失败")
	}
}
//路由器退出
func (router *routerModule) exit() {
	//关闭路由器连接
	if router.routerConnect != nil {
		router.routerConnect.Close()
	}
}


