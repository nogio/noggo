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
	"time"
)


// http driver begin

type (

	//HTTP驱动
	HttpDriver interface {
		Connect(config Map) (HttpConnect,error)
	}

	HttpHandler func(req *http.Request, res http.ResponseWriter)

	//HTTP连接
	HttpConnect interface {
		//打开驱动连接
		Open() error
		//关闭驱动连接
		Close() error

		//注册回调？
		Accept(HttpHandler) error

		//开始
		Start(addr string) error
		//开始TLS
		StartTLS(addr string, certFile, keyFile string) error

	}
)


// http driver end


type (

	//HTTP全局容器
	httpGlobal	struct {
		mutex sync.Mutex

		//驱动
		drivers map[string]HttpDriver
		//中间件
		middlers    map[string]HttpFunc
		middlerNames []string

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
func (module *httpGlobal) connect(config *httpConfig) (HttpConnect,error) {
	if httpDriver,ok := module.drivers[config.Driver]; ok {
		return httpDriver.Connect(config.Config)
	} else {
		panic("HTTP：不支持的驱动 " + config.Driver)
	}
}


//注册HTTP驱动
func (global *httpGlobal) Driver(name string, config HttpDriver) {
	global.mutex.Lock()
	defer global.mutex.Unlock()

	if config == nil {
		panic("HTTP: 驱动不可为空")
	}
	//不做存在判断，因为要支持后注册的驱动替换已注册的驱动
	//框架有可能自带几种默认驱动，并且是默认注册的，用户可以自行注册替换
	global.drivers[name] = config
}

func (global *httpGlobal) Middler(name string, call HttpFunc) {
	global.mutex.Lock()
	defer global.mutex.Unlock()

	//保存配置
	if _,ok := global.middlers[name]; ok == false {
		//没有注册过name，才把name加到列表
		global.middlerNames = append(global.middlerNames, name)
	}
	//可以后注册重写原有路由配置，所以直接保存
	global.middlers[name] = call
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

		//View配置与连接
		viewConfig	*viewConfig
		viewConnect	ViewConnect

		//HTTP配置与连接
		httpConfig	*httpConfig
		httpConnect	HttpConnect


		//所在节点
		node	*Noggo


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
		Url *httpUrl

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
		Debug   bool        //当前请求是否调试模式

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
		Raw     string
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

		//错误信息，error,error,denied上下文中使用
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
	httpBodyJsonp struct {
		//响应的Json对象
		Json Any
		Callback string
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
		Model Any
	}
)



//HTTP模块初始化
func (module *httpModule) run() {
	module.runSession()
	module.runView()
	module.runHttp()
}
func (module *httpModule) runSession() {
	//使用节点会话
	module.sessionConfig = module.node.session.sessionConfig
	module.sessionConnect = module.node.session.sessionConnect
}

func (module *httpModule) runView() {

	module.viewConfig = Config.View
	con,err := View.connect(module.viewConfig)

	if err != nil {
		panic("节点HTTP：连接View失败：" + err.Error())
	} else {

		//打开连接
		err := con.Open()
		if err != nil {
			panic("节点HTTP：打开View失败 " + err.Error())
		} else {

			//保存
			module.viewConnect = con

		}

	}

}
func (module *httpModule) runHttp() {

	module.httpConfig = Config.Http
	con,err := Http.connect(module.httpConfig)

	if err != nil {
		panic("节点HTTP：连接失败：" + err.Error())
	} else {

		//打开会话连接
		err := con.Open()
		if err != nil {
			panic("节点HTTP：打开失败 " + err.Error())
		}

		//注册回调
		con.Accept(module.serveHttp)

		//开始http
		con.Start(module.node.Port)

		//保存连接
		module.httpConnect = con

	}


}


//HTTP模块退出
func (module *httpModule) end() {
	module.endSession()
	module.endView()
	module.endHttp()
}
//退出SESSION
func (module *httpModule) endSession() {
	//使用点节会话，这里不用处理
}
//退出view
func (module *httpModule) endView() {
	if module.viewConnect != nil {
		module.viewConnect.Close()
		module.viewConnect = nil
	}
}
//退出HTTP本身
func (module *httpModule) endHttp() {
	if module.httpConnect != nil {
		module.httpConnect.Close()
		module.httpConnect = nil
	}
}








//HTTP：注册路由
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


//HTTP：获取节点所有路由配置
func (module *httpModule) Routes() (Map) {
	m := Map{}

	for k,v := range module.routes {
		m[k] = v
	}

	return m
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
		node: node,
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


	//节点 拒绝处理器
	nodeDeniedHandlers, nodeDeniedHandlersOK := Http.deniedHandlers[node.Name];
	nodeDeniedHandlerNames, nodeDeniedHandlerNamesOK := Http.deniedHandlerNames[node.Name];
	if nodeDeniedHandlersOK && nodeDeniedHandlerNamesOK {
		for _,n := range nodeDeniedHandlerNames {
			module.DeniedHandler(n, nodeDeniedHandlers[n])
		}
	}
	//全局 拒绝处理器
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
func (module *httpModule) newHttpContext(req *http.Request, res http.ResponseWriter) (*HttpContext) {
	method := strings.ToLower(req.Method)
	host := req.URL.Host
	path := req.URL.Path
	value := Map{}

	return &HttpContext{
		Node: module.node, Module: module,
		next: -1, nexts: []HttpFunc{},

		Req: req, Res: res,

		Method: method, Host: host, Path: path,
		Ajax: false, Lang: "default",

		Name: "", Config: Map{}, Branchs:[]Map{},

		Param: Map{}, Query: Map{}, Form: Map{}, Upload: Map{}, Client: Map{},
		Value: value, Local: Map{}, Item: Map{}, Auth: Map{}, Args: Map{},

		Data: Map{},
	}
}



//HTTPHttp  请求开始
func (module *httpModule) serveHttp(req *http.Request, res http.ResponseWriter) {
	ctx := module.newHttpContext(req, res)
	ctx.Url = &httpUrl{ctx}

	//请求处理
	ctx.handler(module.contextRoute)
	ctx.handler(module.contextRequest)
	//最终所有的响应处理，优先
	ctx.handler(module.contextResponse)

	//中间件
	//用数组保证原始注册顺序
	for _,name := range Http.middlerNames {
		ctx.handler(Http.middlers[name])
	}

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

		//如果uri以/结尾就去掉
		if uri != "/" && strings.HasSuffix(uri, "/") {
			uri = uri[0:len(uri)-1]
		}

		keys, vals := []string{}, []string{}

		//用来匹配的正则表达式
		regs := strings.Replace(uri, ".", `\.`, -1)
		regs = strings.Replace(regs, "/", `\/`, -1)
		if uri != "/" {
			regs = regs + `(|\/)`   //加上可以以斜杠结尾
		}

		//正则对象
		regx := regexp.MustCompile(`(\{[\*A-Za-z0-9_]+\})`)
		//拿到URI中的参数列表
		//kkkkks := regx.FindAllString(uri, -1)

		//替换参数为正则
		regs = regx.ReplaceAllStringFunc(regs, func(p string) string {
			//if (p[1:1] == "*") {
			if (strings.HasPrefix(p, "{*")) {
				//{*abcd}
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
			//因为结尾加了或斜杠，vals的length不一样了
			//if len(keys) == len(vals) {
				for i,k := range keys {
					k = strings.Replace(k, "*", "", -1)
					param[k] = vals[i]
				}
			//}

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
	cookie, err := ctx.Req.Cookie(module.node.Config.Cookie)
	if err != nil || cookie.Value == "" {
		ctx.Id = NewMd5Id()
		m,err := module.sessionConnect.Query(ctx.Id)
		if err == nil {
			ctx.Session = m
		} else {
			ctx.Session = Map{}
		}

		//注意域设置
		//这里在ctx.Cookie封装之后，应当写到ctx.Cookie中，在response写入res
		//待修改
		cookie := http.Cookie{ Name: module.node.Config.Cookie, Value: url.QueryEscape(ctx.Id), Path: "/", HttpOnly: true }
		if module.node.Config.Expiry > 0 {
			cookie.MaxAge = int(module.node.Config.Expiry)
		}
		if module.node.Config.Domain != "" {
			cookie.Domain = module.node.Config.Domain
		}
		http.SetCookie(ctx.Res, &cookie)
	} else {
		ctx.Id, _ = url.QueryUnescape(cookie.Value)

		m,err := module.sessionConnect.Query(ctx.Id)
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
	//复制params
	for k,v := range ctx.Param {
		ctx.Value[k] = v
	}

	ctx.Sign = &Sign{ ctx.Session }
	ctx.Next()
	//module.sessionConnect.Update(ctx.Id, ctx.Session, module.sessionConfig.Expiry)
	module.sessionConnect.Update(ctx.Id, ctx.Session)
}





//处理响应
func (module *httpModule) contextResponse(ctx *HttpContext) {

	//默认一个500，这样才正常
	ctx.Code = 500

	//因为response是在所有请求前的， 所以先调用一下
	//然后对结果进行处理
	ctx.Next()


	//清理执行线
	ctx.cleanup()

	//filter中的response
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
		/*
		//在执行这里不处理item，注意，把item汇总，在真正的action以前执行
		if _,ok := ctx.Config[KeyMapItem]; ok {
			ctx.handler(module.contextItem)
		}
		*/

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
	//复制item，因为item放到action前执行了
	//这里拿到的是全局的item
	var globalItem Map
	if iv,ok := ctx.Config[KeyMapItem].(Map); ok {
		globalItem = iv
	}


	//这里 ctx.Route 和 routing 变换位置
	ctx.Config = Map{}

	//如果有路由
	//HTTP路由不支持多method，非http
	if routeConfig,ok := routing[KeyMapRoute]; ok {

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
			} else {
				ctx.Config = nil
			}
		}
	} else {
		ctx.Config = nil
	}

	if ctx.Config == nil {
		//还是不存在
		ctx.Code = 404
		ctx.handler(module.contextFound)
	} else {

		//复制全局的item
		if globalItem != nil {

			newItems := Map{}
			if iv,ok := ctx.Config[KeyMapItem].(Map); ok {
				for k,v := range iv {
					newItems[k] = v
				}
			}
			for k,v := range globalItem {
				newItems[k] = v
			}

			//写过来
			ctx.Config[KeyMapItem] = newItems
		}



		//走到这了， 肯定 是200了
		ctx.Code = 200

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
	err := Mapping.Parse([]string{}, ctx.Config["args"].(Map), ctx.Value, ctx.Args, argn, false)
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

				//判断是否需要查询数据
				if baseName,ok := authConfig["base"].(string); ok {
					if tableName, ok := authConfig["table"].(string); ok {
						//要查询表
						//不管must是否,都要查表
						db := Data.Base(baseName);
						item,err := db.Table(tableName).Entity(ctx.Sign.Id(authSign))
						db.Close()
						if err != nil {
							if authMust {	//是必要的
								//是否有自定义状态
								err := Const.NewTypeLangStateError(authKey, ctx.Lang, "auth.error", authName)
								if v,ok := authConfig["error"]; ok {
									err = Const.NewTypeLangStateError(authKey, ctx.Lang, v.(string))
								}

								ctx.Denied(err)
								return
							}
						} else {
							saveMap[authKey] = item
						}

					} else if viewName, ok := authConfig["view"].(string); ok {
						//要查询表
						//不管must是否,都要查表
						db := Data.Base(baseName);
						item,err := db.View(viewName).Entity(ctx.Sign.Id(authSign))
						db.Close()
						if err != nil {
							if authMust {	//是必要的
								//是否有自定义状态
								err := Const.NewTypeLangStateError(authKey, ctx.Lang, "auth.error", authName)
								if v,ok := authConfig["error"]; ok {
									err = Const.NewTypeLangStateError(authKey, ctx.Lang, v.(string))
								}

								ctx.Denied(err)
								return
							}
						} else {
							saveMap[authKey] = item
						}

					} else if modelName, ok := authConfig["model"].(string); ok {
						//兼容老代码使用model->TABLE

						//要查询表
						//不管must是否,都要查表
						db := Data.Base(baseName);
						item,err := db.Table(modelName).Entity(ctx.Sign.Id(authSign))
						db.Close()
						if err != nil {
							if authMust {	//是必要的
								//是否有自定义状态
								err := Const.NewTypeLangStateError(authKey, ctx.Lang, "auth.error", authName)
								if v,ok := authConfig["error"]; ok {
									err = Const.NewTypeLangStateError(authKey, ctx.Lang, v.(string))
								}

								ctx.Denied(err)
								return;
							}
						} else {
							saveMap[authKey] = item
						}

					}


				}


			} else {
				ohNo = true
			}

			//到这里是未登录的
			//而且是必须要登录，才显示错误
			if ohNo && authMust {

				//是否有自定义状态
				err := Const.NewTypeLangStateError(authKey, ctx.Lang, "auth.empty", authName)
				if v,ok := authConfig["empty"]; ok {
					err = Const.NewTypeLangStateError(authKey, ctx.Lang, v.(string))
				}

				ctx.Denied(err)
				return

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

			//是否必须
			must := true
			if vv,ok := config["must"].(bool); ok {
				must = vv
			}


			name := config["name"].(string)
			key := k
			var val Any = nil
			if vv,ok := config["args"].(string); ok {
				key = vv
				val = ctx.Args[key]
			} else if vv,ok := config["param"].(string); ok {
				key = vv
				val = ctx.Param[key]
			} else if vv,ok := config["query"].(string); ok {
				key = vv
				val = ctx.Query[key]
			} else if vv,ok := config["value"].(string); ok {
				key = vv
				val = ctx.Value[key]
			} else if vv,ok := config["key"].(string); ok {
				key = vv
				val = ctx.Value[key]
			} else {
				val = nil
			}

			if val == nil && must {
				//参数不为空啊啊
				state := "item.empty"
				//是否有自定义状态
				if v,ok := config["empty"]; ok {
					state = v.(string)
				}
				err := Const.NewTypeLangStateError(k, ctx.Lang, state, name)
				//查询不到东西，也要失败， 接口访问失败
				ctx.Failed(err)
				return
			} else {

				//判断是否需要查询数据
				if baseName,ok := config["base"].(string); ok {



					if tableName,ok := config["table"].(string); ok {

						//要查询库
						db := Data.Base(baseName);
						item,err := db.Table(tableName).Entity(val)
						db.Close()
						if err != nil && must {
							state := "item.error"
							//是否有自定义状态
							if v,ok := config["error"]; ok {
								state = v.(string)
							}
							err := Const.NewTypeLangStateError(k, ctx.Lang, state, name)

							ctx.Failed(err)
							return;
						} else {
							saveMap[k] = item
						}
					} else if viewName,ok := config["view"].(string); ok {

						//要查询库
						db := Data.Base(baseName);
						item,err := db.View(viewName).Entity(val)
						db.Close()
						if err != nil && must {
							state := "item.error"
							//是否有自定义状态
							if v,ok := config["error"]; ok {
								state = v.(string)
							}
							err := Const.NewTypeLangStateError(k, ctx.Lang, state, name)

							ctx.Failed(err)
							return;
						} else {
							saveMap[k] = item
						}
					} else if modelName,ok := config["model"].(string); ok {
						//兼容老代码 model->table

						//要查询库
						db := Data.Base(baseName);
						item,err := db.Table(modelName).Entity(val)
						db.Close()
						if err != nil && must {
							state := "item.error"
							//是否有自定义状态
							if v,ok := config["error"]; ok {
								state = v.(string)
							}
							err := Const.NewTypeLangStateError(k, ctx.Lang, state, name)

							ctx.Failed(err)
							return;
						} else {
							saveMap[k] = item
						}
					}



				}

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
		ctx.Charset = module.node.Config.Charset
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
	case httpBodyJsonp:
		module.jsonpResponder(ctx)
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

	ctx.Code = 404
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

	ctx.Code = 500
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

	ctx.Code = 500

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

	ctx.Code = 500

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
		//这里应该转到error上下文处理
		//待修改
		http.Error(ctx.Res, err.Error(), 500)


	} else {

		if ctx.Type == "" {
			ctx.Type = "json"
		}

		ctx.Res.Header().Set("Content-Type", fmt.Sprintf("%v; charset=%v", Const.MimeType(ctx.Type), ctx.Charset))
		ctx.Res.WriteHeader(ctx.Code)
		fmt.Fprint(ctx.Res, string(bytes))
	}
}

func (module *httpModule) jsonpResponder(ctx *HttpContext) {
	body := ctx.Body.(httpBodyJsonp)

	bytes, err := json.Marshal(body.Json)
	if err != nil {
		//出错啊
		//但是这里已经走完response了。再ctx.Error好像没用了
		//这是一个死循环， 因为走了ctx.Error， 还有可能再返回json
		//又会走到这里，在response中。 继续把response加入调用列表吧。这样保险
		//这里应该转到error上下文处理
		//待修改
		http.Error(ctx.Res, err.Error(), 500)


	} else {

		if ctx.Type == "" {
			ctx.Type = "script"
		}

		resp := fmt.Sprintf("%v(%v)", body.Callback, string(bytes))

		ctx.Res.Header().Set("Content-Type", fmt.Sprintf("%v; charset=%v", Const.MimeType(ctx.Type), ctx.Charset))
		ctx.Res.WriteHeader(ctx.Code)
		fmt.Fprint(ctx.Res, resp)
	}
}
func (module *httpModule) xmlResponder(ctx *HttpContext) {
	body := ctx.Body.(httpBodyXml)

	bytes, err := xml.Marshal(body.Xml)
	if err != nil {
		//这里应该转到error上下文处理
		//待修改
		http.Error(ctx.Res, err.Error(), 500)

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


	if ctx.Type != "file" {
		ctx.Res.Header().Set("Content-Type", fmt.Sprintf("%v; charset=%v", Const.MimeType(ctx.Type), ctx.Charset))
	}

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
	ctx.Res.Write(body.Body)
	//fmt.Fprint(ctx.Res, body.Body)
}
func (module *httpModule) viewResponder(ctx *HttpContext) {
	if ctx.Type == "" {
		ctx.Type = "html"
	}

	body := ctx.Body.(httpBodyView)


	/*
		"agent": func() string{
			return ctx.Agent()
		},
		"ip": func() string{
			return ctx.Ip()
		},
		*/

	/*
	"storage": func(item Map) string {
		return ctx.Url.Storage(item)
	},
	"thumbnail": func(item Map, w,h,t int64) string {
		return ctx.Url.Thumbnail(item, w,h,t)
	},
	*/




	/*
	"status":  func(data,model string, value Any) template.HTML {
		html := ""

		if value == nil {
			html = `<span class="green">正常</span>`
		} else {
			enums := Data.Status(data, model)
			key := fmt.Sprintf("%v", value)
			if v, ok := enums[key]; ok {
				html = fmt.Sprintf(`<span class="red">%v</span>`, v)
			}
		}
		return template.HTML(html)
	},
*/


	parse := &ViewParse{
		Node: ctx.Node.Name, Lang: ctx.Lang,
		Data: ctx.Data, View: body.View, Model: body.Model,
		Args: ctx.Args, Auth: ctx.Auth,
		Setting: Setting, Session: ctx.Session,
		Helpers: Map{
			"backurl": func() string{ return ctx.Url.Back() },
			"lasturl": func() string{ return ctx.Url.Last() },
			"siteurl": func(name string, args ...string) string{ return ctx.Url.Site(name, args...) },
			"nodeurl": func(name string, args ...string) string{ return ctx.Url.Node(name, args...) },
			"rooturl": func(args ...string) string{ return ctx.Url.Root(args...) },
			"route": func(name string, vals ...Any) string{
				args := []Map{}
				for _,v := range vals {
					switch t := v.(type) {
					case Map:
						args = append(args, t)
					default:
						m := Map{}
						e := json.Unmarshal([]byte(fmt.Sprintf("%v", t)), &m)
						if e == nil {
							args = append(args, m)
						}
					}
				}
				return ctx.Url.Route(name, args...)
			},
			"routo": func(site,name string, vals ...Any) string {
				args := []Map{}
				for _,v := range vals {
					switch t := v.(type) {
					case Map:
						args = append(args, t)
					default:
						m := Map{}
						e := json.Unmarshal([]byte(fmt.Sprintf("%v", t)), &m)
						if e == nil {
							args = append(args, m)
						}
					}
				}
				return ctx.Url.Routo(site,name, args...)
			},
			"signed": func(key string) bool {
				return ctx.Sign.Yes(key)
			},
			"signid": func(key string) string {
				return fmt.Sprintf("%v", ctx.Sign.Id(key))
			},
			"signer": func(key string) Any {
				return ctx.Sign.Name(key)
			},
			//多语音字串，使用...Any，因为在view做类型转换比较麻烦，在这里转换
			"lang": func(key string, args ...Any) string {
				strs := []string{}
				for _,v := range args {
					strs = append(strs, fmt.Sprintf("%v", v))
				}
				return ctx.String(key, strs...)
			},
			"enum": func(data,model,field string,v Any) (string) {
				html := ""

				value := fmt.Sprintf("%v", v)

				enums := Data.Enums(data,model,field)
				if v,ok := enums[value]; ok {
					html = fmt.Sprintf("%v", v)
				}
				return html
			},
		},
	}

	//helper要在这里处理
	//可以注册Helper
	for k,v := range View.helpers {
		parse.Helpers[k] = v
	}

	html,err := module.viewConnect.Parse(parse)
	if err != nil {

		//这里应该转到error上下文处理
		//待修改
		http.Error(ctx.Res, err.Error(), 500)

	} else {
		ctx.Res.Header().Set("Content-Type", fmt.Sprintf("%v; charset=%v", Const.MimeType(ctx.Type), ctx.Charset))
		ctx.Res.WriteHeader(ctx.Code)
		fmt.Fprint(ctx.Res, html)
	}

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
//返回错误，args,item,data时都发生
func (ctx *HttpContext) Error(err *Error) {
	ctx.Wrong = err
	ctx.Module.contextError(ctx)
}
//失败时
func (ctx *HttpContext) Failed(err *Error) {
	ctx.Wrong = err
	ctx.Module.contextFailed(ctx)
}
//auth拒绝访问时
func (ctx *HttpContext) Denied(err *Error) {
	ctx.Wrong = err
	ctx.Module.contextDenied(ctx)
}
/* 上下文处理器 end */










/* 上下文响应器 begin */
func (ctx *HttpContext) Goto(url string) {
	ctx.Body = httpBodyGoto{url}
}
func (ctx *HttpContext) Text(text Any, codes ...int) {
	if len(codes) > 0 {
		ctx.Code = codes[0]
	}
	ctx.Type = "text"
	ctx.Body = httpBodyText{fmt.Sprintf("%v", text)}
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
func (ctx *HttpContext) Jsonp(json Any, cbs ...string) {
	cb := "callback"
	if vv,ok := ctx.Query["callback"].(string); ok && vv != "" {
		cb = vv
	}
	if len(cbs) > 0 && cbs[0] != "" {
		cb = cbs[0]
	}
	ctx.Type = "script"
	ctx.Body = httpBodyJsonp{json, cb}
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
func (ctx *HttpContext) TypeFile(mimeType string, file string, names ...string) {
	name := ""
	if len(names) > 0 {
		name = names[0]
	}

	ctx.Code = 200
	ctx.Type = mimeType
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
func (ctx *HttpContext) View(view string, models ...Any) {
	var model Any
	if len(models) > 0 {
		model = models[0]
	}

	ctx.Type = "html"
	ctx.Body = httpBodyView{view, model}
}
func (ctx *HttpContext) TypeView(tttt string, view string, models ...Any) {
	var model Any
	if len(models) > 0 {
		model = models[0]
	}

	if tttt == "" {
		tttt = "html"
	}

	ctx.Type = tttt
	ctx.Body = httpBodyView{view, model}
}
/* 上下文响应器 end */


/* 上下文响应器 end */










func (ctx *HttpContext) Goback() {
	ctx.Goto(ctx.Url.Back())
}
func (ctx *HttpContext) Golast() {
	ctx.Goto(ctx.Url.Last())
}
//跳转到路由
func (ctx *HttpContext) Route(name string, args ...Map) {
	ctx.Goto(ctx.Url.Route(name, args...))
}
//跳转到路由
func (ctx *HttpContext) Routo(site,name string, args ...Map) {
	ctx.Goto(ctx.Url.Routo(site, name, args...))
}



func (ctx *HttpContext) Redirect(url string) {
	ctx.Goto(url)
}

func (ctx *HttpContext) Alert(text string, urls ...string) {

	e := Const.NewLangStateError(ctx.Lang, text)
	if e != nil {
		text = e.Text
	}

	if len(urls) > 0 {
		text = fmt.Sprintf(`<script type="text/javascript">alert("%s"); location.href="%s";</script>`, text, urls[0])
	} else {
		text = fmt.Sprintf(`<script type="text/javascript">alert("%s"); history.back();</script>`, text)
	}
	ctx.Html(text)
}
//展示通用的提示页面
func (ctx *HttpContext) Show(tttt, text string, urls ...string) {

	m := Map{
		"type": tttt,
		"text": text,
		"url": "",
	}
	e := Const.NewLangStateError(ctx.Lang, text)
	if e != nil {
		m["text"] = e.Text
	}
	if len(urls) > 0 {
		m["url"] = urls[0]
	}

	ctx.View("show", m)
}
/* 上下文响应器 end */







//专为接口准备的方法

//返回一个状态，表示失败, 无data节点
func (ctx *HttpContext) State(state string, args ...interface{}) {
	e := Const.NewLangStateError(ctx.Lang, state, args...)
	m := Map{
		"code": e.Code,
		"text": e.Text,
		"time": time.Now().Unix(),
	}
	ctx.Json(m, 500)
}


//返回操作结果，表示成功
//比如，登录，修改密码，等操作类的接口， 成功的时候，使用这个，
//args表示返回给客户端的data
//data 强制改为json格式，因为data有统一加密的可能
//所有数组都要加密。
func (ctx *HttpContext) Result(state string, args ...Map) {
	e := Const.NewLangStateError(ctx.Lang, state)

	//默认应当给Result表示成功，这样state直接写文本，而不用定义状态
	//当code有自定义的时候，才是使用自定义的code
	code := 0
	if e.Code != -1 {
		code = e.Code
	}
	m := Map{
		"code": code, //这里不强制设为0，可以自己在 mains/consts/state.go 中定义自己的状态
		"text": e.Text,
		"time": time.Now().Unix(),
	}
	if len(args) > 0 {
		//这里要对返回的内容对处理
		//待处理一项：data不一定非要Map类型，也可以是数组，或是其它对象
		//暂不处理，后续再改
		data := args[0]

		//如果需要对结果进行处理
		c,cok := ctx.Config["data"].(Map);

		if cok {

			//处理,需不需要整个data节点全加密
			if ctx.Node.Config.Crypto != "" && ctx.Config["nocode"]==nil && ctx.Debug == false {


				newConfig := Map{
					"data": Map{
						"type": "json", "must": true, "encode": ctx.Node.Config.Crypto,
						"json": c,
					},
				}
				newData := Map{
					"data": data,
				}


				v := Map{}

				e := Mapping.Parse([]string{}, newConfig, newData, v, false, false)
				if e != nil {
					//出错了
					ctx.Failed(e)
					return
				} else {
					//处理后的data
					m["data"] = v["data"]
				}


			} else {

				v := Map{}

				e := Mapping.Parse([]string{}, c, data, v, false, false)
				if e != nil {
					//出错了
					ctx.Failed(e)
					return
				} else {
					//处理后的data
					m["data"] = v
				}

			}

		} else {
			//不需要包装值，但是有加密的需求

			//处理,需不需要整个data节点全加密
			if ctx.Node.Config.Crypto != "" && ctx.Config["nocode"]==nil && ctx.Debug == false {


				newConfig := Map{
					"data": Map{
						"type": "json", "must": true, "encode": ctx.Node.Config.Crypto,
					},
				}
				newData := Map{
					"data": data,
				}


				v := Map{}

				e := Mapping.Parse([]string{}, newConfig, newData, v, false, false)
				if e != nil {
					//出错了
					ctx.Failed(e)
					return
				} else {
					//处理后的data
					m["data"] = v["data"]
				}


			} else {
				m["data"] = data
			}
		}

	}

	ctx.Json(m, 200)
}


//返回数据，表示成功
//data必须为json，因为data节点可能统一加密
//如果在data同级返回其它数据，如page信息， 会有泄露数据风险
//所以这里强制data必须为json
func (ctx *HttpContext) Return(data Map) {
	m := Map{
		"code": 0,
		"time": time.Now().Unix(),
	}

	//如果需要对结果进行处理
	//待处理一项：data不一定非要Map类型，也可以是数组，或是其它对象
	//暂不处理，后续再改
	c,cok := ctx.Config["data"].(Map);

	if cok {

		//处理,需不需要整个data节点全加密
		if ctx.Node.Config.Crypto != "" && ctx.Config["nocode"]==nil && ctx.Debug == false {


			newConfig := Map{
				"data": Map{
					"type": "json", "must": true, "encode": ctx.Node.Config.Crypto,
					"json": c,
				},
			}
			newData := Map{
				"data": data,
			}


			v := Map{}

			e := Mapping.Parse([]string{}, newConfig, newData, v, false, false)
			if e != nil {
				//出错了
				ctx.Failed(e)
				return
			} else {
				//处理后的data
				m["data"] = v["data"]
			}


		} else {

			v := Map{}

			e := Mapping.Parse([]string{}, c, data, v, false, false)
			if e != nil {
				//出错了
				ctx.Failed(e)
				return
			} else {
				//处理后的data
				m["data"] = v
			}

		}

	} else {
		//不需要包装值，但是有加密的需求

		//处理,需不需要整个data节点全加密
		if ctx.Node.Config.Crypto != "" && ctx.Config["nocode"]==nil && ctx.Debug == false {

			newConfig := Map{
				"data": Map{
					"type": "json", "must": true, "encode": ctx.Node.Config.Crypto,
				},
			}
			newData := Map{
				"data": data,
			}


			v := Map{}

			e := Mapping.Parse([]string{}, newConfig, newData, v, false, false)
			if e != nil {
				//出错了，用failed方法，因为这和接口有关，好在failedHandler中处理
				ctx.Failed(e)
				return
			} else {
				//处理后的data
				m["data"] = v["data"]
			}


		} else {

			//不需要包装值
			m["data"] = data
		}

	}


	ctx.Json(m, 200)
}


//编码中，主要为接口设计，用于自动生成接口测试数据
func (ctx *HttpContext) Coding() {
	ctx.State("coding")
}














/*
	HTTP上下文方法 end
*/



//获取langString
func (ctx *HttpContext) String(key string, args ...string) string {
	return Const.LangString(key, args...)
}



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

//通用方法
func (ctx *HttpContext) Header(key string, vals ...string) string {

	if len(vals) > 0 {
		//设置header
		ctx.Res.Header().Set(key, vals[0])
		return vals[0]
	} else {
		//读header
		return ctx.Req.Header.Get(key)
	}
	return ""
}

//通用方法
func (ctx *HttpContext) Cookie(key string, vals ...Any) string {

	if len(vals) > 0 {
		//设置header
		switch val := vals[0].(type) {
			case http.Cookie:
				val.Value = url.QueryEscape(val.Value)
				if ctx.Node.Config.Expiry > 0 {
					val.MaxAge = int(ctx.Node.Config.Expiry)
				}
				if ctx.Node.Config.Domain != "" {
					val.Domain = ctx.Node.Config.Domain
				}
				http.SetCookie(ctx.Res, &val)
			case string:
				cookie := http.Cookie{ Name: key, Value: url.QueryEscape(val), Path: "/", HttpOnly: true }
				if ctx.Node.Config.Expiry > 0 {
					cookie.MaxAge = int(ctx.Node.Config.Expiry)
				}
				if ctx.Node.Config.Domain != "" {
					cookie.Domain = ctx.Node.Config.Domain
				}
				http.SetCookie(ctx.Res, &cookie)
			default:
				return ""
		}
	} else {
		//读cookie
		c,e := ctx.Req.Cookie(key)
		if e == nil {
			return c.Value
		}
	}
	return ""
}






/*
	HTTP上下文方法 end
*/












