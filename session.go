package noggo

import (
	. "github.com/nogio/noggo/base"
	"sync"
)

/*
	session 触发器模块
*/

type (
	//会话模块
	sessionModule struct {
		drivers map[string]*SessionDriver
		driversLock sync.RWMutex
	}



	//会话驱动
	SessionDriver interface {
		//打开连接
		//如redis之类的使用之前打开连接
		Open()
		//关闭连接
		//如redis之类的使用之后关闭连接
		Close()

		//生成session唯一id方法
		Id() string
		//创建或查询会话
		Create(id string, expiry int64) Map
		//更新会话数据
		Update(id string, value Map, expiry int64) bool
		//删除会话
		Remove(id string) bool
		//回收会话，系统会每一段时间自动调用此方法
		Recycle(expiry int64) bool
	}

)








//注册会话驱动
func (session *sessionModule) Register(name string, driver *SessionDriver) {
	session.driversLock.Lock()
	defer session.driversLock.Unlock()

	if driver == nil {
		panic("session: Register driver is nil")
	}
	if _, ok := session.drivers[name]; ok {
		panic("session: Registered driver " + name)
	}

	session.drivers[name] = driver
}