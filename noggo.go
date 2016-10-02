package noggo

import (
	. "github.com/nogio/noggo/base"
)

var (
	//编号
	//每一个noggo实例应当设置一个唯一Id，好在分布式环境中区分
	Id	string

	//常量模块
	Const *constModule
	//Map模块
	Mapping	*mappingModule

	//全局配置
	Config configConfig


	//路由器模块
	Router *routerModule
	//会话模块
	Session *sessionModule

	//日志模块
	Logger *loggerModule
	//触发器模块
	Trigger *triggerModule

	//计划模块
	Plan *planModule
	//计划模块
	Task *taskModule
)


func init() {
	Config = readJsonConfig()

	Const = &constModule{}

	Mapping = &mappingModule{
		types: map[string]Map{}, cryptos: map[string]Map{},
	}


	//路由器模块
	Router = &routerModule{
		drivers: map[string]RouterDriver{},
	}
	//会话模块
	Session = &sessionModule{
		drivers: map[string]SessionDriver{},
	}





	//日志模块
	Logger = &loggerModule{
		drivers: map[string]LoggerDriver{},
	}
	//触发器模块
	Trigger = &triggerModule{
		routes: map[string]Map{}, routeNames:[]string{}, routeUris: map[string]string{},
	}

	//计划模块
	Plan = &planModule{
		drivers: map[string]PlanDriver{},
		routes: map[string]Map{}, routeNames:[]string{}, routeUris: map[string]string{}, routeTimes:map[string][]string{},
	}
	//任务模块
	Task = &taskModule{
		drivers: map[string]TaskDriver{},
		routes: map[string]Map{}, routeNames:[]string{}, routeUris: map[string]string{}, routeTimes:map[string][]string{},
	}
}