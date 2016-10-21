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

		//数据连接
		connects    map[string]driver.DataConnect

		models      map[string]map[string]Map
		views      map[string]map[string]Map
	}
)


//连接驱动
func (global *dataGlobal) connect(config *dataConfig) (driver.DataConnect,error) {
	if dataDriver,ok := global.drivers[config.Driver]; ok {
		return dataDriver.Connect(config.Config)
	} else {
		panic("数据：不支持的驱动 " + config.Driver)
	}
}


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










//数据初始化
func (global *dataGlobal) init() {
	global.initData()
}
func (global *dataGlobal) initData() {

	//遍历数据配置
	for name,config := range Config.Data {
		conn,err := global.connect(config)
		if err != nil {
			panic("数据：连接失败：" + err.Error())
		} else {
			err := conn.Open()
			if err != nil {
				panic("数据：打开连接失败：" + err.Error())
			} else {

				//注册模型
				//先全局
				for k,v := range global.models[ConstNodeGlobal] {
					conn.Model(k, v)
				}
				//再库
				for k,v := range global.models[name] {
					conn.Model(k, v)
				}


				//注册讽图
				//先全局
				for k,v := range global.views[ConstNodeGlobal] {
					conn.View(k, v)
				}
				//再库
				for k,v := range global.views[name] {
					conn.View(k, v)
				}




				//保存连接
				global.connects[name] = conn

			}
		}
	}
}

//数据退出
func (global *dataGlobal) exit() {
	global.exitData()
}
func (global *dataGlobal) exitData() {
	for _,conn := range global.connects {
		conn.Close()
	}
}
























//注册模型
func (global *dataGlobal) Model(name string, config Map) {
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



//注册视图
func (global *dataGlobal) View(name string, config Map) {
	global.mutex.Lock()
	defer global.mutex.Unlock()

	if global.views == nil {
		global.views = map[string]map[string]Map{}
	}

	//节点
	nodeName := ConstNodeGlobal
	if Current != "" {
		nodeName = Current
	}

	//如果节点配置不存在，创建
	if global.views[nodeName] == nil {
		global.views[nodeName] = map[string]Map{}
	}


	//可以后注册重写原有配置，所以直接保存
	global.views[nodeName][name] = config
}












//返回DB对象
func (global *dataGlobal) Base(name string) (driver.DataBase) {
	if conn,ok := global.connects[name]; ok {
		db,err := conn.Base(name)
		if err != nil {
			panic("数据：打开DB失败：" + err.Error())
		} else {
			//返回
			return db
		}
	} else {
		panic("数据：未定义的数据库")
	}
}