package noggo

import (
	. "github.com/nogio/noggo/base"
	"sync"
	"github.com/nogio/noggo/driver"
)

type (
	dataGlobal    struct {
		mutex       sync.Mutex
		drivers     map[string]driver.DataDriver

		models      map[string]map[string]Map
	}
)



//注册数据驱动
func (global *dataGlobal) Driver(name string, config driver.DataDriver) {
	global.mutex.Lock()
	defer global.mutex.Unlock()

	if config == nil {
		panic("数据: 驱动不可为空")
	}
	//不做存在判断，因为要支持后注册的驱动替换已注册的驱动
	//框架有可能自带几种默认驱动，并且是默认注册的，用户可以自行注册替换
	global.drivers[name] = config
}


//连接驱动
func (global *dataGlobal) connect(config *dataConfig) (error,driver.DataConnect) {
	if dataDriver,ok := global.drivers[config.Driver]; ok {
		return dataDriver.Connect(config.Config)
	} else {
		panic("数据：不支持的驱动 " + config.Driver)
	}
}




//注册模型
func (global *dataGlobal) Register(name string, config Map) {
	global.mutex.Lock()
	defer global.mutex.Unlock()

	if global.models == nil {
		global.models = map[string]map[string]Map{}
	}

	//节点
	nodeName := ConstNodeGlobal
	if Current != "" {
		nodeName = Current
	}

	//如果节点配置不存在，创建
	if global.models[nodeName] == nil {
		global.models[nodeName] = map[string]Map{}
	}


	//可以后注册重写原有路由配置，所以直接保存
	global.models[nodeName][name] = config
}



//返回DB对象
func (global *dataGlobal) DB(name string) (driver.DataDB) {
	return nil
}