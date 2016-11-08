/*
	session 会话模块
	会话模块，是一个通用模块

	session运行方式：
	sessionGlobal会建立一个全局的会话连接，用于全局的Trigger, Task
	另外每个节点都会建立一个节点的会话连接，用于节点的Plan,Event,Queue,Http

*/


package noggo

import (
	. "github.com/nogio/noggo/base"
	"sync"
)


// session driver begin

type (
	//会话驱动
	SessionDriver interface {
		Connect(config Map) (SessionConnect,error)
	}
	//会话连接
	SessionConnect interface {
		//打开连接
		Open() error
		//关闭连接
		Close() error


		//查询会话
		Query(id string) (Map,error)
		//更新会话数据，不存在则创建，存在就更新
		Update(id string, value Map, exps ...int64) error
		//删除会话
		Remove(id string) error
	}
)


// session driver end



type (
	//会话全局
	sessionGlobal struct {
		mutex sync.Mutex
		drivers         map[string]SessionDriver

		sessionConfig   *sessionConfig
		sessionConnect  SessionConnect
	}
)








//注册会话驱动
func (global *sessionGlobal) Driver(name string, config SessionDriver) {
	global.mutex.Lock()
	defer global.mutex.Unlock()

	if global.drivers == nil {
		global.drivers = map[string]SessionDriver{}
	}

	if config == nil {
		panic("会话: 驱动不可为空")
	}
	global.drivers[name] = config
}


//连接驱动
func (global *sessionGlobal) connect(config *sessionConfig) (SessionConnect,error) {
	if sessionDriver,ok := global.drivers[config.Driver]; ok {
		return sessionDriver.Connect(config.Config)
	} else {
		panic("会话：不支持的驱动 " + config.Driver)
	}
}

//会话初始化
func (global *sessionGlobal) init() {
	global.initSession()
}

//初始化会话驱动
func (global *sessionGlobal) initSession() {

	//连接会话
	global.sessionConfig = Config.Session
	con,err := Session.connect(global.sessionConfig)

	if err != nil {
		panic("会话：连接失败 " + err.Error())
	} else {

		//打开会话连接
		err := con.Open()
		if err != nil {
			panic("会话：打开失败 " + err.Error())
		}

		//保存连接
		global.sessionConnect = con

	}
}
//会话退出
func (global *sessionGlobal) exit() {
	global.exitSession()
}
//任务退出，会话
func (global *sessionGlobal) exitSession() {
	if global.sessionConnect != nil {
		global.sessionConnect.Close()
	}
}





//-------------------------------------------------------------------------

type (
	//会话模块
	sessionModule struct {
		node	*Noggo
		sessionConfig   *sessionConfig
		sessionConnect  SessionConnect
	}
)





//会话模块初始化
func (module *sessionModule) run() {
	module.runSession()
}
func (module *sessionModule) runSession() {
	if module.node.Config.Session != nil {
		//使用节点的会话配置
		module.sessionConfig = module.node.Config.Session
	} else {
		//使用默认的会话配置
		module.sessionConfig = Config.Session
	}

	//连接会话
	con,err := Session.connect(module.sessionConfig)

	if err != nil {
		panic("节点会话：连接失败：" + err.Error())
	} else {

		//打开会话连接
		err := con.Open()
		if err != nil {
			panic("节点会话：打开失败 " + err.Error())
		}

		//保存连接
		module.sessionConnect = con
	}
}


//会话模块退出
func (module *sessionModule) end() {
	module.endSession()
}
//退出SESSION
func (module *sessionModule) endSession() {
	if module.sessionConnect != nil {
		module.sessionConnect.Close()
	}
}






//新建SESSION模块
func newSessionModule(node *Noggo) (*sessionModule) {
	return &sessionModule{ node: node }
}
