package noggo

import (
	. "github.com/nogio/noggo/base"
	"sync"
	"time"
)

/*
	session 会话模块
*/

type (

	//会话值
	SessionValue struct {
		Value	Map
		Expiry	time.Time
	}

	//会话驱动
	SessionDriver interface {
		Connect(config Map) (SessionConnect)
	}
	//会话连接
	SessionConnect interface {
		//打开连接
		Open() error
		//关闭连接
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

	//会话模块
	sessionModule struct {
		drivers map[string]SessionDriver
		driversMutex sync.Mutex

		//默认会话连接
		sessionConfig *sessionConfig
		sessionConnect SessionConnect
	}
)








//注册会话驱动
func (session *sessionModule) Register(name string, driver SessionDriver) {
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


//连接驱动
func (session *sessionModule) connect(config *sessionConfig) (SessionConnect) {
	if sessionDriver,ok := session.drivers[config.Driver]; ok {
		return sessionDriver.Connect(config.Config)
	}
	return nil
}

//会话初始化
func (session *sessionModule) init() {

	//默认配置
	session.sessionConfig = Config.Session
	session.sessionConnect = session.connect(session.sessionConfig)

	err := session.sessionConnect.Open()
	if err != nil {
		panic("打开会话连接失败")
	}
}
//会话退出
func (session *sessionModule) exit() {
	//关闭日志连接
	if session.sessionConnect != nil {
		session.sessionConnect.Close()
	}
}


