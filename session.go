/*
	session 会话模块
	会话模块，是一个通用模块
*/


package noggo

import (
	//. "github.com/nogio/noggo/base"
	"sync"
	"github.com/nogio/noggo/driver"
)

type (

	/*
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
		Close() error


		//查询会话，不存在就创建新的返回
		Query(id string, expiry int64) (error,Map)
		//更新会话数据，不存在则创建，存在就更新
		Update(id string, value Map, expiry int64) error
		//删除会话
		Remove(id string) error
		//回收会话，系统会每一段时间自动调用此方法
		Recycle(expiry int64) error
	}
	*/

	//会话模块
	sessionGlobal struct {
		mutex sync.Mutex

		//会话驱动们
		drivers map[string]driver.SessionDriver
	}
)








//注册会话驱动
func (global *sessionGlobal) Driver(name string, config driver.SessionDriver) {
	global.mutex.Lock()
	defer global.mutex.Unlock()

	if global.drivers == nil {
		global.drivers = map[string]driver.SessionDriver{}
	}

	if config == nil {
		panic("会话: 驱动不可为空")
	}
	global.drivers[name] = config
}


//连接驱动
func (global *sessionGlobal) connect(config *sessionConfig) (driver.SessionConnect) {
	if sessionDriver,ok := global.drivers[config.Driver]; ok {
		return sessionDriver.Connect(config.Config)
	} else {
		panic("会话：不支持的驱动 " + config.Driver)
	}
}

//会话初始化
func (global *sessionGlobal) init() {
	//会话全局容器，不需要处理任何东西
}

//会话退出
func (global *sessionGlobal) exit() {
	//会话全局容器，不需要处理任何东西
}


