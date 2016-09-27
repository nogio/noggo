package noggo


import (
	"sync"
	"fmt"
	"time"
	. "github.com/nogio/noggo/base"
)


type (
	//计划函数
	PlanCall func(*PlanCtx)
	PlanMatch func(*PlanCtx) bool


	//计划驱动
	PlanDriver interface {
		Connect(config Map) (PlanConnect)
	}
	//计划连接
	PlanConnect interface {
		Open() error
		Close() error
		Accept(job PlanJob) error
		Remove(id string) error
		Start() error
		Stop() error
	}
	//计划工作实体
	PlanJob struct {
		Id		string
		Name	string
		Time	string
		Call	func()
	}
	//计划上下文
	PlanContext interface {
		//请求拦截器，在请求一开始执行
		RequestFilter(*PlanCtx)
		//响应拦截器，在响应开始前执行
		ResponseFilter(*PlanCtx)
		//执行拦截器，在执行action前执行
		ExecuteFilter(*PlanCtx)

		//404处理器，找不到请求时执行
		FoundHandler(*PlanCtx)
		//错误处理器，发生错误时执行
		ErrorHandler(*PlanCtx)
		//失败处理器，失败时执行，如参数解析失败
		FailedHandler(*PlanCtx)
		//拒绝处理器，拒绝时执行，主要用于Sign签名认证
		DeniedHandler(*PlanCtx)
	}





	//计划模块
	planModule struct {
		//驱动
		drivers map[string]PlanDriver
		driversMutex sync.Mutex
		//上下文
		contexts map[string]PlanContext
		contextsMutex sync.Mutex


		//路由器连接
		routerConfig	*routerConfig
		routerConnect	RouterConnect
		//会话连接
		sessionConfig	*sessionConfig
		sessionConnect	SessionConnect


		//计划连接
		planConfig		*planConfig
		planConnect		PlanConnect

		//路由
		routes 		map[string]Map			//路由定义
		routeNames	[]string				//路由名称原始顺序，因为map是无序的
		routeUris 	map[string]string		//记录所有uris指定	map[uri]name
		routeTimes	map[string][]string		//记录所有times定义，不同的计划可能会有相同的time定义,  map[name][times]
	}

	//计划上下文
	PlanCtx struct {
		//执行线
		nexts []PlanCall		//方法列表
		next int				//下一个索引

		//基础
		Id	string			//Session Id  会话时使用
		Session Map			//存储Session值
		Sign	*Sign		//签名功能，基于session

		//请求相关
		Method	string		//请求的method， 继承之web请求， 暂时无用
		Path	string		//请求的路径，演变自web， 暂时等于plan的名称


		//路由相关
		Name string			//解析路由后得到的name
		Config Map			//解析后得到的路由配置
		Branchs []Map		//解析后得到的路由分支配置

		//数据相关
		Params	Map			//路由解析后uri中的参数
		Value	Map			//所有请求过来的原始参数
		Locals	Map			//在ctx中传递数据用的
		Args	Map			//经过args处理后的参数
		Items	Map			//单条记录查询对象
		Auths	Map			//签名认证对象

		//响应相关
		//Code	int			//返回的状态码
		//Type	Type		//响应类型
		Body	Any			//响应内容
		Error	*Error		//响应错误
	}
)



/*
	计划模块方法 begin
*/




//计划初始化
func (plan *planModule) init() {
	plan.initRouter()
	plan.initSession()
	plan.initPlan()
}
//初始化路由驱动
func (plan *planModule) initRouter() {
	if Config.Plan.Router == nil {
		//使用默认的由路连接
		plan.routerConfig = Plan.routerConfig
		plan.routerConnect = Plan.routerConnect
	} else {
		//使用自定义的由路连接
		plan.routerConfig = Config.Plan.Router
		plan.routerConnect = Router.connect(plan.routerConfig)

		if plan.routerConnect == nil {
			panic("计划连接路由服务失败")
		} else {
			err := plan.routerConnect.Open()
			if err != nil {
				panic("计划打开路由服务失败 " + err.Error())
			}
		}
	}
}

//初始化会话驱动
func (plan *planModule) initSession() {
	if Config.Plan.Session == nil {
		//使用默认的会话连接
		plan.sessionConfig = Session.sessionConfig
		plan.sessionConnect = Session.sessionConnect
	} else {
		//使用自定义的会话连接
		plan.sessionConfig = Config.Plan.Session
		plan.sessionConnect = Session.connect(plan.sessionConfig)


		if plan.sessionConnect == nil {
			panic("计划连接会话服务失败")
		} else {
			//打开会话连接
			err := plan.sessionConnect.Open()
			if err != nil {
				panic("计划打开会话服务失败 " + err.Error())
			}
		}
	}
}


//初始化计划驱动
func (plan *planModule) initPlan() {

	plan.planConfig = Config.Plan
	plan.planConnect = plan.connect(plan.planConfig)

	//打开计划连接
	err := plan.planConnect.Open()
	if err != nil {
		panic("打开计划服务失败 " + err.Error())
	}

	//处理计划
	for name,times := range plan.routeTimes {
		for index,time := range times {
			//监听计划
			plan.planConnect.Accept(PlanJob{
				Id: fmt.Sprintf("%v.%v", name, index),
				Name: name, Time: time,
				Call: func() {
					plan.Touch(name)
				},
			})
		}
	}

	//开始计划
	plan.planConnect.Start()
}



//计划退出
func (plan *planModule) exit() {
	plan.exitRouter()
	plan.exitSession()
	plan.exitPlan()
}
//计划退出，路由器
func (plan *planModule) exitRouter() {
	//关闭路由
	if plan.routerConnect != nil {
		plan.routerConnect.Close()
		plan.routerConnect = nil
	}
}
//计划退出，会话
func (plan *planModule) exitSession() {
	//关闭会话
	if plan.sessionConnect != nil {
		plan.sessionConnect.Close()
		plan.sessionConnect = nil
	}
}
//计划退出，计划
func (plan *planModule) exitPlan() {
	//关闭计划
	if plan.planConnect != nil {
		plan.planConnect.Close()
		plan.planConnect = nil
	}
}


















//Run是大多cron库支持的
func (job PlanJob) Run() {
	if job.Call != nil {
		job.Call()
	}
}
//Launch是大多cron库支持的
func (job PlanJob) Launch() {
	if job.Call != nil {
		job.Call()
	}
}





//连接驱动
func (plan *planModule) connect(config *planConfig) (PlanConnect) {
	if planDriver,ok := plan.drivers[config.Driver]; ok {
		return planDriver.Connect(config.Config)
	} else {
		panic("不支持的计划驱动： " + config.Driver)
	}
}




//注册驱动
func (plan *planModule) Driver(name string, driver PlanDriver) {
	plan.driversMutex.Lock()
	defer plan.driversMutex.Unlock()

	if driver == nil {
		panic("plan: Register driver is nil")
	}
	if _, ok := plan.drivers[name]; ok {
		panic("plan: Registered driver " + name)
	}

	plan.drivers[name] = driver
}
//注册上下文
func (plan *planModule) Context(name string, context PlanContext) {
	plan.contextsMutex.Lock()
	defer plan.contextsMutex.Unlock()

	if context == nil {
		panic("plan: register context is nil")
	}
	if _, ok := plan.contexts[name]; ok {
		panic("plan: registered context " + name)
	}

	plan.contexts[name] = context
}
//注册路由
func (plan *planModule) Route(name string, config Map) {
	//保存配置
	plan.routes[name] = config
	plan.routeNames = append(plan.routeNames, name)

	//处理uri
	plan.routeUris[name] = name
	if v,ok := config[KeyMapUri]; ok {

		switch uris := v.(type) {
		case string:
			plan.routeUris[uris] = name
		case []string:
			for _,uri := range uris {
				plan.routeUris[uri] = name
			}
		}
	}

	//处理time
	if v,ok := config[KeyMapTime]; ok {
		switch times := v.(type) {
		case string:
			plan.routeTimes[name] = []string { times }
		case []string:
			plan.routeTimes[name] = times
		}
	}
}


















//创建Plan上下文
func (plan *planModule) newPlan(method, path string, value Map) (*PlanCtx) {
	return &PlanCtx{
		next: -1, nexts: []PlanCall{},

		Method: method, Path: path,
		Name: "", Config: Map{}, Branchs:[]Map{},

		Params: Map{}, Value: value, Locals: Map{},
		Args: Map{}, Items: Map{}, Auths: Map{},
	}
}



//计划Plan  请求开始
func (plan *planModule) servePlan(method, path string, value Map) {
	ctx := plan.newPlan(method, path, value)

	//请求处理
	for _,v := range plan.contexts {
		ctx.handler(v.RequestFilter)
	}
	ctx.handler(plan.handlerRequest)

	//响应处理
	ctx.handler(plan.handlerResponse)
	for _,v := range plan.contexts {
		ctx.handler(v.ResponseFilter)
	}

	//开始执行
	ctx.handler(plan.handlerExecute)
	ctx.Next()
}














//发起触发
func (plan *planModule) Touch(path string, args ...Map) {

	value := Map{}
	if len(args) > 0 {
		value = args[0]
	}

	go plan.servePlan("", path, value)
}

/*
	计划模块方法  end
*/




/*
	计划模块处理方法 begin
*/



//plan 计划处理器
//请求处理
//包含：route解析、request处理、session处理
func (plan *planModule) handlerRequest(ctx *PlanCtx) {

	//路由解析
	//目前暂不支持driver
	//直接使用name相等就匹配
	if name,ok := plan.routeUris[ctx.Path]; ok {
		ctx.Name = name
		ctx.Config = plan.routes[name]
	}


	//请求处理
	//主要是SessionId处理、处理传过来的值或表单
	ctx.Id = ctx.Name	//使用name做为id，以便在同一个计划之下共享session

	//会话处理
	ctx.Session = plan.sessionConnect.Create(ctx.Id, plan.sessionConfig.Expiry)
	ctx.Sign = &Sign{ ctx.Session }
	ctx.Next()
	plan.sessionConnect.Update(ctx.Id, ctx.Session, plan.sessionConfig.Expiry)
}

//处理响应
func (plan *planModule) handlerResponse(ctx *PlanCtx) {
	ctx.Next()

	if ctx.Body == nil {
		//没有响应，应该走到found流程
	} else {

		switch body := ctx.Body.(type) {
		case BodyPlanFinish:
			//完成不做任何处理
		case BodyPlanReplan:
			//目前直接调度，可调整，以后做到task中统一调整
			//因为万一delay很久。中间正好程序重新或是其它，就丢了
			//所以有必要使用task机制重新调度
			time.AfterFunc(body.Delay, func() {
				Plan.Touch(ctx.Path, ctx.Value)
			})
		default:
			//默认，也没有什么好处理的
		}
	}
}



//路由执行，处理
func (plan *planModule) handlerExecute(ctx *PlanCtx) {

	//解析路由，拿到actions
	if ctx.Config == nil {

		//找不到路由
		//ctx.handler(plan.handlerFound)
	} else {





		//先走，filter拦截器
		/*
		for _,n := range plan.Service.filterNames {
			plan.handler(plan.Service.FilterActions(n)...)
		}
		*/

		//验证，参数，数据处理
		//验证处理，数据处理， 可以考虑走中间件
		/*
		if _,ok := plan.Config["args"]; ok {
			plan.handler(planArgs)
		}
		if _,ok := plan.Config["auth"]; ok {
			plan.handler(planAuth)
		}
		if _,ok := plan.Config["item"]; ok {
			plan.handler(planItem)
		}
		*/




		//最终都由分支处理
		ctx.handler(plan.handlerBranch)
	}

	ctx.Next()
}


//计划处理：处理分支
func (plan *planModule) handlerBranch(ctx *PlanCtx) {

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
			case func(*PlanCtx)bool:
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
	//复制顶层的路由配置
	//noggo更新， 应该复制一下， 这样可以省一个handler，execute直接不要了。就直接分支
	//顶层的复制主要是 auth, item 的处理
	for k,v := range ctx.Route {
		if k != "uri" && k != "match" && k != "route" && k != "branch" {
			routing[k] = v
		}
	}
	*/


	//这里 ctx.Route 和 routing 变换位置
	ctx.Config = Map{}

	//如果有路由
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
			}
		}
	}




	//先处理参数，验证等的东西
	/*
	if _,ok := ctx.Config["args"]; ok {
		ctx.handler(planArgs)
	}
	if _,ok := ctx.Config["auth"]; ok {
		ctx.handler(planAuth)
	}
	if _,ok := ctx.Config["item"]; ok {
		ctx.handler(planItem)
	}
	*/


	//action之前Execute拦截器
	for _,v := range plan.contexts {
		ctx.handler(v.ExecuteFilter)
	}

	//把action加入调用列表
	if actionConfig,ok := ctx.Config[KeyMapAction]; ok {
		switch actions:=actionConfig.(type) {
		case func(*PlanCtx):
			ctx.handler(actions)
		case []func(*PlanCtx):
			for _,action := range actions {
				ctx.handler(action)
			}
		case PlanCall:
			ctx.handler(actions)
		case []PlanCall:
			ctx.handler(actions...)
		default:
		}
	}

	ctx.Next()
}




/*
	计划模块方法 end
*/






/*
	计划上下文处理方法 begin
*/



//添加执行线
func (plan *PlanCtx) handler(handlers ...PlanCall) {
	for _,handler := range handlers {
		plan.nexts = append(plan.nexts, handler)
	}
}
//清空执行线
func (plan *PlanCtx) cleanup() {
	plan.next = -1
	plan.nexts = make([]PlanCall, 0)
}

/* 执行下一个 */
func (plan *PlanCtx) Next() {
	plan.next++
	if len(plan.nexts) > plan.next {
		next := plan.nexts[plan.next]
		if next != nil {
			next(plan)
		} else {
			plan.Next()
		}
	} else {
		//没有了，不再执行，Response会处理为404
	}
}


//计划响应
//完成操作
func (ctx *PlanCtx) Finish() {
	ctx.Body = BodyPlanFinish{}
}

//计划响应
//重新触发
func (ctx *PlanCtx) Replan(delays ...time.Duration) {
	if len(delays) > 0 {
		//延时重新触发
		ctx.Body = BodyPlanReplan{ Delay: delays[0] }
	} else {
		//立即重新触发
		ctx.Body = BodyPlanReplan{ Delay: time.Second*0 }
	}
}

/*
	计划上下文方法 end
*/

