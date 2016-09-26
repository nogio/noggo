package noggo

import (
	//. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo/driver"
	"sync"
)

/*
	session 会话模块
*/

type (
	//会话模块
	sessionModule struct {
		drivers map[string]driver.SessionDriver
		driversMutex sync.Mutex
	}
)








//注册会话驱动
func (session *sessionModule) Register(name string, driver driver.SessionDriver) {
	session.driversMutex.Lock()
	defer session.driversMutex.Unlock()

	if driver == nil {
		panic("session: Register driver is nil")
	}
	if _, ok := session.drivers[name]; ok {
		panic("session: Registered driver " + name)
	}

	session.drivers[name] = driver
}

//会话初始化
func (session *sessionModule) init() {

}
//会话退出
func (session *sessionModule) exit() {

}




//连接会话驱动
func (session *sessionModule) connect(config *sessionConfig) (driver.SessionConnect) {
	if sessionDriver,ok := Session.drivers[config.Driver]; ok {
		return sessionDriver.Connect(config.Config)
	}
	return nil
}