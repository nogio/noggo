package noggo


import (
	. "github.com/nogio/noggo/base"
	"sync"
)


type (

	//路由器驱动
	LoggerDriver interface {
		Connect(config Map) (LoggerConnect)
	}
	//路由器连接
	LoggerConnect interface {
		//打开连接
		Open() error
		//关闭连接
		Close()

		//输出调试
		Debug(args ...interface{})
		//输出信息
		Info(args ...interface{})
		//输出错误
		Error(args ...interface{})
	}

	//日志模块
	loggerModule struct {
		drivers map[string]LoggerDriver
		driversMutex sync.Mutex

		loggerConfig *loggerConfig
		loggerConnect LoggerConnect
	}
)






//注册日志器驱动
func (logger *loggerModule) Register(name string, driver LoggerDriver) {
	logger.driversMutex.Lock()
	defer logger.driversMutex.Unlock()

	if driver == nil {
		panic("logger: Register driver is nil")
	}
	if _, ok := logger.drivers[name]; ok {
		panic("logger: Registered driver " + name)
	}

	logger.drivers[name] = driver
}



//连接驱动
func (logger *loggerModule) connect(config *loggerConfig) (LoggerConnect) {

	if loggerDriver,ok := logger.drivers[config.Driver]; ok {
		return loggerDriver.Connect(config.Config)
	}
	return nil
}

//日志初始化
func (logger *loggerModule) init() {

	//先拿到默认的配置
	logger.loggerConfig = Config.Logger
	logger.loggerConnect = Logger.connect(logger.loggerConfig)

	err := logger.loggerConnect.Open()
	if err != nil {
		panic("打开日志连接失败")
	}

}
//日志退出
func (logger *loggerModule) exit() {
	//关闭日志连接
	if logger.loggerConnect != nil {
		logger.loggerConnect.Close()
	}
}






//调试
func (logger *loggerModule) Debug(args ...interface{}) {
	if logger.loggerConnect != nil && Config.Debug {
		logger.loggerConnect.Debug(args...)
	}
}
//信息
func (logger *loggerModule) Info(args ...interface{}) {
	if logger.loggerConnect != nil {
		logger.loggerConnect.Info(args...)
	}
}
//错误
func (logger *loggerModule) Error(args ...interface{}) {
	if logger.loggerConnect != nil {
		logger.loggerConnect.Error(args...)
	}
}