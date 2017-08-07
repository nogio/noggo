/*
	view 视图模块
	视图模块，是一个全局模块
	用于解析HTTP页面View层
*/

package noggo


import (
	"sync"
	. "github.com/nogio/noggo/base"
)


// view driver begin


type (
	//视图驱动
	ViewDriver interface {
		Connect(config Map) (ViewConnect,error)
	}

	ViewParse struct {
		Node    string  //当前节点
		Lang    string  //当前语言
		Charset string  //当前字符集
		Data    Map     //ctx.Data
		Setting Map
		View    string  //视图文件
		Model   Any     //视图模型
		Args    Map     //参数
		Auth    Map     //登录对象信息
		Helpers Map     //工具方法
	}

	//视图连接
	ViewConnect interface {
		//打开连接
		Open() error
		//关闭连接
		Close() error

		//解析
		Parse(*ViewParse) (string,error)
	}
)


// view driver end



type (

	//视图模块
	viewGlobal struct {
		mutex sync.Mutex

		//视图驱动
		drivers map[string]ViewDriver

		//视图函数库
		//因为可能各种各样的函数，所以用Any类型
		helpers map[string]Any


		//日志配置，日志连接
		viewConfig *viewConfig
		viewConnect ViewConnect

	}
)



//连接驱动
func (global *viewGlobal) connect(config *viewConfig) (ViewConnect,error) {
	if viewDriver,ok := global.drivers[config.Driver]; ok {
		return viewDriver.Connect(config.Config)
	} else {
		panic("视图：不支持的驱动 " + config.Driver)
	}
}

//注册视图驱动
func (global *viewGlobal) Driver(name string, config ViewDriver) {
	global.mutex.Lock()
	defer global.mutex.Unlock()

	if config == nil {
		panic("视图: 驱动不可为空")
	}
	//不做存在判断，因为要支持后注册的驱动替换已注册的驱动
	//框架有可能自带几种默认驱动，并且是默认注册的，用户可以自行注册替换
	global.drivers[name] = config
}
func (global *viewGlobal) Helper(name string, helper Any) {
	global.mutex.Lock()
	defer global.mutex.Unlock()


	if global.helpers == nil {
		global.helpers = map[string]Any{}
	}

	if helper == nil {
		panic("视图: 方法不可为空")
	}
	//不做存在判断，因为要支持后注册的驱动替换已注册的驱动
	//框架有可能自带几种默认驱动，并且是默认注册的，用户可以自行注册替换
	global.helpers[name] = helper
}






//视图初始化
func (global *viewGlobal) init() {
	//全局视图不需要做什么
}
//视图退出
func (global *viewGlobal) exit() {
	//全局视图不需要做什么
}
