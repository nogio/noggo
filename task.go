/*
	task 任务模块
	任务模块，是一个全局模块
	用于进程内的一些触发，比如数据被创建，修改，删除等等
	并且任务不需要三方驱动
*/

package noggo


import (
	"time"
	"sync"
	. "github.com/nogio/noggo/base"
	"errors"
)



// task driver begin


type (
	//任务驱动
	TaskDriver interface {
		Connect(config Map) (TaskConnect,error)
	}

	//回调函数
	TaskHandler func(*TaskRequest, TaskResponse)

	//任务连接
	TaskConnect interface {
		//打开连接
		Open() error
		//关闭连接
		Close() error

		//注册任务
		Register(string) error

		//注册回调
		Accept(TaskHandler) error

		//开始任务
		Start() error

		//触发任务
		After(name string, delay time.Duration, value Map) error
	}

	//任务请求实体
	TaskRequest struct {
		Id string
		Name string
		Delay time.Duration
		Value Map
	}

	//任务响应接口
	TaskResponse interface {
		//完成任务
		Finish(id string) error
		//重新开始任务
		Retask(id string, delay time.Duration) error
	}
)

// task driver end


type (
	//任务函数
	TaskFunc func(*TaskContext)

	//响应完成
	taskBodyFinish struct {
	}
	//响应重新触发
	taskBodyRetask struct {
		Delay time.Duration
	}

	//任务模块
	taskGlobal struct {
		mutex sync.Mutex

		//任务驱动容器
		drivers map[string]TaskDriver

		//中间件
		middlers    map[string]TaskFunc
		middlerNames []string

		//路由
		routes 		map[string]Map			//路由定义
		routeNames	[]string				//路由名称原始顺序，因为map是无序的

		//拦截器们
		requestFilters, executeFilters, responseFilters map[string]TaskFunc
		requestFilterNames, executeFilterNames, responseFilterNames []string

		//处理器们
		foundHandlers, errorHandlers map[string]TaskFunc
		foundHandlerNames, errorHandlerNames []string




		//会话配置与连接
		sessionConfig   *sessionConfig
		sessionConnect	SessionConnect
		//日志配置，日志连接
		taskConfig  *taskConfig
		taskConnect TaskConnect
	}

	//任务上下文
	TaskContext struct {
		Global	*taskGlobal

		//执行线
		nexts []TaskFunc		//方法列表
		next int				//下一个索引

		req *TaskRequest
		res TaskResponse

		//基础
		Id	string			//Session Id  会话时使用
		Session Map			//存储Session值
		Sign	*Sign		//签名功能，基于session

		//配置相关
		Name string			//解析路由后得到的name
		Config Map			//解析后得到的路由配置
		Branchs []Map		//解析后得到的路由分支配置

		//数据相关
		Delay	time.Duration	//延时
		Value	Map			//所有请求过来的原始参数汇总
		Local	Map			//在ctx中传递数据用的
		Item	Map			//单条记录查询对象
		Args	Map			//经过args处理后的参数

		//响应相关
		Body	Any			//响应内容

		Wrong	*Error		//错误信息
	}
)



/*
	任务模块方法 begin
*/


//连接驱动
func (global *taskGlobal) connect(config *taskConfig) (TaskConnect,error) {
	if taskDriver,ok := global.drivers[config.Driver]; ok {
		return taskDriver.Connect(config.Config)
	} else {
		panic("任务：不支持的驱动 " + config.Driver)
	}
}

//注册任务驱动
func (global *taskGlobal) Driver(name string, config TaskDriver) {
	global.mutex.Lock()
	defer global.mutex.Unlock()

	if config == nil {
		panic("任务: 驱动不可为空")
	}
	//不做存在判断，因为要支持后注册的驱动替换已注册的驱动
	//框架有可能自带几种默认驱动，并且是默认注册的，用户可以自行注册替换
	global.drivers[name] = config
}
func (global *taskGlobal) Middler(name string, call TaskFunc) {
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






//任务初始化
func (global *taskGlobal) init() {
	global.initSession()
	global.initTask()
}

//初始化会话驱动
func (global *taskGlobal) initSession() {
	//现在直接使用全局的会话连接
	global.sessionConfig = Session.sessionConfig
	global.sessionConnect = Session.sessionConnect
}

//初始化驱动
func (global *taskGlobal) initTask() {

	//先拿到默认的配置
	con,err := global.connect(Config.Task)

	if err != nil {
		panic("任务：连接失败：" + err.Error())
	} else {

		err := con.Open()
		if err != nil {
			panic("任务：打开失败 " + err.Error())
		}



		//注册回调

		//注册任务
		//貌似不需要注册了，因为只要注册一个NAME。 貌似没意义
		for _,name := range global.routeNames {
			con.Register(name)
		}

		//开始任务
		con.Accept(global.serveTask);


		//保存连接
		global.taskConnect = con
	}
}




//任务退出
func (global *taskGlobal) exit() {
	global.exitSession()
	global.exitTask()
}
//任务退出，会话
func (global *taskGlobal) exitSession() {
	//使用全局会话，不需要在这里关闭了
}
//任务退出
func (global *taskGlobal) exitTask() {
	//关闭连接
	if global.taskConnect != nil {
		global.taskConnect.Close()
	}
}













//任务：注册路由
func (global *taskGlobal) Route(name string, config Map) {
	global.mutex.Lock()
	defer global.mutex.Unlock()


	if global.routes == nil {
		global.routes = map[string]Map{}
	}
	if global.routeNames == nil {
		global.routeNames = []string{}
	}

	//保存配置
	if _,ok := global.routes[name]; ok == false {
		//没有注册过name，才把name加到列表
		global.routeNames = append(global.routeNames, name)
	}
	//可以后注册重写原有路由配置，所以直接保存
	global.routes[name] = config
}








/* 注册拦截器 begin */
func (global *taskGlobal) RequestFilter(name string, call TaskFunc) {
	global.mutex.Lock()
	defer global.mutex.Unlock()

	if global.requestFilters == nil {
		global.requestFilters = make(map[string]TaskFunc)
	}
	if global.requestFilterNames == nil {
		global.requestFilterNames = make([]string, 0)
	}

	//如果没有注册个此name，才加入数组
	if _,ok := global.requestFilters[name]; ok == false {
		global.requestFilterNames = append(global.requestFilterNames, name)
	}
	//函数直接写， 因为可以使用同名替换现有的
	global.requestFilters[name] = call
}
func (global *taskGlobal) ExecuteFilter(name string, call TaskFunc) {
	global.mutex.Lock()
	defer global.mutex.Unlock()

	if global.executeFilters == nil {
		global.executeFilters = make(map[string]TaskFunc)
	}
	if global.executeFilterNames == nil {
		global.executeFilterNames = make([]string, 0)
	}

	//如果没有注册个此name，才加入数组
	if _,ok := global.executeFilters[name]; ok == false {
		global.executeFilterNames = append(global.executeFilterNames, name)
	}
	//函数直接写， 因为可以使用同名替换现有的
	global.executeFilters[name] = call
}
func (global *taskGlobal) ResponseFilter(name string, call TaskFunc) {
	global.mutex.Lock()
	defer global.mutex.Unlock()

	if global.responseFilters == nil {
		global.responseFilters = make(map[string]TaskFunc)
	}
	if global.responseFilterNames == nil {
		global.responseFilterNames = make([]string, 0)
	}

	//如果没有注册个此name，才加入数组
	if _,ok := global.responseFilters[name]; ok == false {
		global.responseFilterNames = append(global.responseFilterNames, name)
	}
	//函数直接写， 因为可以使用同名替换现有的
	global.responseFilters[name] = call
}
/* 注册拦截器 end */


/* 注册处理器 begin */
func (global *taskGlobal) FoundHandler(name string, call TaskFunc) {
	global.mutex.Lock()
	defer global.mutex.Unlock()

	if global.foundHandlers == nil {
		global.foundHandlers = make(map[string]TaskFunc)
	}
	if global.foundHandlerNames == nil {
		global.foundHandlerNames = make([]string, 0)
	}

	//如果没有注册个此name，才加入数组
	if _,ok := global.foundHandlers[name]; ok == false {
		global.foundHandlerNames = append(global.foundHandlerNames, name)
	}
	//函数直接写， 因为可以使用同名替换现有的
	global.foundHandlers[name] = call
}
func (global *taskGlobal) ErrorHandler(name string, call TaskFunc) {
	global.mutex.Lock()
	defer global.mutex.Unlock()

	if global.errorHandlers == nil {
		global.errorHandlers = make(map[string]TaskFunc)
	}
	if global.errorHandlerNames == nil {
		global.errorHandlerNames = make([]string, 0)
	}

	//如果没有注册个此name，才加入数组
	if _,ok := global.errorHandlers[name]; ok == false {
		global.errorHandlerNames = append(global.errorHandlerNames, name)
	}
	//函数直接写， 因为可以使用同名替换现有的
	global.errorHandlers[name] = call
}
/* 注册处理器 end */










//创建Task上下文
//func (global *taskGlobal) newTaskContext(id string, name string, delay time.Duration, value Map) (*TaskContext) {
func (global *taskGlobal) newTaskContext(req *TaskRequest, res TaskResponse) (*TaskContext) {
	return &TaskContext{
		Global: global,
		next: -1, nexts: []TaskFunc{},

		req: req, res: res,

		Id: req.Id, Name: req.Name, Config: nil, Branchs:nil,
		Delay: req.Delay, Value: Map{}, Local: Map{}, Item: Map{}, Args: Map{},
	}
}



//任务Task  请求开始
//func (global *taskGlobal) serveTask(id string, name string, delay time.Duration, value Map) {
func (global *taskGlobal) serveTask(req *TaskRequest, res TaskResponse) {
	ctx := global.newTaskContext(req, res)

	//请求处理
	ctx.handler(global.contextRequest)
	//响应处理
	ctx.handler(global.contextResponse)

	//中间件
	//用数组保证原始注册顺序
	for _,name := range Task.middlerNames {
		ctx.handler(Task.middlers[name])
	}


	//filter中的request
	//用数组保证原始注册顺序
	for _,name := range global.requestFilterNames {
		ctx.handler(global.requestFilters[name])
	}

	//开始执行
	ctx.handler(global.contextExecute)
	ctx.Next()
}














//任务：触发
func (global *taskGlobal) After(name string, delay time.Duration, args ...Map) (error) {

	if global.taskConnect == nil {
		return errors.New("任务未连接")
	} else {
		value := Map{}
		if len(args) > 0 {
			value = args[0]
		}

		return global.taskConnect.After(name, delay, value)
	}
}

/*
	任务模块方法  end
*/



















/*
	任务模块处理方法 begin
*/



//task 任务处理器
//请求处理
//包含：route解析、request处理、session处理
func (global *taskGlobal) contextRequest(ctx *TaskContext) {

	//任务不需要路由解析，直接new的时候就有name了
	if config,ok := global.routes[ctx.Name]; ok {
		ctx.Config = config
	} else {
		ctx.Config = nil
	}

	//请求处理
	//Id已经有了


	//会话处理
	m,e := global.sessionConnect.Query(ctx.Id)
	if e == nil {
		ctx.Session = m
	} else {
		ctx.Session = Map{}
	}
	ctx.Sign = &Sign{ ctx.Session }
	ctx.Next()
	//global.sessionConnect.Update(ctx.Id, ctx.Session, global.sessionConfig.Expiry)
	global.sessionConnect.Update(ctx.Id, ctx.Session)
}

//处理响应
func (global *taskGlobal) contextResponse(ctx *TaskContext) {
	//因为response是在所有请求前的， 所以先调用一下
	//然后对结果进行处理
	ctx.Next()


	//清理执行线
	ctx.cleanup()

	//filter中的request
	//用数组保证原始注册顺序
	for _,name := range global.responseFilterNames {
		ctx.handler(global.responseFilters[name])
	}

	//这个函数才是真正响应的处理函数
	ctx.handler(global.contextResponder)

	ctx.Next()
}


//路由执行，处理
func (global *taskGlobal) contextExecute(ctx *TaskContext) {

	//解析路由，拿到actions
	if ctx.Config == nil {
		//找不到路由
		ctx.handler(global.contextFound)
	} else {


		//验证，参数，数据处理
		//验证处理，数据处理， 可以考虑走外部中间件
		if _,ok := ctx.Config[KeyMapArgs]; ok {
			ctx.handler(global.contextArgs)
		}
		if _,ok := ctx.Config[KeyMapItem]; ok {
			ctx.handler(global.contextItem)
		}

		//最终都由分支处理
		ctx.handler(global.contextBranch)
	}

	ctx.Next()
}


//任务处理：处理分支
func (global *taskGlobal) contextBranch(ctx *TaskContext) {

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
			case func(*TaskContext)bool:
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
	//任务路由不支持多method，非http
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
		ctx.handler(global.contextFound)
	} else {




		//先处理参数，验证等的东西
		if _,ok := ctx.Config[KeyMapArgs]; ok {
			ctx.handler(global.contextArgs)
		}
		if _,ok := ctx.Config[KeyMapItem]; ok {
			ctx.handler(global.contextItem)
		}


		//action之前的拦截器
		//filter中的execute
		//用数组保证原始注册顺序
		for _,name := range global.executeFilterNames {
			ctx.handler(global.executeFilters[name])
		}

		//把action加入调用列表
		if actionConfig,ok := ctx.Config[KeyMapAction]; ok {
			switch actions:=actionConfig.(type) {
			case func(*TaskContext):
				ctx.handler(actions)
			case []func(*TaskContext):
				for _,action := range actions {
					ctx.handler(action)
				}
			case TaskFunc:
				ctx.handler(actions)
			case []TaskFunc:
				ctx.handler(actions...)
			default:
			}
		}
	}

	ctx.Next()
}

















//自带中间件，参数处理
func (global *taskGlobal) contextArgs(ctx *TaskContext) {

	//argn表示参数都可为空
	argn := false
	if v,ok := ctx.Config["argn"].(bool); ok {
		argn = v
	}

	//所有值都会放在 global.Value 中
	err := Mapping.Parse([]string{}, ctx.Config["args"].(Map), ctx.Value, ctx.Args, argn, false)
	if err != nil {
		ctx.Error(err)
	} else {
		ctx.Next()
	}
}



//Entity实体处理
func (global *taskGlobal) contextItem(ctx *TaskContext) {
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
func (global *taskGlobal) contextFound(ctx *TaskContext) {
	//清理执行线
	ctx.cleanup()

	//如果路由配置中有found，就自定义处理
	if v,ok := ctx.Config[KeyMapFound]; ok {
		switch c := v.(type) {
		case TaskFunc: {
			ctx.handler(c)
		}
		case []TaskFunc: {
			for _,v := range c {
				ctx.handler(v)
			}
		}
		case func(*TaskContext): {
			ctx.handler(c)
		}
		case []func(*TaskContext): {
			for _,v := range c {
				ctx.handler(v)
			}
		}
		default:
		}
	}

	//handler中的found
	//用数组保证原始注册顺序
	for _,name := range global.foundHandlerNames {
		ctx.handler(global.foundHandlers[name])
	}

	//最后是默认found中间件
	ctx.handler(global.foundDefaultHandler)

	ctx.Next()
}



//路由执行，error
func (global *taskGlobal) contextError(ctx *TaskContext) {
	//清理执行线
	ctx.cleanup()

	//如果路由配置中有found，就自定义处理
	if v,ok := ctx.Config[KeyMapError]; ok {
		switch c := v.(type) {
		case TaskFunc: {
			ctx.handler(c)
		}
		case []TaskFunc: {
			for _,v := range c {
				ctx.handler(v)
			}
		}
		case func(*TaskContext): {
			ctx.handler(c)
		}
		case []func(*TaskContext): {
			for _,v := range c {
				ctx.handler(v)
			}
		}
		default:
		}
	}


	//handler中的error
	//用数组保证原始注册顺序
	for _,name := range global.errorHandlerNames {
		ctx.handler(global.errorHandlers[name])
	}

	//最后是默认error中间件
	ctx.handler(global.errorDefaultHandler)

	ctx.Next()
}


/*
	任务模块方法 end
*/




//处理响应
func (global *taskGlobal) contextResponder(ctx *TaskContext) {

	if ctx.Body == nil {
		//没有响应，应该走到found流程
		global.contextFound(ctx)
	}


	switch ctx.Body.(type) {
	case taskBodyFinish:
		global.finishResponder(ctx)
	case taskBodyRetask:
		global.retaskResponder(ctx)
	default:
		global.defaultResponder(ctx)
	}
}
















/* 默认响应器 begin */
func (global *taskGlobal) finishResponder(ctx *TaskContext) {
	//通知驱动，任务完成
	ctx.res.Finish(ctx.Id)
}

//目前直接调度，可调整，以后做到task中统一调整
//因为万一delay很久。中间正好程序重新或是其它，就丢了
//所以有必要使用task机制重新调度
func (global *taskGlobal) retaskResponder(ctx *TaskContext) {
	body := ctx.Body.(taskBodyRetask)
	//重新处理任务
	ctx.res.Retask(ctx.Id, body.Delay)
}
func (global *taskGlobal) defaultResponder(ctx *TaskContext) {
	//默认处理器， 一般执行不到。 默认完成吧
	ctx.res.Finish(ctx.Id)
}
/* 默认响应器 end */




/* 默认处理器 begin */
//代码中没有指定相关的处理器，才会执行到默认处理器
func (global *taskGlobal) foundDefaultHandler(ctx *TaskContext) {
	//当找不到任务时，应当通知驱动，完成此任务，以免重复调用
	ctx.res.Finish(ctx.Id)
}
func (global *taskGlobal) errorDefaultHandler(ctx *TaskContext) {
	//出错，此任务就完成了
	ctx.res.Finish(ctx.Id)
}
/* 默认处理器 end */









































/*
	任务上下文处理方法 begin
*/



//添加执行线
func (ctx *TaskContext) handler(handlers ...TaskFunc) {
	for _,handler := range handlers {
		ctx.nexts = append(ctx.nexts, handler)
	}
}
//清空执行线
func (ctx *TaskContext) cleanup() {
	ctx.next = -1
	ctx.nexts = make([]TaskFunc, 0)
}

/* 执行下一个 */
func (ctx *TaskContext) Next() {
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
func (ctx *TaskContext) Found() {
	ctx.Global.contextFound(ctx)
}
//失败, 就是参数处理失败为主
func (ctx *TaskContext) Error(err *Error) {
	ctx.Wrong = err
	ctx.Global.contextError(ctx)
}
/* 上下文处理器 end */







/* 上下文响应器 begin */
//完成操作
func (ctx *TaskContext) Finish() {
	ctx.Body = taskBodyFinish{}
}
//重新触发
func (ctx *TaskContext) Retask(delays ...time.Duration) {
	if len(delays) > 0 {
		//延时重新触发
		ctx.Body = taskBodyRetask{ Delay: delays[0] }
	} else {
		//立即重新触发
		ctx.Body = taskBodyRetask{ Delay: time.Second*0 }
	}
}
/* 上下文响应器 end */











/*
	任务上下文方法 end
*/

