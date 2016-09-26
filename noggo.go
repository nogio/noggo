package noggo

import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo/driver"
)

var (
	//编号
	//每一个noggo实例应当设置一个唯一Id，好在分布式环境中区分
	Id	string

	Config configConfig

	//路由器模块
	Router *routerModule
	//会话模块
	Session *sessionModule

	//触发器模块
	Trigger *triggerModule
)


func init() {
	Config = readJsonConfig()


	//路由器模块
	Router = &routerModule{
		drivers: map[string]driver.RouterDriver{},
	}
	//会话模块
	Session = &sessionModule{
		drivers: map[string]driver.SessionDriver{},
	}
	//触发器模块
	Trigger = &triggerModule{
		routes: map[string]Map{}, routeNames:[]string{}, routeUris: map[string]string{},
	}
}