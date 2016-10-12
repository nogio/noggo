package noggo


import (
	"sync"
	. "github.com/nogio/noggo/base"
	"net/http"
)


type (
	//WEB函数
	WebFunc func(*WebContext)
	WebAcceptFunc func(res http.ResponseWriter, req *http.Request)



	//WEB驱动
	WebDriver interface {
		Connect(config Map) (WebConnect)
	}
	//WEB连接
	WebConnect interface {
		Open() error
		Close() error

		Accept(call WebAcceptFunc) error

		Start(addr string) error
		StartTLS(addr string, certFile, keyFile string) error
	}

	//WEB全局
	httpGlobal struct {
		mutex sync.Mutex

		//驱动
		drivers map[string]WebDriver

		//路由
		routes 		map[string]map[string]Map	//路由定义							map[node]map[name]Map
		routeNames	map[string][]string			//路由名称原始顺序，因为map是无序的		map[node][]string

		//拦截器们
		requestFilters, executeFilters, responseFilters map[string]map[string]WebFunc
		requestFilterNames, executeFilterNames, responseFilterNames map[string][]string

		//处理器们
		foundHandlers, errorHandlers, failedHandlers, deniedHandlers map[string]map[string]WebFunc
		foundHandlerNames, errorHandlerNames, failedHandlerNames, deniedHandlerNames map[string][]string
	}


	WebContext struct {

	}

)



//注册WEB驱动
func (global *httpGlobal) Driver(name string, driver WebDriver) {
	global.mutex.Lock()
	defer global.mutex.Unlock()

	if global.drivers == nil {
		global.drivers = map[string]WebDriver{}
	}

	if driver == nil {
		panic("WEB: 驱动不可为空")
	}
	//不做存在判断，因为要支持后注册的驱动替换已注册的驱动
	//框架有可能自带几种默认驱动，并且是默认注册的，用户可以自行注册替换
	global.drivers[name] = driver
}


//连接驱动
func (global *httpGlobal) connect(config *httpConfig) (WebConnect) {
	if httpDriver,ok := global.drivers[config.Driver]; ok {
		return httpDriver.Connect(config.Config)
	} else {
		panic("WEB：不支持的驱动 " + config.Driver)
	}
}




//WEB初始化
func (global *httpGlobal) init() {
	//全局的WEB不需要初始化处理
}




//WEB退出
func (global *httpGlobal) exit() {
	//全局的WEB不需要初始化处理
}








//WEB：注册路由
//注册路由到全局容器
//Current将标明表示属于哪一个节点
//如果Current为空，表示全局，相当于注册到所有节点
func (global *httpGlobal) Route(name string, config Map) {
	global.mutex.Lock()
	defer global.mutex.Unlock()

	if global.routes == nil {
		global.routes = map[string]map[string]Map{}
	}
	if global.routeNames == nil {
		global.routeNames = map[string][]string{}
	}


	//节点
	nodeName := ConstNodeGlobal
	if Current != "" {
		nodeName = Current
	}

	//如果节点配置不存在，创建
	if global.routes[nodeName] == nil {
		global.routes[nodeName] = map[string]Map{}
	}
	if global.routeNames[nodeName] == nil {
		global.routeNames[nodeName] = []string{}
	}


	//保存配置
	if _,ok := global.routes[nodeName][name]; ok == false {
		//没有注册过name，才把name加到列表
		global.routeNames[nodeName] = append(global.routeNames[nodeName], name)
	}
	//可以后注册重写原有路由配置，所以直接保存
	global.routes[nodeName][name] = config
}








/* 注册拦截器 begin */
//请求拦截器
func (global *httpGlobal) RequestFilter(name string, call WebFunc) {
	global.mutex.Lock()
	defer global.mutex.Unlock()

	if global.requestFilters == nil {
		global.requestFilters = map[string]map[string]WebFunc{}
	}
	if global.requestFilterNames == nil {
		global.requestFilterNames =  map[string][]string{}
	}


	//节点
	nodeName := ConstNodeGlobal
	if Current != "" {
		nodeName = Current
	}


	//如果节点配置不存在，创建
	if global.requestFilters[nodeName] == nil {
		global.requestFilters[nodeName] = map[string]WebFunc{}
	}
	if global.requestFilterNames[nodeName] == nil {
		global.requestFilterNames[nodeName] = []string{}
	}




	//如果没有注册个此name，才加入数组
	if _,ok := global.requestFilters[nodeName][name]; ok == false {
		global.requestFilterNames[nodeName] = append(global.requestFilterNames[nodeName], name)
	}
	//函数直接写， 因为可以使用同名替换现有的
	global.requestFilters[nodeName][name] = call
}
//执行拦截器
func (global *httpGlobal) ExecuteFilter(name string, call WebFunc) {
	global.mutex.Lock()
	defer global.mutex.Unlock()

	if global.executeFilters == nil {
		global.executeFilters = map[string]map[string]WebFunc{}
	}
	if global.executeFilterNames == nil {
		global.executeFilterNames =  map[string][]string{}
	}


	//节点
	nodeName := ConstNodeGlobal
	if Current != "" {
		nodeName = Current
	}


	//如果节点配置不存在，创建
	if global.executeFilters[nodeName] == nil {
		global.executeFilters[nodeName] = map[string]WebFunc{}
	}
	if global.executeFilterNames[nodeName] == nil {
		global.executeFilterNames[nodeName] = []string{}
	}




	//如果没有注册个此name，才加入数组
	if _,ok := global.executeFilters[nodeName][name]; ok == false {
		global.executeFilterNames[nodeName] = append(global.executeFilterNames[nodeName], name)
	}
	//函数直接写， 因为可以使用同名替换现有的
	global.executeFilters[nodeName][name] = call
}
//响应拦截器
func (global *httpGlobal) ResponseFilter(name string, call WebFunc) {
	global.mutex.Lock()
	defer global.mutex.Unlock()

	if global.responseFilters == nil {
		global.responseFilters = map[string]map[string]WebFunc{}
	}
	if global.responseFilterNames == nil {
		global.responseFilterNames =  map[string][]string{}
	}


	//节点
	nodeName := ConstNodeGlobal
	if Current != "" {
		nodeName = Current
	}


	//如果节点配置不存在，创建
	if global.responseFilters[nodeName] == nil {
		global.responseFilters[nodeName] = map[string]WebFunc{}
	}
	if global.responseFilterNames[nodeName] == nil {
		global.responseFilterNames[nodeName] = []string{}
	}




	//如果没有注册个此name，才加入数组
	if _,ok := global.responseFilters[nodeName][name]; ok == false {
		global.responseFilterNames[nodeName] = append(global.responseFilterNames[nodeName], name)
	}
	//函数直接写， 因为可以使用同名替换现有的
	global.responseFilters[nodeName][name] = call
}
/* 注册拦截器 end */


/* 注册处理器 begin */
//找不到处理
func (global *httpGlobal) FoundHandler(name string, call WebFunc) {
	global.mutex.Lock()
	defer global.mutex.Unlock()

	if global.foundHandlers == nil {
		global.foundHandlers = map[string]map[string]WebFunc{}
	}
	if global.foundHandlerNames == nil {
		global.foundHandlerNames =  map[string][]string{}
	}


	//节点
	nodeName := ConstNodeGlobal
	if Current != "" {
		nodeName = Current
	}


	//如果节点配置不存在，创建
	if global.foundHandlers[nodeName] == nil {
		global.foundHandlers[nodeName] = map[string]WebFunc{}
	}
	if global.foundHandlerNames[nodeName] == nil {
		global.foundHandlerNames[nodeName] = []string{}
	}




	//如果没有注册个此name，才加入数组
	if _,ok := global.foundHandlers[nodeName][name]; ok == false {
		global.foundHandlerNames[nodeName] = append(global.foundHandlerNames[nodeName], name)
	}
	//函数直接写， 因为可以使用同名替换现有的
	global.foundHandlers[nodeName][name] = call
}
//错误处理
func (global *httpGlobal) ErrorHandler(name string, call WebFunc) {
	global.mutex.Lock()
	defer global.mutex.Unlock()

	if global.errorHandlers == nil {
		global.errorHandlers = map[string]map[string]WebFunc{}
	}
	if global.errorHandlerNames == nil {
		global.errorHandlerNames =  map[string][]string{}
	}


	//节点
	nodeName := ConstNodeGlobal
	if Current != "" {
		nodeName = Current
	}


	//如果节点配置不存在，创建
	if global.errorHandlers[nodeName] == nil {
		global.errorHandlers[nodeName] = map[string]WebFunc{}
	}
	if global.errorHandlerNames[nodeName] == nil {
		global.errorHandlerNames[nodeName] = []string{}
	}




	//如果没有注册个此name，才加入数组
	if _,ok := global.errorHandlers[nodeName][name]; ok == false {
		global.errorHandlerNames[nodeName] = append(global.errorHandlerNames[nodeName], name)
	}
	//函数直接写， 因为可以使用同名替换现有的
	global.errorHandlers[nodeName][name] = call
}
func (global *httpGlobal) FailedHandler(name string, call WebFunc) {
	global.mutex.Lock()
	defer global.mutex.Unlock()

	if global.failedHandlers == nil {
		global.failedHandlers = map[string]map[string]WebFunc{}
	}
	if global.failedHandlerNames == nil {
		global.failedHandlerNames =  map[string][]string{}
	}


	//节点
	nodeName := ConstNodeGlobal
	if Current != "" {
		nodeName = Current
	}


	//如果节点配置不存在，创建
	if global.failedHandlers[nodeName] == nil {
		global.failedHandlers[nodeName] = map[string]WebFunc{}
	}
	if global.failedHandlerNames[nodeName] == nil {
		global.failedHandlerNames[nodeName] = []string{}
	}




	//如果没有注册个此name，才加入数组
	if _,ok := global.failedHandlers[nodeName][name]; ok == false {
		global.failedHandlerNames[nodeName] = append(global.failedHandlerNames[nodeName], name)
	}
	//函数直接写， 因为可以使用同名替换现有的
	global.failedHandlers[nodeName][name] = call
}
func (global *httpGlobal) DeniedHandler(name string, call WebFunc) {
	global.mutex.Lock()
	defer global.mutex.Unlock()

	if global.deniedHandlers == nil {
		global.deniedHandlers = map[string]map[string]WebFunc{}
	}
	if global.deniedHandlerNames == nil {
		global.deniedHandlerNames =  map[string][]string{}
	}


	//节点
	nodeName := ConstNodeGlobal
	if Current != "" {
		nodeName = Current
	}


	//如果节点配置不存在，创建
	if global.deniedHandlers[nodeName] == nil {
		global.deniedHandlers[nodeName] = map[string]WebFunc{}
	}
	if global.deniedHandlerNames[nodeName] == nil {
		global.deniedHandlerNames[nodeName] = []string{}
	}




	//如果没有注册个此name，才加入数组
	if _,ok := global.deniedHandlers[nodeName][name]; ok == false {
		global.deniedHandlerNames[nodeName] = append(global.deniedHandlerNames[nodeName], name)
	}
	//函数直接写， 因为可以使用同名替换现有的
	global.deniedHandlers[nodeName][name] = call
}
/* 注册处理器 end */



