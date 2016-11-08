/*
	plan	计划模块
	计划功能是一个周期性的功能，定时周期执行
	主要使用场景：定时提醒，备份啊等等

	具体的计划执行代码，都在节点中，而非全局
*/

package noggo

import (
	. "github.com/nogio/noggo/base"
	"sync"
	"time"
)


// plan driver begin


type (
	//计划驱动
	PlanDriver interface {
		Connect(config Map) (PlanConnect,error)
	}

	PlanHandler func(*PlanRequest, PlanResponse)

	//计划连接
	PlanConnect interface {
		//打开连接
		Open() error
		//关闭连接
		Close() error

		//注册回调
		Accept(PlanHandler) error

		//注册计划
		Register(name string, time string) error

		//开始计划
		Start() error
	}


	//计划请求实体
	PlanRequest struct {
		Id string
		Name string
		Time string
		Value Map
	}

	//计划响应接口
	PlanResponse interface {
		//完成
		Finish(id string) error
		//重新开始
		Replan(id string, delay time.Duration) error
	}
)

// plan driver end



type (
	//计划全局容器
	planGlobal	struct {
		mutex sync.Mutex
		//驱动
		drivers map[string]PlanDriver

		//中间件
		middlers    map[string]PlanFunc
		middlerNames []string

		//路由
		routes 		map[string]map[string]Map	//路由定义							map[node]map[name]Map
		routeNames	map[string][]string			//路由名称原始顺序，因为map是无序的		map[node][]string

		//拦截器们
		requestFilters, executeFilters, responseFilters map[string]map[string]PlanFunc
		requestFilterNames, executeFilterNames, responseFilterNames map[string][]string

		//处理器们
		foundHandlers, errorHandlers map[string]map[string]PlanFunc
		foundHandlerNames, errorHandlerNames map[string][]string
	}

)

//计划：连接驱动
func (module *planGlobal) connect(config *planConfig) (PlanConnect,error) {
	if planDriver,ok := module.drivers[config.Driver]; ok {
		return planDriver.Connect(config.Config)
	} else {
		panic("计划：不支持的驱动 " + config.Driver)
	}
}


//注册计划驱动
func (global *planGlobal) Driver(name string, config PlanDriver) {
	global.mutex.Lock()
	defer global.mutex.Unlock()

	if global.drivers == nil {
		global.drivers = map[string]PlanDriver{}
	}

	if config == nil {
		panic("计划: 驱动不可为空")
	}
	//不做存在判断，因为要支持后注册的驱动替换已注册的驱动
	//框架有可能自带几种默认驱动，并且是默认注册的，用户可以自行注册替换
	global.drivers[name] = config
}


func (global *planGlobal) Middler(name string, call PlanFunc) {
	global.mutex.Lock()
	defer global.mutex.Unlock()


	if global.middlers == nil {
		global.middlers = map[string]PlanFunc{}
	}
	if global.middlerNames == nil {
		global.middlerNames = []string{}
	}

	//保存配置
	if _,ok := global.middlers[name]; ok == false {
		//没有注册过name，才把name加到列表
		global.middlerNames = append(global.middlerNames, name)
	}
	//可以后注册重写原有路由配置，所以直接保存
	global.middlers[name] = call
}





//计划全局初始化
func (global *planGlobal) init() {
	//计划全局无需处理任何东西
}

//计划全局退出
func (global *planGlobal) exit() {
	//计划全局无需处理任何东西
}








//计划：注册路由
//注册路由到全局容器
//Current将标明表示属于哪一个节点
//如果Current为空，表示全局，相当于注册到所有节点
func (global *planGlobal) Route(name string, config Map) {
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
func (global *planGlobal) RequestFilter(name string, call PlanFunc) {
	global.mutex.Lock()
	defer global.mutex.Unlock()

	if global.requestFilters == nil {
		global.requestFilters = map[string]map[string]PlanFunc{}
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
		global.requestFilters[nodeName] = map[string]PlanFunc{}
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
func (global *planGlobal) ExecuteFilter(name string, call PlanFunc) {
	global.mutex.Lock()
	defer global.mutex.Unlock()

	if global.executeFilters == nil {
		global.executeFilters = map[string]map[string]PlanFunc{}
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
		global.executeFilters[nodeName] = map[string]PlanFunc{}
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
func (global *planGlobal) ResponseFilter(name string, call PlanFunc) {
	global.mutex.Lock()
	defer global.mutex.Unlock()

	if global.responseFilters == nil {
		global.responseFilters = map[string]map[string]PlanFunc{}
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
		global.responseFilters[nodeName] = map[string]PlanFunc{}
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
func (global *planGlobal) FoundHandler(name string, call PlanFunc) {
	global.mutex.Lock()
	defer global.mutex.Unlock()

	if global.foundHandlers == nil {
		global.foundHandlers = map[string]map[string]PlanFunc{}
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
		global.foundHandlers[nodeName] = map[string]PlanFunc{}
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
func (global *planGlobal) ErrorHandler(name string, call PlanFunc) {
	global.mutex.Lock()
	defer global.mutex.Unlock()

	if global.errorHandlers == nil {
		global.errorHandlers = map[string]map[string]PlanFunc{}
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
		global.errorHandlers[nodeName] = map[string]PlanFunc{}
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


//-----------------------------------------------------------------------------------------------------------------------//












type (

	//计划上下文函数类型
	PlanFunc func(*PlanContext)

	//响应完成
	planBodyFinish struct {
	}
	//响应重新触发
	planBodyReplan struct {
		Delay time.Duration
	}

	//计划模块
	planModule struct {
		mutex sync.Mutex

		//会话配置与连接
		sessionConfig	*sessionConfig
		sessionConnect	SessionConnect

		//计划配置与连接
		planConfig	*planConfig
		planConnect	PlanConnect


		//所在节点
		node	*Noggo


		//路由
		routes 		map[string]Map			//路由定义
		routeNames	[]string				//路由名称原始顺序，因为map是无序的
		routeTimes	map[string]string		//计划的时间，一个计划可以多个时间。map[name]time

		//拦截器们
		requestFilters, executeFilters, responseFilters map[string]PlanFunc
		requestFilterNames, executeFilterNames, responseFilterNames []string

		//处理器们
		foundHandlers, errorHandlers map[string]PlanFunc
		foundHandlerNames, errorHandlerNames []string
	}

	//计划上下文
	PlanContext struct {
		Node	*Noggo
		Module	*planModule

		//执行线
		nexts []PlanFunc		//方法列表
		next int				//下一个索引

		req *PlanRequest
		res PlanResponse

		//基础
		Id	string			//Session Id  会话时使用
		Session Map			//存储Session值
		Sign	*Sign		//签名功能，基于session

		//配置相关
		Name string			//解析路由后得到的name
		Config Map			//解析后得到的路由配置
		Branchs []Map		//解析后得到的路由分支配置

		//此计划的时间
		Time    string

		//数据相关
		Value	Map			//所有请求过来的原始参数汇总
		Local	Map			//在ctx中传递数据用的
		Item	Map			//单条记录查询对象
		Auth	Map			//签名认证对象
		Args	Map			//经过args处理后的参数

		//响应相关
		Body	Any			//响应内容

		Wrong	*Error		//错误信息
	}
)



//计划模块初始化
func (module *planModule) run() {
	module.runSession()
	module.runPlan()
}
func (module *planModule) runSession() {
	//使用节点的会话了
	module.sessionConfig = module.node.session.sessionConfig
	module.sessionConnect = module.node.session.sessionConnect
}
func (module *planModule) runPlan() {

	module.planConfig = Config.Plan
	con,err := Plan.connect(module.planConfig)


	if err != nil {
		panic("节点计划：连接失败：" + err.Error())
	} else {

		//打开会话连接
		err := con.Open()
		if err != nil {
			panic("节点计划：打开失败 " + err.Error())
		} else {

			//注册回调
			con.Accept(module.servePlan)

			//注册计划
			for name,time := range module.routeTimes {
				con.Register(name, time)
			}

			con.Start()

			//保存
			module.planConnect = con

		}
	}
}


//计划模块退出
func (module *planModule) end() {
	module.endSession()
	module.endPlan()

}
//退出SESSION
func (module *planModule) endSession() {
	//使用节点的会话，此处无需处理
}
//退出计划本身
func (module *planModule) endPlan() {
	if module.planConnect != nil {
		module.planConnect.Close()
	}
}








//任务：注册路由
func (module *planModule) Route(name string, config Map) {
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





	if module.routeTimes == nil {
		module.routeTimes = map[string]string{}
	}


	//处理time
	if v,ok := config[KeyMapTime]; ok {
		switch ttt := v.(type) {
		case string:
			module.routeTimes[name] = ttt
		}
	}
}








/* 注册拦截器 begin */
func (module *planModule) RequestFilter(name string, call PlanFunc) {
	module.mutex.Lock()
	defer module.mutex.Unlock()

	if module.requestFilters == nil {
		module.requestFilters = make(map[string]PlanFunc)
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
func (module *planModule) ExecuteFilter(name string, call PlanFunc) {
	module.mutex.Lock()
	defer module.mutex.Unlock()

	if module.executeFilters == nil {
		module.executeFilters = make(map[string]PlanFunc)
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
func (module *planModule) ResponseFilter(name string, call PlanFunc) {
	module.mutex.Lock()
	defer module.mutex.Unlock()

	if module.responseFilters == nil {
		module.responseFilters = make(map[string]PlanFunc)
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
func (module *planModule) FoundHandler(name string, call PlanFunc) {
	module.mutex.Lock()
	defer module.mutex.Unlock()

	if module.foundHandlers == nil {
		module.foundHandlers = make(map[string]PlanFunc)
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
func (module *planModule) ErrorHandler(name string, call PlanFunc) {
	module.mutex.Lock()
	defer module.mutex.Unlock()

	if module.errorHandlers == nil {
		module.errorHandlers = make(map[string]PlanFunc)
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
/* 注册处理器 end */








//创建计划模块
func newPlanModule(node *Noggo) (*planModule) {
	module := &planModule{
		node: node,
	}


	//复制路由，拦截器，处理器


	//全局路由
	routes, routesOK := Plan.routes[ConstNodeGlobal];
	routeNames, routeNamesOK := Plan.routeNames[ConstNodeGlobal];
	if routesOK && routeNamesOK {
		for _,n := range routeNames {
			module.Route(n, routes[n])
		}
	}
	//节点路由
	nodeRoutes, nodeRoutesOK := Plan.routes[node.Name];
	nodeRouteNames, nodeRouteNamesOK := Plan.routeNames[node.Name];
	if nodeRoutesOK && nodeRouteNamesOK {
		for _,n := range nodeRouteNames {
			module.Route(n, nodeRoutes[n])
		}
	}



	//全局 请求拦截器
	requestFilters, requestFiltersOK := Plan.requestFilters[ConstNodeGlobal];
	requestFilterNames, requestFilterNamesOK := Plan.requestFilterNames[ConstNodeGlobal];
	if requestFiltersOK && requestFilterNamesOK {
		for _,n := range requestFilterNames {
			module.RequestFilter(n, requestFilters[n])
		}
	}
	//节点 请求拦截器
	nodeRequestFilters, nodeRequestFiltersOK := Plan.requestFilters[node.Name];
	nodeRequestFilterNames, nodeRequestFilterNamesOK := Plan.requestFilterNames[node.Name];
	if nodeRequestFiltersOK && nodeRequestFilterNamesOK {
		for _,n := range nodeRequestFilterNames {
			module.RequestFilter(n, nodeRequestFilters[n])
		}
	}
	//全局 执行拦截器
	executeFilters, executeFiltersOK := Plan.executeFilters[ConstNodeGlobal];
	executeFilterNames, executeFilterNamesOK := Plan.executeFilterNames[ConstNodeGlobal];
	if executeFiltersOK && executeFilterNamesOK {
		for _,n := range executeFilterNames {
			module.ExecuteFilter(n, executeFilters[n])
		}
	}
	//节点 执行拦截器
	nodeExecuteFilters, nodeExecuteFiltersOK := Plan.executeFilters[node.Name];
	nodeExecuteFilterNames, nodeExecuteFilterNamesOK := Plan.executeFilterNames[node.Name];
	if nodeExecuteFiltersOK && nodeExecuteFilterNamesOK {
		for _,n := range nodeExecuteFilterNames {
			module.ExecuteFilter(n, nodeExecuteFilters[n])
		}
	}
	//全局 响应拦截器
	responseFilters, responseFiltersOK := Plan.responseFilters[ConstNodeGlobal];
	responseFilterNames, responseFilterNamesOK := Plan.responseFilterNames[ConstNodeGlobal];
	if responseFiltersOK && responseFilterNamesOK {
		for _,n := range responseFilterNames {
			module.ResponseFilter(n, responseFilters[n])
		}
	}
	//节点 响应拦截器
	nodeResponseFilters, nodeResponseFiltersOK := Plan.responseFilters[node.Name];
	nodeResponseFilterNames, nodeResponseFilterNamesOK := Plan.responseFilterNames[node.Name];
	if nodeResponseFiltersOK && nodeResponseFilterNamesOK {
		for _,n := range nodeResponseFilterNames {
			module.ResponseFilter(n, nodeResponseFilters[n])
		}
	}







	//注意，处理器，优先注册节点的，在前面
	//全局的在后面， 这样可以做一个全局的默认处理器
	//然后在有些节点，可以优先自定义处理


	//节点 不存在处理器
	nodeFoundHandlers, nodeFoundHandlersOK := Plan.foundHandlers[node.Name];
	nodeFoundHandlerNames, nodeFoundHandlerNamesOK := Plan.foundHandlerNames[node.Name];
	if nodeFoundHandlersOK && nodeFoundHandlerNamesOK {
		for _,n := range nodeFoundHandlerNames {
			module.FoundHandler(n, nodeFoundHandlers[n])
		}
	}
	//全局 不存在处理器
	foundHandlers, foundHandlersOK := Plan.foundHandlers[ConstNodeGlobal];
	foundHandlerNames, foundHandlerNamesOK := Plan.foundHandlerNames[ConstNodeGlobal];
	if foundHandlersOK && foundHandlerNamesOK {
		for _,n := range foundHandlerNames {
			//不存在才注册，这样全局不会替换节点的处理器
			if _,ok := module.foundHandlers[n]; ok == false {
				module.FoundHandler(n, foundHandlers[n])
			}
		}
	}




	//节点 错误处理器
	nodeErrorHandlers, nodeErrorHandlersOK := Plan.errorHandlers[node.Name];
	nodeErrorHandlerNames, nodeErrorHandlerNamesOK := Plan.errorHandlerNames[node.Name];
	if nodeErrorHandlersOK && nodeErrorHandlerNamesOK {
		for _,n := range nodeErrorHandlerNames {
			module.ErrorHandler(n, nodeErrorHandlers[n])
		}
	}
	//全局 错误处理器
	errorHandlers, errorHandlersOK := Plan.errorHandlers[ConstNodeGlobal];
	errorHandlerNames, errorHandlerNamesOK := Plan.errorHandlerNames[ConstNodeGlobal];
	if errorHandlersOK && errorHandlerNamesOK {
		for _,n := range errorHandlerNames {
			//不存在才注册，这样全局不会替换节点的处理器
			if _,ok := module.errorHandlers[n]; ok == false {
				module.ErrorHandler(n, errorHandlers[n])
			}
		}
	}


	return module
}


















//创建Plan上下文
//func (module *planModule) newPlanContext(id string, name string, time string, value Map) (*PlanContext) {
func (module *planModule) newPlanContext(req *PlanRequest, res PlanResponse) (*PlanContext) {
	return &PlanContext{
		Node: module.node, Module: module,
		next: -1, nexts: []PlanFunc{},

		req: req, res: res,

		Id: req.Id, Name: req.Name, Config: nil, Branchs:nil,

		Value: req.Value, Local: Map{}, Item: Map{}, Auth: Map{}, Args: Map{},
	}
}



//计划Plan  请求开始
//func (module *planModule) servePlan(id string, name string, time string, value Map) {
func (module *planModule) servePlan(req *PlanRequest, res PlanResponse) {

	ctx := module.newPlanContext(req, res)

	ctx.handler(module.contextRequest)
	//最终所有的响应处理，优先
	ctx.handler(module.contextResponse)


	//中间件
	//用数组保证原始注册顺序
	for _,name := range Plan.middlerNames {
		ctx.handler(Plan.middlers[name])
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
	计划模块处理方法 begin
*/



//plan 计划处理器
//请求处理
//包含：route解析、request处理、session处理
func (module *planModule) contextRequest(ctx *PlanContext) {

	//计划不需要路由解析，直接new的时候就有name了
	if config,ok := module.routes[ctx.Name]; ok {
		ctx.Config = config
	} else {
		ctx.Config = nil
	}

	//请求处理
	//Id已经有了

	//会话处理
	m,err := module.sessionConnect.Query(ctx.Id)
	if err == nil {
		ctx.Session = m
	} else {
		ctx.Session = Map{}
	}
	ctx.Sign = &Sign{ ctx.Session }
	ctx.Next()
	//module.sessionConnect.Update(ctx.Id, ctx.Session, module.sessionConfig.Expiry)
	module.sessionConnect.Update(ctx.Id, ctx.Session)
}

//处理响应
func (module *planModule) contextResponse(ctx *PlanContext) {
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
func (module *planModule) contextExecute(ctx *PlanContext) {

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
		//if _,ok := ctx.Config[KeyMapAuth]; ok {
		//	ctx.handler(module.contextAuth)
		//}
		if _,ok := ctx.Config[KeyMapItem]; ok {
			ctx.handler(module.contextItem)
		}

		//最终都由分支处理
		ctx.handler(module.contextBranch)
	}

	ctx.Next()
}


//计划处理：处理分支
func (module *planModule) contextBranch(ctx *PlanContext) {

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
			case func(*PlanContext)bool:
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
	//计划路由不支持多method，非http
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
	} else {
		ctx.Config = nil
	}

	if ctx.Config == nil {
		//还是不存在的
		ctx.handler(module.contextFound)
	} else {




		//先处理参数，验证等的东西
		if _,ok := ctx.Config[KeyMapArgs]; ok {
			ctx.handler(module.contextArgs)
		}
		//if _,ok := ctx.Config[KeyMapAuth]; ok {
		//	ctx.handler(module.contextAuth)
		//}
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
			case func(*PlanContext):
				ctx.handler(actions)
			case []func(*PlanContext):
				for _,action := range actions {
					ctx.handler(action)
				}
			case PlanFunc:
				ctx.handler(actions)
			case []PlanFunc:
				ctx.handler(actions...)
			default:
			}
		}
	}

	ctx.Next()
}

















//自带中间件，参数处理
func (module *planModule) contextArgs(ctx *PlanContext) {

	//argn表示参数都可为空
	argn := false
	if v,ok := ctx.Config["argn"].(bool); ok {
		argn = v
	}

	//所有值都会放在 module.Value 中
	err := Mapping.Parse([]string{}, ctx.Config["args"].(Map), ctx.Value, ctx.Args, argn)
	if err != nil {
		ctx.Error(err)
	} else {
		ctx.Next()
	}
}



//Entity实体处理
func (module *planModule) contextItem(ctx *PlanContext) {
	if ctx.Config["item"] != nil {
		cfg := ctx.Config["item"].(Map)

		saveMap := Map{}

		for k,v := range cfg {
			config := v.(Map)

			name := config["name"].(string)
			key := k
			if config["value"] != nil && config["value"] != "" {
				key = config["value"].(string)
			}

			if ctx.Value[key] == nil {
				//参数不为空啊啊
				state := "item.empty"
				//是否有自定义状态
				if v,ok := config["empty"]; ok {
					state = v.(string)
				}
				err := Const.NewTypeStateError(k, state, name)
				//查询不到东西，也要失败， 接口访问失败
				ctx.Error(err)
				return
			} else {

				//判断是否需要查询数据
				dataName,dok := config["base"].(string); modelName,mok := config["model"].(string);
				if dok && mok {

					//要查询库
					db := Data.Base(dataName);
					item,err := db.Model(modelName).Entity(ctx.Value[key])
					db.Close()
					if err != nil {
						state := "item.error"
						//是否有自定义状态
						if v,ok := config["error"]; ok {
							state = v.(string)
						}
						err := Const.NewTypeStateError(k, state, name)

						ctx.Error(err)
						return;
					} else {
						saveMap[k] = item
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

































//路由执行，found
func (module *planModule) contextFound(ctx *PlanContext) {
	//清理执行线
	ctx.cleanup()

	//如果路由配置中有found，就自定义处理
	if v,ok := ctx.Config[KeyMapFound]; ok {
		switch c := v.(type) {
		case PlanFunc: {
			ctx.handler(c)
		}
		case []PlanFunc: {
			for _,v := range c {
				ctx.handler(v)
			}
		}
		case func(*PlanContext): {
			ctx.handler(c)
		}
		case []func(*PlanContext): {
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
func (module *planModule) contextError(ctx *PlanContext) {
	//清理执行线
	ctx.cleanup()

	//如果路由配置中有found，就自定义处理
	if v,ok := ctx.Config[KeyMapError]; ok {
		switch c := v.(type) {
		case PlanFunc: {
			ctx.handler(c)
		}
		case []PlanFunc: {
			for _,v := range c {
				ctx.handler(v)
			}
		}
		case func(*PlanContext): {
			ctx.handler(c)
		}
		case []func(*PlanContext): {
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




/*
	计划模块方法 end
*/




/* 默认响应器 begin */
//这个才是真的响应处理
func (module *planModule) contextResponder(ctx *PlanContext) {

	if ctx.Body == nil {
		//没有响应，应该走到found流程
		module.contextFound(ctx)
	}


	switch ctx.Body.(type) {
	case planBodyFinish:
		module.finishResponder(ctx)
	case planBodyReplan:
		module.replanResponder(ctx)
	default:
		module.defaultResponder(ctx)
	}

}





/* 默认响应器 begin */
func (module *planModule) finishResponder(ctx *PlanContext) {
	ctx.res.Finish(ctx.Id)
}
func (module *planModule) replanResponder(ctx *PlanContext) {
	body := ctx.Body.(planBodyReplan)
	ctx.res.Replan(ctx.Id, body.Delay)
}
func (module *planModule) defaultResponder(ctx *PlanContext) {
	ctx.res.Finish(ctx.Id)
}
/* 默认响应器 end */




/* 默认处理器 begin */
//代码中没有指定相关的处理器，才会执行到默认处理器
func (module *planModule) foundDefaultHandler(ctx *PlanContext) {
	ctx.res.Finish(ctx.Id)
}
func (module *planModule) errorDefaultHandler(ctx *PlanContext) {
	ctx.res.Finish(ctx.Id)
}
/* 默认处理器 end */









































/*
	计划上下文处理方法 begin
*/



//添加执行线
func (ctx *PlanContext) handler(handlers ...PlanFunc) {
	for _,handler := range handlers {
		ctx.nexts = append(ctx.nexts, handler)
	}
}
//清空执行线
func (ctx *PlanContext) cleanup() {
	ctx.next = -1
	ctx.nexts = make([]PlanFunc, 0)
}

/* 执行下一个 */
func (ctx *PlanContext) Next() {
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
func (ctx *PlanContext) Found() {
	ctx.Module.contextFound(ctx)
}
//返回错误
func (ctx *PlanContext) Error(err *Error) {
	ctx.Wrong = err
	ctx.Module.contextError(ctx)
}
/* 上下文处理器 end */







/* 上下文响应器 begin */
//完成操作
func (ctx *PlanContext) Finish() {
	ctx.Body = planBodyFinish{}
}
//重新触发
func (ctx *PlanContext) Replan(delays ...time.Duration) {
	if len(delays) > 0 {
		//延时重新触发
		ctx.Body = planBodyReplan{ Delay: delays[0] }
	} else {
		//立即重新触发
		ctx.Body = planBodyReplan{ Delay: time.Second*0 }
	}
}
/* 上下文响应器 end */











/*
	计划上下文方法 end
*/












