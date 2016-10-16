/*
	日志模块
	日志模块是一个全局模块，不属于任何一个节点

	主要功能是用于输出各种日志
*/

package noggo


import (
	//. "github.com/nogio/noggo/base"
	"sync"
	"github.com/nogio/noggo/driver"
)


type (

	/*
	//日志驱动
	LoggerDriver interface {
		Connect(config Map) (LoggerConnect)
	}
	//日志连接
	LoggerConnect interface {
		//打开连接
		Open() error
		//关闭连接
		Close() error

		//输出调试
		Debug(args ...interface{})
		//输出信息
		Info(args ...interface{})
		//输出错误
		Error(args ...interface{})
	}
	*/


	//日志模块
	loggerGlobal struct {
		mutex sync.Mutex

		//日志驱动容器
		drivers map[string]driver.LoggerDriver

		//日志配置，日志连接
		loggerConfig *loggerConfig
		loggerConnect driver.LoggerConnect
	}
)






//注册日志驱动
func (global *loggerGlobal) Driver(name string, config driver.LoggerDriver) {
	global.mutex.Lock()
	defer global.mutex.Unlock()

	if config == nil {
		panic("日志: 驱动不可为空")
	}
	//不做存在判断，因为要支持后注册的驱动替换已注册的驱动
	//框架有可能自带几种默认驱动，并且是默认注册的，用户可以自行注册替换
	global.drivers[name] = config
}


//连接驱动
func (global *loggerGlobal) connect(config *loggerConfig) (error,driver.LoggerConnect) {
	if loggerDriver,ok := global.drivers[config.Driver]; ok {
		return loggerDriver.Connect(config.Config)
	} else {
		panic("日志：不支持的驱动 " + config.Driver)
	}
}

//日志初始化
func (global *loggerGlobal) init() {

	//先拿到默认的配置
	global.loggerConfig = Config.Logger
	err,con := global.connect(global.loggerConfig)

	if err != nil {
		panic("日志：连接失败：" + err.Error())
	} else {

		global.loggerConnect = con

		err := global.loggerConnect.Open()
		if err != nil {
			panic("日志：打开失败 " + err.Error())
		}
	}


}
//日志退出
func (global *loggerGlobal) exit() {
	//关闭日志连接
	if global.loggerConnect != nil {
		global.loggerConnect.Close()
		global.loggerConnect = nil
	}
}









//调试
func (global *loggerGlobal) Debug(args ...interface{}) {
	if global.loggerConnect != nil && Config.Debug {
		global.loggerConnect.Debug(args...)
	}
}
//信息
func (global *loggerGlobal) Info(args ...interface{}) {
	if global.loggerConnect != nil {
		global.loggerConnect.Info(args...)
	}
}
//错误
func (global *loggerGlobal) Error(args ...interface{}) {
	if global.loggerConnect != nil {
		global.loggerConnect.Error(args...)
	}
}