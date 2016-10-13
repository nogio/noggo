/*
	http	HTTP模块
	HTTP功能是一个周期性的功能，定时周期执行
	主要使用场景：定时提醒，备份啊等等

	具体的HTTP执行代码，都在节点中，而非全局
*/

package noggo

import (
	. "github.com/nogio/noggo/base"
	"sync"
	"net/http"
	"strings"
	"regexp"
	"net/url"
	"fmt"
	"encoding/json"
	"encoding/xml"
)


type (

	HttpAcceptFunc func(res http.ResponseWriter, req *http.Request)

	//HTTP驱动
	HttpDriver interface {
		Connect(config Map) (HttpConnect)
	}
	//HTTP连接
	HttpConnect interface {
		//打开
		Open() error
		//关闭
		Close() error
		//注册
		Accept(call HttpAcceptFunc) error

		//开始
		Start(addr string) error
		//开始SSL
		StartTLS(addr string, certFile, keyFile string) error
	}
	//HTTP全局容器
	httpGlobal	struct {
		mutex sync.Mutex
		                                             //驱动
		drivers map[string]HttpDriver


		//路由
		routes 		map[string]map[string]Map	//路由定义							map[node]map[name]Map
		routeNames	map[string][]string			//路由名称原始顺序，因为map是无序的		map[node][]string

		//拦截器们
		requestFilters, executeFilters, responseFilters map[string]map[string]HttpFunc
		requestFilterNames, executeFilterNames, responseFilterNames map[string][]string

		                                             //处理器们
		foundHandlers, errorHandlers, failedHandlers, deniedHandlers map[string]map[string]HttpFunc
		foundHandlerNames, errorHandlerNames, failedHandlerNames, deniedHandlerNames map[string][]string
	}

)

//HTTP：连接驱动
func (module *httpGlobal) connect(config *httpConfig) (HttpConnect) {
	if httpDriver,ok := module.drivers[config.Driver]; ok {
		return httpDriver.Connect(config.Config)
	} else {
		panic("HTTP：不支持的驱动 " + config.Driver)
	}
}


//注册HTTP驱动
func (global *httpGlobal) Driver(name string, driver HttpDriver) {
	global.mutex.Lock()
	defer global.mutex.Unlock()

	if global.drivers == nil {
		global.drivers = map[string]HttpDriver{}
	}

	if driver == nil {
		panic("HTTP: 驱动不可为空")
	}
	//不做存在判断，因为要支持后注册的驱动替换已注册的驱动
	//框架有可能自带几种默认驱动，并且是默认注册的，用户可以自行注册替换
	global.drivers[name] = driver
}





//HTTP全局初始化
func (global *httpGlobal) init() {
	//HTTP全局无需处理任何东西
}

//HTTP全局退出
func (global *httpGlobal) exit() {
	//HTTP全局无需处理任何东西
}








//HTTP：注册路由
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
func (global *httpGlobal) RequestFilter(name string, call HttpFunc) {
	global.mutex.Lock()
	defer global.mutex.Unlock()

	if global.requestFilters == nil {
		global.requestFilters = map[string]map[string]HttpFunc{}
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
		global.requestFilters[nodeName] = map[string]HttpFunc{}
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
func (global *httpGlobal) ExecuteFilter(name string, call HttpFunc) {
	global.mutex.Lock()
	defer global.mutex.Unlock()

	if global.executeFilters == nil {
		global.executeFilters = map[string]map[string]HttpFunc{}
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
		global.executeFilters[nodeName] = map[string]HttpFunc{}
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
func (global *httpGlobal) ResponseFilter(name string, call HttpFunc) {
	global.mutex.Lock()
	defer global.mutex.Unlock()

	if global.responseFilters == nil {
		global.responseFilters = map[string]map[string]HttpFunc{}
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
		global.responseFilters[nodeName] = map[string]HttpFunc{}
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



//找不到处理器
func (global *httpGlobal) FoundHandler(name string, call HttpFunc) {
	global.mutex.Lock()
	defer global.mutex.Unlock()

	if global.foundHandlers == nil {
		global.foundHandlers = map[string]map[string]HttpFunc{}
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
		global.foundHandlers[nodeName] = map[string]HttpFunc{}
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
//错误处理器
func (global *httpGlobal) ErrorHandler(name string, call HttpFunc) {
	global.mutex.Lock()
	defer global.mutex.Unlock()

	if global.errorHandlers == nil {
		global.errorHandlers = map[string]map[string]HttpFunc{}
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
		global.errorHandlers[nodeName] = map[string]HttpFunc{}
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



//失败处理器
func (global *httpGlobal) FailedHandler(name string, call HttpFunc) {
	global.mutex.Lock()
	defer global.mutex.Unlock()

	if global.failedHandlers == nil {
		global.failedHandlers = map[string]map[string]HttpFunc{}
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
		global.failedHandlers[nodeName] = map[string]HttpFunc{}
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

//拒绝处理器
func (global *httpGlobal) DeniedHandler(name string, call HttpFunc) {
	global.mutex.Lock()
	defer global.mutex.Unlock()

	if global.deniedHandlers == nil {
		global.deniedHandlers = map[string]map[string]HttpFunc{}
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
		global.deniedHandlers[nodeName] = map[string]HttpFunc{}
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
//-----------------------------------------------------------------------------------------------------------------------//












type (

	//HTTP上下文函数类型
	HttpFunc func(*HttpContext)


	httpUri struct {
		name string
		uri string
	}
	//HTTP模块
	httpModule struct {
		mutex sync.Mutex

		//会话配置与连接
		sessionConfig	*sessionConfig
		sessionConnect	SessionConnect

		//HTTP配置与连接
		httpConfig	*httpConfig
		httpConnect	HttpConnect


		//所在节点
		Node	*Noggo


		//路由
		routes 		map[string]Map			//路由定义
		routeNames	[]string				//路由名称原始顺序，因为map是无序的
		routeUris	[]httpUri		//路由指定

		//拦截器们
		requestFilters, executeFilters, responseFilters map[string]HttpFunc
		requestFilterNames, executeFilterNames, responseFilterNames []string

		//处理器们
		foundHandlers, errorHandlers, failedHandlers, deniedHandlers map[string]HttpFunc
		foundHandlerNames, errorHandlerNames, failedHandlerNames, deniedHandlerNames []string
	}

	//HTTP上下文
	HttpContext struct {
		Node	*Noggo
		Module	*httpModule

		//执行线
		nexts []HttpFunc		//方法列表
		next int				//下一个索引


		Req	*http.Request
		Res	http.ResponseWriter

		//基础
		Id	string			//Session Id  会话时使用
		Session Map			//存储Session值
		Sign	*Sign		//签名功能，基于session


		//请求相关
		Method	string		//请求的method， 继承之http请求， 暂时无用
		Host	string		//请求的域名
		Path	string		//请求的路径，演变自http， 暂时等于http的名称
		Lang	string		//当前上下文的语言，默认应为default
		Ajax	bool		//标记当前请求是否ajax

		//配置相关
		Name string			//解析路由后得到的name
		Config Map			//解析后得到的路由配置
		Branchs []Map		//解析后得到的路由分支配置

		//表单相关
		Param	Map			//路由解析后uri中的参数
		Query	Map			//请求中的GET参数
		Form	Map			//请求的表单数据
		Upload	Map			//从表单中提取的文件信息
		Client	Map			//从请求中解析得到客户端信息保存于此，比如，包名，版本，或是浏览器信息等等，可自行在拦截器中处理


		//数据相关
		Value	Map			//所有请求过来的原始参数汇总
		Local	Map			//在ctx中传递数据用的
		Item	Map			//单条记录查询对象
		Auth	Map			//签名认证对象
		Args	Map			//经过args处理后的参数


		//响应相关
		Charset	string
		Code	int			//返回的状态码
		Type	string		//响应类型
		Body	Any			//响应内容
		Data	Map			//返回给view层的数据

		//错误信息，failed,error,denied上下文中使用
		Wrong	*Error		//错误信息
	}


	//响应跳转
	httpBodyGoto struct {
		Url string
	}
	//响应文本
	httpBodyText struct {
		Text	string
	}
	//响应HTML
	httpBodyHtml struct {
		Html	string
	}
	//响应脚本
	httpBodyScript struct {
		Script	string
	}
	httpBodyJson struct {
		//响应的Json对象
		Json Any
	}
	httpBodyXml struct {
		//响应的Xml对象
		Xml Any
	}
	//响应文件下载
	httpBodyFile struct {
		//要下载的文件路径
		File string
		//自定义下载文件名
		Name string
	}
	//响应二进制下载
	httpBodyDown struct {
		//要下载的文件内容
		Body []byte
		//自定义下载文件名
		Name string
	}
	//响应视图
	httpBodyView struct {
		View string
		Model Map
	}
)



//HTTP模块初始化
func (module *httpModule) run() {
	module.runSession()
	module.runHttp()
}
func (module *httpModule) runSession() {
	if Config.Http.Session != nil {
		//使用HTTP中的会话配置
		module.sessionConfig = Config.Http.Session
	} else {
		//使用默认的会话配置
		module.sessionConfig = Config.Session
	}

	//连接会话
	module.sessionConnect = Session.connect(module.sessionConfig)

	if module.sessionConnect == nil {
		panic("节点HTTP：连接会话失败")
	} else {
		//打开会话连接
		err := module.sessionConnect.Open()
		if err != nil {
			panic("节点HTTP：打开会话失败 " + err.Error())
		}
	}
}
func (module *httpModule) runHttp() {

	module.httpConfig = Config.Http
	module.httpConnect = Http.connect(module.httpConfig)

	if module.httpConnect == nil {
		panic("节点HTTP：连接失败")
	} else {
		//打开会话连接
		err := module.httpConnect.Open()
		if err != nil {
			panic("节点HTTP：打开失败 " + err.Error())
		}
	}

	//监听
	module.httpConnect.Accept(module.serveHttp)

	//开始HTTP
	//这里要判断是否SSL，如果是应该开始SSL
	//注意，connect.Start不应该阻塞线程
	module.httpConnect.Start(module.Node.Port)

}


//HTTP模块退出
func (module *httpModule) end() {
	module.endSession()
	module.endHttp()
}
//退出SESSION
func (module *httpModule) endSession() {
	if module.sessionConnect != nil {
		module.sessionConnect.Close()
		module.sessionConnect = nil
	}
}
//退出HTTP本身
func (module *httpModule) endHttp() {
	if module.httpConnect != nil {
		module.httpConnect.Close()
		module.httpConnect = nil
	}
}








//任务：注册路由
func (module *httpModule) Route(name string, config Map) {
	module.mutex.Lock()
	defer module.mutex.Unlock()


	if module.routes == nil {
		module.routes = map[string]Map{}
	}
	if module.routeNames == nil {
		module.routeNames = []string{}
	}

	//保存配置
	if _,ok := module.routes[name]; ok == false {
		//没有注册过name，才把name加到列表
		module.routeNames = append(module.routeNames, name)
	}
	//可以后注册重写原有路由配置，所以直接保存
	module.routes[name] = config





	if module.routeUris == nil {
		module.routeUris = []httpUri{}
	}


	//处理uri
	if v,ok := config[KeyMapUri]; ok {
		switch uris := v.(type) {
		case string:
			module.routeUris = append(module.routeUris, httpUri{ name, uris })
		case []string:
			for _,uri := range uris {
				module.routeUris = append(module.routeUris, httpUri{ name, uri })
			}
		}
	}
}








/* 注册拦截器 begin */
func (module *httpModule) RequestFilter(name string, call HttpFunc) {
	module.mutex.Lock()
	defer module.mutex.Unlock()

	if module.requestFilters == nil {
		module.requestFilters = make(map[string]HttpFunc)
	}
	if module.requestFilterNames == nil {
		module.requestFilterNames = make([]string, 0)
	}

	//如果没有注册个此name，才加入数组
	if _,ok := module.requestFilters[name]; ok == false {
		module.requestFilterNames = append(module.requestFilterNames, name)
	}
	//函数直接写， 因为可以使用同名替换现有的
	module.requestFilters[name] = call
}
func (module *httpModule) ExecuteFilter(name string, call HttpFunc) {
	module.mutex.Lock()
	defer module.mutex.Unlock()

	if module.executeFilters == nil {
		module.executeFilters = make(map[string]HttpFunc)
	}
	if module.executeFilterNames == nil {
		module.executeFilterNames = make([]string, 0)
	}

	//如果没有注册个此name，才加入数组
	if _,ok := module.executeFilters[name]; ok == false {
		module.executeFilterNames = append(module.executeFilterNames, name)
	}
	//函数直接写， 因为可以使用同名替换现有的
	module.executeFilters[name] = call
}
func (module *httpModule) ResponseFilter(name string, call HttpFunc) {
	module.mutex.Lock()
	defer module.mutex.Unlock()

	if module.responseFilters == nil {
		module.responseFilters = make(map[string]HttpFunc)
	}
	if module.responseFilterNames == nil {
		module.responseFilterNames = make([]string, 0)
	}

	//如果没有注册个此name，才加入数组
	if _,ok := module.responseFilters[name]; ok == false {
		module.responseFilterNames = append(module.responseFilterNames, name)
	}
	//函数直接写， 因为可以使用同名替换现有的
	module.responseFilters[name] = call
}
/* 注册拦截器 end */


/* 注册处理器 begin */
func (module *httpModule) FoundHandler(name string, call HttpFunc) {
	module.mutex.Lock()
	defer module.mutex.Unlock()

	if module.foundHandlers == nil {
		module.foundHandlers = make(map[string]HttpFunc)
	}
	if module.foundHandlerNames == nil {
		module.foundHandlerNames = make([]string, 0)
	}

	//如果没有注册个此name，才加入数组
	if _,ok := module.foundHandlers[name]; ok == false {
		module.foundHandlerNames = append(module.foundHandlerNames, name)
	}
	//函数直接写， 因为可以使用同名替换现有的
	module.foundHandlers[name] = call
}
func (module *httpModule) ErrorHandler(name string, call HttpFunc) {
	module.mutex.Lock()
	defer module.mutex.Unlock()

	if module.errorHandlers == nil {
		module.errorHandlers = make(map[string]HttpFunc)
	}
	if module.errorHandlerNames == nil {
		module.errorHandlerNames = make([]string, 0)
	}

	//如果没有注册个此name，才加入数组
	if _,ok := module.errorHandlers[name]; ok == false {
		module.errorHandlerNames = append(module.errorHandlerNames, name)
	}
	//函数直接写， 因为可以使用同名替换现有的
	module.errorHandlers[name] = call
}
func (module *httpModule) FailedHandler(name string, call HttpFunc) {
	module.mutex.Lock()
	defer module.mutex.Unlock()

	if module.failedHandlers == nil {
		module.failedHandlers = make(map[string]HttpFunc)
	}
	if module.failedHandlerNames == nil {
		module.failedHandlerNames = make([]string, 0)
	}

	//如果没有注册个此name，才加入数组
	if _,ok := module.failedHandlers[name]; ok == false {
		module.failedHandlerNames = append(module.failedHandlerNames, name)
	}
	//函数直接写， 因为可以使用同名替换现有的
	module.failedHandlers[name] = call
}
func (module *httpModule) DeniedHandler(name string, call HttpFunc) {
	module.mutex.Lock()
	defer module.mutex.Unlock()

	if module.deniedHandlers == nil {
		module.deniedHandlers = make(map[string]HttpFunc)
	}
	if module.deniedHandlerNames == nil {
		module.deniedHandlerNames = make([]string, 0)
	}

	//如果没有注册个此name，才加入数组
	if _,ok := module.deniedHandlers[name]; ok == false {
		module.deniedHandlerNames = append(module.deniedHandlerNames, name)
	}
	//函数直接写， 因为可以使用同名替换现有的
	module.deniedHandlers[name] = call
}
/* 注册处理器 end */








//创建HTTP模块
func newHttpModule(node *Noggo) (*httpModule) {
	module := &httpModule{
		Node: node,
	}


	//复制路由，拦截器，处理器


	//全局路由
	routes, routesOK := Http.routes[ConstNodeGlobal];
	routeNames, routeNamesOK := Http.routeNames[ConstNodeGlobal];
	if routesOK && routeNamesOK {
		for _,n := range routeNames {
			module.Route(n, routes[n])
		}
	}
	//节点路由
	nodeRoutes, nodeRoutesOK := Http.routes[node.Name];
	nodeRouteNames, nodeRouteNamesOK := Http.routeNames[node.Name];
	if nodeRoutesOK && nodeRouteNamesOK {
		for _,n := range nodeRouteNames {
			module.Route(n, nodeRoutes[n])
		}
	}



	//全局 请求拦截器
	requestFilters, requestFiltersOK := Http.requestFilters[ConstNodeGlobal];
	requestFilterNames, requestFilterNamesOK := Http.requestFilterNames[ConstNodeGlobal];
	if requestFiltersOK && requestFilterNamesOK {
		for _,n := range requestFilterNames {
			module.RequestFilter(n, requestFilters[n])
		}
	}
	//节点 请求拦截器
	nodeRequestFilters, nodeRequestFiltersOK := Http.requestFilters[node.Name];
	nodeRequestFilterNames, nodeRequestFilterNamesOK := Http.requestFilterNames[node.Name];
	if nodeRequestFiltersOK && nodeRequestFilterNamesOK {
		for _,n := range nodeRequestFilterNames {
			module.RequestFilter(n, nodeRequestFilters[n])
		}
	}
	//全局 执行拦截器
	executeFilters, executeFiltersOK := Http.executeFilters[ConstNodeGlobal];
	executeFilterNames, executeFilterNamesOK := Http.executeFilterNames[ConstNodeGlobal];
	if executeFiltersOK && executeFilterNamesOK {
		for _,n := range executeFilterNames {
			module.ExecuteFilter(n, executeFilters[n])
		}
	}
	//节点 执行拦截器
	nodeExecuteFilters, nodeExecuteFiltersOK := Http.executeFilters[node.Name];
	nodeExecuteFilterNames, nodeExecuteFilterNamesOK := Http.executeFilterNames[node.Name];
	if nodeExecuteFiltersOK && nodeExecuteFilterNamesOK {
		for _,n := range nodeExecuteFilterNames {
			module.ExecuteFilter(n, nodeExecuteFilters[n])
		}
	}
	//全局 响应拦截器
	responseFilters, responseFiltersOK := Http.responseFilters[ConstNodeGlobal];
	responseFilterNames, responseFilterNamesOK := Http.responseFilterNames[ConstNodeGlobal];
	if responseFiltersOK && responseFilterNamesOK {
		for _,n := range responseFilterNames {
			module.ResponseFilter(n, responseFilters[n])
		}
	}
	//节点 响应拦截器
	nodeResponseFilters, nodeResponseFiltersOK := Http.responseFilters[node.Name];
	nodeResponseFilterNames, nodeResponseFilterNamesOK := Http.responseFilterNames[node.Name];
	if nodeResponseFiltersOK && nodeResponseFilterNamesOK {
		for _,n := range nodeResponseFilterNames {
			module.ResponseFilter(n, nodeResponseFilters[n])
		}
	}







	//注意，处理器，优先注册节点的，在前面
	//全局的在后面， 这样可以做一个全局的默认处理器
	//然后在有些节点，可以优先自定义处理


	//节点 不存在处理器
	nodeFoundHandlers, nodeFoundHandlersOK := Http.foundHandlers[node.Name];
	nodeFoundHandlerNames, nodeFoundHandlerNamesOK := Http.foundHandlerNames[node.Name];
	if nodeFoundHandlersOK && nodeFoundHandlerNamesOK {
		for _,n := range nodeFoundHandlerNames {
			module.FoundHandler(n, nodeFoundHandlers[n])
		}
	}
	//全局 不存在处理器
	foundHandlers, foundHandlersOK := Http.foundHandlers[ConstNodeGlobal];
	foundHandlerNames, foundHandlerNamesOK := Http.foundHandlerNames[ConstNodeGlobal];
	if foundHandlersOK && foundHandlerNamesOK {
		for _,n := range foundHandlerNames {
			//不存在才注册，这样全局不会替换节点的处理器
			if _,ok := module.foundHandlers[n]; ok == false {
				module.FoundHandler(n, foundHandlers[n])
			}
		}
	}




	//节点 错误处理器
	nodeErrorHandlers, nodeErrorHandlersOK := Http.errorHandlers[node.Name];
	nodeErrorHandlerNames, nodeErrorHandlerNamesOK := Http.errorHandlerNames[node.Name];
	if nodeErrorHandlersOK && nodeErrorHandlerNamesOK {
		for _,n := range nodeErrorHandlerNames {
			module.ErrorHandler(n, nodeErrorHandlers[n])
		}
	}
	//全局 错误处理器
	errorHandlers, errorHandlersOK := Http.errorHandlers[ConstNodeGlobal];
	errorHandlerNames, errorHandlerNamesOK := Http.errorHandlerNames[ConstNodeGlobal];
	if errorHandlersOK && errorHandlerNamesOK {
		for _,n := range errorHandlerNames {
			//不存在才注册，这样全局不会替换节点的处理器
			if _,ok := module.errorHandlers[n]; ok == false {
				module.ErrorHandler(n, errorHandlers[n])
			}
		}
	}


	//节点 失败处理器
	nodeFailedHandlers, nodeFailedHandlersOK := Http.failedHandlers[node.Name];
	nodeFailedHandlerNames, nodeFailedHandlerNamesOK := Http.failedHandlerNames[node.Name];
	if nodeFailedHandlersOK && nodeFailedHandlerNamesOK {
		for _,n := range nodeFailedHandlerNames {
			module.FailedHandler(n, nodeFailedHandlers[n])
		}
	}
	//全局 失败处理器
	failedHandlers, failedHandlersOK := Http.failedHandlers[ConstNodeGlobal];
	failedHandlerNames, failedHandlerNamesOK := Http.failedHandlerNames[ConstNodeGlobal];
	if failedHandlersOK && failedHandlerNamesOK {
		for _,n := range failedHandlerNames {
			//不存在才注册，这样全局不会替换节点的处理器
			if _,ok := module.failedHandlers[n]; ok == false {
				module.FailedHandler(n, failedHandlers[n])
			}
		}
	}



	//节点 请求拦截器
	nodeDeniedHandlers, nodeDeniedHandlersOK := Http.deniedHandlers[node.Name];
	nodeDeniedHandlerNames, nodeDeniedHandlerNamesOK := Http.deniedHandlerNames[node.Name];
	if nodeDeniedHandlersOK && nodeDeniedHandlerNamesOK {
		for _,n := range nodeDeniedHandlerNames {
			module.DeniedHandler(n, nodeDeniedHandlers[n])
		}
	}
	//全局 请求拦截器
	deniedHandlers, deniedHandlersOK := Http.deniedHandlers[ConstNodeGlobal];
	deniedHandlerNames, deniedHandlerNamesOK := Http.deniedHandlerNames[ConstNodeGlobal];
	if deniedHandlersOK && deniedHandlerNamesOK {
		for _,n := range deniedHandlerNames {
			//不存在才注册，这样全局不会替换节点的处理器
			if _,ok := module.deniedHandlers[n]; ok == false {
				module.DeniedHandler(n, deniedHandlers[n])
			}
		}
	}


	return module
}


















//创建Http上下文
func (module *httpModule) newHttpContext(res http.ResponseWriter, req *http.Request) (*HttpContext) {
	method := strings.ToLower(req.Method)
	host := req.URL.Host
	path := req.URL.Path
	value := Map{}

	return &HttpContext{
		Module: module,
		next: -1, nexts: []HttpFunc{},

		Res: res, Req: req,

		Method: method, Host: host, Path: path,
		Ajax: false, Lang: "default",

		Name: "", Config: Map{}, Branchs:[]Map{},

		Param: Map{}, Query: Map{}, Form: Map{}, Upload: Map{}, Client: Map{},
		Value: value, Local: Map{}, Item: Map{}, Auth: Map{}, Args: Map{},
	}
}



//HTTPHttp  请求开始
func (module *httpModule) serveHttp(res http.ResponseWriter, req *http.Request) {
	ctx := module.newHttpContext(res, req)

	//请求处理
	ctx.handler(module.contextRoute)
	ctx.handler(module.contextRequest)
	//最终所有的响应处理，优先
	ctx.handler(module.contextResponse)
	//filter中的request
	//用数组保证原始注册顺序
	for _,name := range module.requestFilterNames {
		ctx.handler(module.requestFilters[name])
	}

	//开始执行
	ctx.handler(module.contextExecute)
	ctx.Next()
}







/*
	HTTP模块处理方法 begin
*/


//http HTTP处理器
//直接路由了，暂不支持三方路由器驱动了
func (module *httpModule) contextRoute(ctx *HttpContext) {

	/*
	//路由解析
	routerResult := module.routerConnect.Parse(ctx.Host, ctx.Path)
	if routerResult != nil {
		ctx.Name = routerResult.Name
		ctx.Config = module.routes[ctx.Name]
		ctx.Param = routerResult.Param
	} else {
		ctx.Config = nil
	}
	*/


	//解析成功
	ctx.Name = ""
	ctx.Config = nil


	//路由解析
	for _,obj := range module.routeUris {
		uri := obj.uri

		keys, vals := []string{}, []string{}

		//用来匹配的正则表达式
		regs := strings.Replace(uri, ".", "\\.", -1)
		regs = strings.Replace(regs, "/", "\\/", -1)


		//正则对象
		regx := regexp.MustCompile(`\{[\*A-Za-z0-9_]+\}`)
		//拿到URI中的参数列表
		//keys := regx.FindAllString(uri, -1)
		//替换参数为正则
		regs = regx.ReplaceAllStringFunc(regs, func(p string) string {
			if (p[1:1] == "*") {
				keys = append(keys, p[2:len(p)-1])
				return `(.*)`
			} else {
				keys = append(keys, p[1:len(p)-1])
				return `([-_A-Za-z0-9\.]+)`
			}
		})



		//正则匹配当前URL，是否从域名开始匹配
		//忽略query部分
		url := ctx.Path
		if (uri[0:1] != "/") {
			url = ctx.Host + ctx.Path
		}


		regx = regexp.MustCompile("^"+regs+"$")
		if regx.MatchString(url) {

			matchs := regx.FindAllStringSubmatch(url, -1)
			if matchs != nil && len(matchs) > 0 && len(matchs[0]) > 0 {
				vals = matchs[0][1:]
			}

			param := Map{}

			//处理params
			if len(keys) == len(vals) {
				for i,k := range keys {
					param[k] = vals[i]
				}
			}

			//解析成功
			ctx.Name = obj.name
			ctx.Config = module.routes[ctx.Name]
			ctx.Param = param
			break
		}
	}

	ctx.Next()
}


//http HTTP处理器
//请求处理
//包含：route解析、request处理、session处理
func (module *httpModule) contextRequest(ctx *HttpContext) {

	//会话处理相关
	cookie, err := ctx.Req.Cookie(Config.Http.Cookie)
	if err != nil || cookie.Value == "" {
		ctx.Id = NewMd5Id()
		err,m := module.sessionConnect.Query(ctx.Id, module.sessionConfig.Expiry)
		if err == nil {
			ctx.Session = m
		} else {
			ctx.Session = Map{}
		}

		//注意域设置
		//这里在ctx.Cookie封装之后，应当写到ctx.Cookie中，在response写入res
		//待修改
		cookie := http.Cookie{ Name: Config.Http.Cookie, Value: url.QueryEscape(ctx.Id), Path: "/", HttpOnly: true }
		if Config.Session.Expiry > 0 {
			cookie.MaxAge = int(Config.Session.Expiry)
		}
		if Config.Http.Domain != "" {
			cookie.Domain = Config.Http.Domain
		}
		http.SetCookie(ctx.Res, &cookie)
	} else {
		ctx.Id, _ = url.QueryUnescape(cookie.Value)

		err,m := module.sessionConnect.Query(ctx.Id, module.sessionConfig.Expiry)
		if err == nil {
			ctx.Session = m
		} else {
			ctx.Session = Map{}
		}
	}



	//基本的表单处理
	//如 query form

	//处理Query，不管任何method都要处理
	if querys, err := url.ParseQuery(ctx.Req.URL.RawQuery); err == nil {
		for k,v := range querys {
			if len(v) > 1 {
				ctx.Query[k] = v
				ctx.Value[k] = v
			} else {
				ctx.Query[k] = v[0]
				ctx.Value[k] = v[0]
			}
		}
	}



	ctx.Sign = &Sign{ ctx.Session }
	ctx.Next()
	module.sessionConnect.Update(ctx.Id, ctx.Session, module.sessionConfig.Expiry)
}





//处理响应
func (module *httpModule) contextResponse(ctx *HttpContext) {
	//因为response是在所有请求前的， 所以先调用一下
	//然后对结果进行处理
	ctx.Next()


	//清理执行线
	ctx.cleanup()

	//filter中的request
	//用数组保证原始注册顺序
	for _,name := range module.responseFilterNames {
		ctx.handler(module.responseFilters[name])
	}

	//这个函数才是真正响应的处理函数
	ctx.handler(module.contextResponder)

	ctx.Next()

}







//路由执行，处理
func (module *httpModule) contextExecute(ctx *HttpContext) {

	//解析路由，拿到actions
	if ctx.Config == nil {
		//找不到路由
		ctx.handler(module.contextFound)
	} else {


		//验证，参数，数据处理
		//验证处理，数据处理， 可以考虑走外部中间件
		if _,ok := ctx.Config[KeyMapArgs]; ok {
			ctx.handler(module.contextArgs)
		}
		if _,ok := ctx.Config[KeyMapAuth]; ok {
			ctx.handler(module.contextAuth)
		}
		if _,ok := ctx.Config[KeyMapItem]; ok {
			ctx.handler(module.contextItem)
		}

		//最终都由分支处理
		ctx.handler(module.contextBranch)
	}

	ctx.Next()
}


//HTTP处理：处理分支
func (module *httpModule) contextBranch(ctx *HttpContext) {

	//执行线重置
	ctx.cleanup()
	ctx.Branchs = []Map{}

	//总体思路，考虑路由和分支
	//把路由本身，做为一个匹配所有的分支，放到最后一个执行

	//如果有分支，来
	if branchConfig,ok := ctx.Config[KeyMapBranch]; ok {
		//遍历分支
		for _,v := range branchConfig.(Map) {
			//保存了：branch.xxx { match, route }
			ctx.Branchs = append(ctx.Branchs, v.(Map))
		}
	}
	//如果有路由
	if routeConfig,ok := ctx.Config[KeyMapRoute]; ok {
		//保存{ match, route }
		ctx.Branchs = append(ctx.Branchs, Map{
			KeyMapMatch:	true,	//默认路由直接匹配
			KeyMapRoute:	routeConfig,
		})
	}

	var routing Map

	forBranchs:
	for _,b := range ctx.Branchs {
		if match,ok := b[KeyMapMatch]; ok {

			switch match:=match.(type) {
			case bool:
				if (match) {
					routing = b
					break forBranchs;
				}
			case func(*HttpContext)bool:
				if (match(ctx)) {
					routing = b
					break forBranchs;
				}
			default:
			}
		}
	}


	/*
	//先不复制了吧，因为顶级的，在已经处理过 params,args,auth等的东西，再复制会重复处理
	//而且复制的话， 还得判断auth, item的子级map， 合并到一起
	for k,v := range ctx.Route {
		if k != "uri" && k != "match" && k != "route" && k != "branch" {
			routing[k] = v
		}
	}
	*/


	//这里 ctx.Route 和 routing 变换位置
	ctx.Config = Map{}

	//如果有路由
	//HTTP路由不支持多method，非http
	if routeConfig,ok := routing[KeyMapRoute]; ok {

		for k,v := range routeConfig.(Map) {
			ctx.Config[k] = v
		}

		/*
		//如果是method=*版
		if _,ok := routeConfig.(Map)[KeyMapAction]; ok {
			for k,v := range routeConfig.(Map) {
				ctx.Config[k] = v
			}
		} else {	//否则为方法版：get,post
			if methodConfig,ok := routeConfig.(Map)[ctx.Method]; ok {
				for k,v := range methodConfig.(Map) {
					ctx.Config[k] = v
				}
			}
		}
		*/
	}




	//先处理参数，验证等的东西
	if _,ok := ctx.Config[KeyMapArgs]; ok {
		ctx.handler(module.contextArgs)
	}
	if _,ok := ctx.Config[KeyMapAuth]; ok {
		ctx.handler(module.contextAuth)
	}
	if _,ok := ctx.Config[KeyMapItem]; ok {
		ctx.handler(module.contextItem)
	}


	//action之前的拦截器
	//filter中的execute
	//用数组保证原始注册顺序
	for _,name := range module.executeFilterNames {
		ctx.handler(module.executeFilters[name])
	}

	//把action加入调用列表
	if actionConfig,ok := ctx.Config[KeyMapAction]; ok {
		switch actions:=actionConfig.(type) {
		case func(*HttpContext):
			ctx.handler(actions)
		case []func(*HttpContext):
			for _,action := range actions {
				ctx.handler(action)
			}
		case HttpFunc:
			ctx.handler(actions)
		case []HttpFunc:
			ctx.handler(actions...)
		default:
		}
	}

	ctx.Next()
}

















//自带中间件，参数处理
func (module *httpModule) contextArgs(ctx *HttpContext) {

	//argn表示参数都可为空
	argn := false
	if v,ok := ctx.Config["argn"].(bool); ok {
		argn = v
	}

	//所有值都会放在 module.Value 中
	err := Mapping.Parse([]string{}, ctx.Config["args"].(Map), ctx.Value, ctx.Args, argn)
	if err != nil {
		ctx.Failed(err)
	} else {
		ctx.Next()
	}
}



//Auth验证处理
func (module *httpModule) contextAuth(ctx *HttpContext) {

	if auths,ok := ctx.Config["auth"]; ok {
		saveMap := Map{}

		for authKey,authMap := range auths.(Map) {

			ohNo := false
			authConfig := authMap.(Map)

			if authConfig["sign"] == nil {
				continue
			}

			authSign := authConfig["sign"].(string)
			authMust := false
			authName := authSign

			if authConfig["must"] != nil {
				authMust = authConfig["must"].(bool)
			}
			if authConfig["name"] != nil {
				authName = authConfig["name"].(string)
			}

			//判断是否登录
			if ctx.Sign.Yes(authSign) {

				/*
				因为数据层还没上， 所以暂不支持，以下查询数据库的操作
				//判断是否需要查询数据
				dataName,dok := authConfig["data"]; modelName,mok := authConfig["model"];
				if dok && mok {

					//要查询库
					//不管must是否,都要查库
					db := Data.Data(dataName.(string)); defer db.Close()
					item := db.Model(modelName.(string)).Entity(ctx.Sign.Id(authSign))
					if item != nil {
						saveMap[authKey] = item
					} else {
						if authMust {	//是必要的
							//是否有自定义状态
							err := NewStateError("auth.error", authName)
							if v,ok := authConfig["error"]; ok {
								err = NewStateError(v.(string))
							}

							err.Data = authConfig
							ctx.Denied(authKey, err)
							return;
						}
					}


				} else {
					//无需data, model， 不管
				}
				*/

			} else {
				ohNo = true
			}

			//到这里是未登录的
			//而且是必须要登录，才显示错误
			if ohNo && authMust {

				//是否有自定义状态
				err := Const.NewStateError("auth.empty", authName)
				if v,ok := authConfig["empty"]; ok {
					err = Const.NewStateError(v.(string))
				}

				//貌似不需要这个
				//err.Data = authConfig

				//指定错误类型为authKey
				err.Type = authKey
				ctx.Denied(err)
				return;

			}
		}

		//存入
		for k,v := range saveMap {
			ctx.Auth[k] = v
		}
	}

	ctx.Next()
}
//Entity实体处理
func (module *httpModule) contextItem(ctx *HttpContext) {
	if ctx.Config["item"] != nil {
		cfg := ctx.Config["item"].(Map)

		saveMap := Map{}

		for k,v := range cfg {
			config := v.(Map)

			name := config["name"].(string)
			key := k
			if config["key"] != nil && config["key"] != "" {
				key = config["key"].(string)
			}

			if ctx.Value[key] == nil {
				//参数不为空啊啊
				state := "item.empty"
				//是否有自定义状态
				if v,ok := config["empty"]; ok {
					state = v.(string)
				}
				err := Const.NewStateError(state, name)

				//指定错误类型为item的key，好在处理时区分
				err.Type = k
				//查询不到东西，也要失败， 接口访问失败
				ctx.Failed(err)
				return
			} else {

				/*
				由于数据层还未完工，暂不支持数据查询
				//判断是否需要查询数据
				dataName,dok := config["data"]; modelName,mok := config["model"];
				if dok && mok {

					//要查询库
					db := Data.Data(dataName.(string)); defer db.Close()
					item := db.Model(modelName.(string)).Entity(ctx.Value[key])
					if item != nil {
						saveMap[k] = item
					} else {
						state := "item.error"
						//是否有自定义状态
						if v,ok := config["error"]; ok {
							state = v.(string)
						}
						err := Const.NewStateError(state, name)

						//这个不需要了吧
						//err.Data = config

						//错误类型等于item.key,方便处理
						err.Type = k
						ctx.Failed(err)
						return;
					}
				}

				*/
			}
		}

		//存入
		for k,v := range saveMap {
			ctx.Item[k] = v
		}
	}
	ctx.Next()
}







//处理响应
//这个才是真的响应处理
func (module *httpModule) contextResponder(ctx *HttpContext) {

	//对响应的默认处理
	if ctx.Code == 0 {
		ctx.Code = 200
	}
	if ctx.Charset == "" {
		ctx.Charset = Config.Http.Charset
		if ctx.Charset == "" {
			ctx.Charset = "utf-8"
		}
	}

	/*
	if ctx.Body == nil {
		//没有响应，应该走到found流程
		module.contextFound(ctx)
	}
	*/


	switch body := ctx.Body.(type) {
	case httpBodyGoto:
		module.gotoResponder(ctx)
	case httpBodyText:
		module.textResponder(ctx)
	case httpBodyHtml:
		module.htmlResponder(ctx)
	case httpBodyScript:
		module.scriptResponder(ctx)
	case httpBodyJson:
		module.jsonResponder(ctx)
	case httpBodyXml:
		module.xmlResponder(ctx)
	case httpBodyFile:
		module.fileResponder(ctx)
	case httpBodyDown:
		module.downResponder(ctx)
	case httpBodyView:
		module.viewResponder(ctx)

	//以下是一些常用的数据类型要处理的
	case string:{
		ctx.Body = httpBodyText{ body }
		module.textResponder(ctx)
	}
	case Map,[]Map:{
		ctx.Body = httpBodyJson{ body }
		module.textResponder(ctx)
	}
	case []byte:{
		ctx.Body = httpBodyDown{ body,"" }
		module.textResponder(ctx)
	}

	default:
		module.defaultResponder(ctx)
	}
}



























//路由执行，found
func (module *httpModule) contextFound(ctx *HttpContext) {
	//清理执行线
	ctx.cleanup()

	//如果路由配置中有found，就自定义处理
	if v,ok := ctx.Config[KeyMapFound]; ok {
		switch c := v.(type) {
		case HttpFunc: {
			ctx.handler(c)
		}
		case []HttpFunc: {
			for _,v := range c {
				ctx.handler(v)
			}
		}
		case func(*HttpContext): {
			ctx.handler(c)
		}
		case []func(*HttpContext): {
			for _,v := range c {
				ctx.handler(v)
			}
		}
		default:
		}
	}


	//handler中的found
	//用数组保证原始注册顺序
	for _,name := range module.foundHandlerNames {
		ctx.handler(module.foundHandlers[name])
	}

	//最后是默认found中间件
	ctx.handler(module.foundDefaultHandler)

	ctx.Next()
}


//路由执行，error
func (module *httpModule) contextError(ctx *HttpContext) {
	//清理执行线
	ctx.cleanup()

	//如果路由配置中有found，就自定义处理
	if v,ok := ctx.Config[KeyMapError]; ok {
		switch c := v.(type) {
		case HttpFunc: {
			ctx.handler(c)
		}
		case []HttpFunc: {
			for _,v := range c {
				ctx.handler(v)
			}
		}
		case func(*HttpContext): {
			ctx.handler(c)
		}
		case []func(*HttpContext): {
			for _,v := range c {
				ctx.handler(v)
			}
		}
		default:
		}
	}


	//handler中的error
	//用数组保证原始注册顺序
	for _,name := range module.errorHandlerNames {
		ctx.handler(module.errorHandlers[name])
	}

	//最后是默认error中间件
	ctx.handler(module.errorDefaultHandler)

	ctx.Next()
}


//路由执行，failed
func (module *httpModule) contextFailed(ctx *HttpContext) {
	//清理执行线
	ctx.cleanup()

	//如果路由配置中有found，就自定义处理
	if v,ok := ctx.Config[KeyMapFailed]; ok {
		switch c := v.(type) {
		case HttpFunc: {
			ctx.handler(c)
		}
		case []HttpFunc: {
			for _,v := range c {
				ctx.handler(v)
			}
		}
		case func(*HttpContext): {
			ctx.handler(c)
		}
		case []func(*HttpContext): {
			for _,v := range c {
				ctx.handler(v)
			}
		}
		default:
		}
	}


	//handler中的failed
	//用数组保证原始注册顺序
	for _,name := range module.failedHandlerNames {
		ctx.handler(module.failedHandlers[name])
	}

	//最后是默认failed中间件
	ctx.handler(module.failedDefaultHandler)

	ctx.Next()
}



//路由执行，denied
func (module *httpModule) contextDenied(ctx *HttpContext) {
	//清理执行线
	ctx.cleanup()

	//如果路由配置中有found，就自定义处理
	if v,ok := ctx.Config[KeyMapDenied]; ok {
		switch c := v.(type) {
		case HttpFunc: {
			ctx.handler(c)
		}
		case []HttpFunc: {
			for _,v := range c {
				ctx.handler(v)
			}
		}
		case func(*HttpContext): {
			ctx.handler(c)
		}
		case []func(*HttpContext): {
			for _,v := range c {
				ctx.handler(v)
			}
		}
		default:
		}
	}

	//handler中的denied
	//用数组保证原始注册顺序
	for _,name := range module.deniedHandlerNames {
		ctx.handler(module.deniedHandlers[name])
	}

	//最后是默认denied中间件
	ctx.handler(module.deniedDefaultHandler)

	ctx.Next()
}



/*
	HTTP模块方法 end
*/






/* 默认处理器 begin */
//代码中没有指定相关的处理器，才会执行到默认处理器
func (module *httpModule) foundDefaultHandler(ctx *HttpContext) {
	ctx.Text("http not found")
}
func (module *httpModule) errorDefaultHandler(ctx *HttpContext) {
	ctx.Text(fmt.Sprintf("http error %v", ctx.Wrong))
}
func (module *httpModule) failedDefaultHandler(ctx *HttpContext) {
	ctx.Text(fmt.Sprintf("http failed %v", ctx.Wrong))
}
func (module *httpModule) deniedDefaultHandler(ctx *HttpContext) {
	ctx.Text(fmt.Sprintf("http denied %v", ctx.Wrong))
}
/* 默认处理器 end */








/* 响应器 begin */
func (module *httpModule) gotoResponder(ctx *HttpContext) {
	body := ctx.Body.(httpBodyGoto)
	http.Redirect(ctx.Res, ctx.Req, body.Url, http.StatusFound)
}
func (module *httpModule) textResponder(ctx *HttpContext) {
	body := ctx.Body.(httpBodyText)

	if ctx.Type == "" {
		ctx.Type = "text"
	}

	ctx.Res.Header().Set("Content-Type", fmt.Sprintf("%v; charset=%v", Const.MimeType(ctx.Type), ctx.Charset))
	ctx.Res.WriteHeader(ctx.Code)
	fmt.Fprint(ctx.Res, body.Text)
}
func (module *httpModule) htmlResponder(ctx *HttpContext) {
	body := ctx.Body.(httpBodyHtml)

	if ctx.Type == "" {
		ctx.Type = "html"
	}

	ctx.Res.Header().Add("Content-Type", fmt.Sprintf("%v; charset=%v", Const.MimeType(ctx.Type), ctx.Charset))
	ctx.Res.WriteHeader(ctx.Code)
	fmt.Fprint(ctx.Res, body.Html)
}
func (module *httpModule) scriptResponder(ctx *HttpContext) {
	body := ctx.Body.(httpBodyScript)

	if ctx.Type == "" {
		ctx.Type = "script"
	}


	ctx.Res.Header().Add("Content-Type", fmt.Sprintf("%v; charset=%v", Const.MimeType(ctx.Type), ctx.Charset))
	ctx.Res.WriteHeader(ctx.Code)
	fmt.Fprint(ctx.Res, body.Script)
}
func (module *httpModule) jsonResponder(ctx *HttpContext) {
	body := ctx.Body.(httpBodyJson)

	bytes, err := json.Marshal(body.Json)
	if err != nil {
		//出错啊
		//但是这里已经走完response了。再ctx.Error好像没用了
		//这是一个死循环， 因为走了ctx.Error， 还有可能再返回json
		//又会走到这里，在response中。 继续把response加入调用列表吧。这样保险
		ctx.Error(NewError(err.Error()))
	} else {

		if ctx.Type == "" {
			ctx.Type = "json"
		}

		ctx.Res.Header().Set("Content-Type", fmt.Sprintf("%v; charset=%v", Const.MimeType(ctx.Type), ctx.Charset))
		ctx.Res.WriteHeader(ctx.Code)
		fmt.Fprint(ctx.Res, string(bytes))
	}
}
func (module *httpModule) xmlResponder(ctx *HttpContext) {
	body := ctx.Body.(httpBodyXml)

	bytes, err := xml.Marshal(body.Xml)
	if err != nil {
		//出错啊
		ctx.Error(NewError(err.Error()))
	} else {

		if ctx.Type == "" {
			ctx.Type = "xml"
		}

		ctx.Res.Header().Set("Content-Type", fmt.Sprintf("%v; charset=%v", Const.MimeType(ctx.Type), ctx.Charset))
		ctx.Res.WriteHeader(ctx.Code)
		fmt.Fprint(ctx.Res, string(bytes))
	}
}
func (module *httpModule) fileResponder(ctx *HttpContext) {
	body := ctx.Body.(httpBodyFile)

	//加入自定义文件名
	if body.Name != "" {
		ctx.Res.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%v;", body.Name))
	}
	http.ServeFile(ctx.Res, ctx.Req, body.File)
}
func (module *httpModule) downResponder(ctx *HttpContext) {
	if ctx.Type == "" {
		ctx.Type = "down"
	}

	body := ctx.Body.(httpBodyDown)

	ctx.Res.Header().Set("Content-Type", fmt.Sprintf("%v; charset=%v", Const.MimeType(ctx.Type), ctx.Charset))
	//加入自定义文件名
	if body.Name != "" {
		ctx.Res.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%v;", body.Name))
	}
	ctx.Res.WriteHeader(ctx.Code)
	fmt.Fprint(ctx.Res, body.Body)
}
func (module *httpModule) viewResponder(ctx *HttpContext) {
	//先不支持VIEW的样子？
}
func (module *httpModule) defaultResponder(ctx *HttpContext) {
	if ctx.Type == "" {
		ctx.Type = "down"
	}

	ctx.Res.Header().Set("Content-Type", fmt.Sprintf("%v; charset=%v", Const.MimeType(ctx.Type), ctx.Charset))
	ctx.Res.WriteHeader(ctx.Code)
	fmt.Fprint(ctx.Res, ctx.Body)
}
/* 响应器 end */






































/*
	HTTP上下文处理方法 begin
*/



//添加执行线
func (ctx *HttpContext) handler(handlers ...HttpFunc) {
	for _,handler := range handlers {
		ctx.nexts = append(ctx.nexts, handler)
	}
}
//清空执行线
func (ctx *HttpContext) cleanup() {
	ctx.next = -1
	ctx.nexts = make([]HttpFunc, 0)
}

/* 执行下一个 */
func (ctx *HttpContext) Next() {
	ctx.next++
	if len(ctx.nexts) > ctx.next {
		next := ctx.nexts[ctx.next]
		if next != nil {
			next(ctx)
		} else {
			ctx.Next()
		}
	} else {
		//没有了，不再执行，Response会处理为404
	}
}





/* 上下文处理器 begin */
//不存在
func (ctx *HttpContext) Found() {
	ctx.Module.contextFound(ctx)
}
//返回错误
func (ctx *HttpContext) Error(err *Error) {
	ctx.Wrong = err
	ctx.Module.contextError(ctx)
}

//失败, 就是参数处理失败为主
func (ctx *HttpContext) Failed(err *Error) {
	ctx.Wrong = err
	ctx.Module.contextFailed(ctx)
}
//拒绝,主要是 auth
func (ctx *HttpContext) Denied(err *Error) {
	ctx.Wrong = err
	ctx.Module.contextFailed(ctx)
}
/* 上下文处理器 end */










/* 上下文响应器 begin */
func (ctx *HttpContext) Goto(url string) {
	ctx.Body = httpBodyGoto{url}
}
func (ctx *HttpContext) Text(text string, codes ...int) {
	if len(codes) > 0 {
		ctx.Code = codes[0]
	}
	ctx.Type = "text"
	ctx.Body = httpBodyText{text}
}
func (ctx *HttpContext) Html(html string, codes ...int) {
	if len(codes) > 0 {
		ctx.Code = codes[0]
	}
	ctx.Type = "html"
	ctx.Body = httpBodyHtml{html}
}
func (ctx *HttpContext) Script(script string, codes ...int) {
	if len(codes) > 0 {
		ctx.Code = codes[0]
	}
	ctx.Type = "script"
	ctx.Body = httpBodyScript{script}
}
func (ctx *HttpContext) Json(json Any, codes ...int) {
	if len(codes) > 0 {
		ctx.Code = codes[0]
	}
	ctx.Type = "json"
	ctx.Body = httpBodyJson{json}
}
func (ctx *HttpContext) Xml(xml Any, codes ...int) {
	if len(codes) > 0 {
		ctx.Code = codes[0]
	}
	ctx.Type = "xml"
	ctx.Body = httpBodyXml{xml}
}
func (ctx *HttpContext) File(file string, names ...string) {
	name := ""
	if len(names) > 0 {
		name = names[0]
	}

	ctx.Code = 200
	ctx.Type = "file"
	ctx.Body = httpBodyFile{file,name}
}
func (ctx *HttpContext) Down(body []byte, mimeType string, names ...string) {
	name := ""
	if len(names) > 0 {
		name = names[0]
	}

	ctx.Code = 200
	ctx.Type = mimeType
	ctx.Body = httpBodyDown{body,name}
}
func (ctx *HttpContext) View(view string, models ...Map) {
	model := Map{};
	if len(models) > 0 {
		model = models[0]
	}

	ctx.Type = "html"
	ctx.Body = httpBodyView{view, model}
}
/* 上下文响应器 end */











/*
	HTTP上下文方法 end
*/





//通用方法
func (ctx *HttpContext) Ip() string {
	ip := "127.0.0.1"

	if realIp := ctx.Req.Header.Get("X-Real-IP"); realIp != "" {
		ip = realIp
	} else if forwarded := ctx.Req.Header.Get("x-forwarded-for"); forwarded != "" {
		ip = forwarded
	} else {
		//GO默认带端口,要去掉端口
		//如果是IPV6,应该不要去,这个以后要判断处理
		ip := ctx.Req.RemoteAddr
		pos := strings.Index(ip, ":")
		if (pos >= 0) {
			ip = ip[0:pos]
		}

	}
	return ip
}





/*
	HTTP上下文方法 end
*/











