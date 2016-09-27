package noggo


import (
	"time"
	"sync"
	. "github.com/nogio/noggo/base"
)


type (

	//触发器函数
	TriggerCall func(*TriggerCtx)
	TriggerMatch func(*TriggerCtx) bool

	TriggerContext interface {
		//请求拦截器，在请求一开始执行
		RequestFilter(*TriggerCtx)
		//响应拦截器，在响应开始前执行
		ResponseFilter(*TriggerCtx)
		//执行拦截器，在执行action前执行
		ExecuteFilter(*TriggerCtx)

		//404处理器，找不到请求时执行
		FoundHandler(*TriggerCtx)
		//错误处理器，发生错误时执行
		ErrorHandler(*TriggerCtx)
		//失败处理器，失败时执行，如参数解析失败
		FailedHandler(*TriggerCtx)
		//拒绝处理器，拒绝时执行，主要用于Sign签名认证
		DeniedHandler(*TriggerCtx)
	}

	//触发器模块
	triggerModule struct {
		//上下文
		contexts map[string]TriggerContext
		contextsMutex sync.Mutex


		//路由连接
		routerConfig	*routerConfig
		routerConnect	RouterConnect
		//会话连接
		sessionConfig	*sessionConfig
		sessionConnect	SessionConnect

		//路由
		routes map[string]Map			//路由定义
		routeNames []string				//路由名称原始顺序，因为map是无序的
		routeUris map[string]string		//记录所有uris指定
	}

	//触发器上下文
	TriggerCtx struct {
		//执行线
		nexts []TriggerCall		//方法列表
		next int				//下一个索引

		//基础
		Id	string			//Session Id  会话时使用
		Session Map			//存储Session值
		Sign	*Sign		//签名功能，基于session

		//请求相关
		Method	string		//请求的method， 继承之web请求， 暂时无用
		Path	string		//请求的路径，演变自web， 暂时等于trigger的名称


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
	触发器模块方法 begin
*/




//注册上下文
func (trigger *triggerModule) Context(name string, context TriggerContext) {
	trigger.contextsMutex.Lock()
	defer trigger.contextsMutex.Unlock()

	if context == nil {
		panic("trigger: register context is nil")
	}
	if _, ok := trigger.contexts[name]; ok {
		panic("trigger: registered context " + name)
	}

	trigger.contexts[name] = context
}






















//触发器初始化
func (trigger *triggerModule) init() {
	trigger.initRouter()
	trigger.initSession()
}
//初始化路由驱动
func (trigger *triggerModule) initRouter() {
	if Config.Trigger.Router == nil {
		//使用默认的由路连接
		trigger.routerConfig = Router.routerConfig
		trigger.routerConnect = Router.routerConnect
	} else {
		//使用自定义的由路连接
		trigger.routerConfig = Config.Trigger.Router
		trigger.routerConnect = Router.connect(trigger.routerConfig)

		if trigger.routerConnect == nil {
			panic("触发器连接路由服务失败")
		} else {
			err := trigger.routerConnect.Open()
			if err != nil {
				panic("触发器打开会话服务失败 " + err.Error())
			}
		}
	}
}

//初始化会话驱动
func (trigger *triggerModule) initSession() {
	if Config.Trigger.Session == nil {
		//使用默认的会话连接
		trigger.sessionConfig = Session.sessionConfig
		trigger.sessionConnect = Session.sessionConnect
	} else {
		//使用自定义的会话连接
		trigger.sessionConfig = Config.Trigger.Session
		trigger.sessionConnect = Session.connect(trigger.sessionConfig)

		if trigger.sessionConnect == nil {
			panic("触发器连接会话服务失败")
		} else {
			err := trigger.sessionConnect.Open()
			if err != nil {
				panic("触发器打开会话服务失败 " + err.Error())
			}
		}
	}
}

//触发器退出
func (trigger *triggerModule) exit() {
	//关闭路由
	if trigger.routerConnect != nil {
		trigger.routerConnect.Close()
		trigger.routerConnect = nil
	}
	//关闭会话
	if trigger.sessionConnect != nil {
		trigger.sessionConnect.Close()
		trigger.sessionConnect = nil
	}
}


















//触发器，注册路由
func (trigger *triggerModule) Route(name string, config Map) {
	//保存配置
	trigger.routes[name] = config
	trigger.routeNames = append(trigger.routeNames, name)

	//处理uri
	trigger.routeUris[name] = name
	if v,ok := config[KeyMapUri]; ok {

		switch uris := v.(type) {
		case string:
			trigger.routeUris[uris] = name
		case []string:
			for _,uri := range uris {
				trigger.routeUris[uri] = name
			}
		}
	}
}














//创建Trigger上下文
func (trigger *triggerModule) newTrigger(method, path string, value Map) (*TriggerCtx) {
	return &TriggerCtx{
		next: -1, nexts: []TriggerCall{},

		Method: method, Path: path,
		Name: "", Config: Map{}, Branchs:[]Map{},

		Params: Map{}, Value: value, Locals: Map{},
		Args: Map{}, Items: Map{}, Auths: Map{},
	}
}



//触发器Trigger  请求开始
func (trigger *triggerModule) serveTrigger(method, path string, value Map) {
	ctx := trigger.newTrigger(method, path, value)

	//请求处理
	for _,v := range trigger.contexts {
		ctx.handler(v.RequestFilter)
	}
	ctx.handler(trigger.handlerRequest)

	//响应处理
	ctx.handler(trigger.handlerResponse)
	for _,v := range trigger.contexts {
		ctx.handler(v.ResponseFilter)
	}

	//开始执行
	ctx.handler(trigger.handlerExecute)
	ctx.Next()
}














//发起触发
func (trigger *triggerModule) Touch(path string, args ...Map) {

	value := Map{}
	if len(args) > 0 {
		value = args[0]
	}

	go trigger.serveTrigger("", path, value)
}

/*
	触发器模块方法  end
*/




/*
	触发器模块处理方法 begin
*/



//trigger 触发器处理器
//请求处理
//包含：route解析、request处理、session处理
func (trigger *triggerModule) handlerRequest(ctx *TriggerCtx) {

	//路由解析
	//目前暂不支持driver
	//直接使用name相等就匹配
	if name,ok := trigger.routeUris[ctx.Path]; ok {
		ctx.Name = name
		ctx.Config = trigger.routes[name]
	}


	//请求处理
	//主要是SessionId处理、处理传过来的值或表单
	ctx.Id = ctx.Name	//使用name做为id，以便在同一个触发器之下共享session

	//会话处理
	ctx.Session = trigger.sessionConnect.Create(ctx.Id, trigger.sessionConfig.Expiry)
	ctx.Sign = &Sign{ ctx.Session }
	ctx.Next()
	trigger.sessionConnect.Update(ctx.Id, ctx.Session, trigger.sessionConfig.Expiry)
}

//处理响应
func (trigger *triggerModule) handlerResponse(ctx *TriggerCtx) {
	ctx.Next()

	if ctx.Body == nil {
		//没有响应，应该走到found流程
	} else {

		switch body := ctx.Body.(type) {
		case BodyTriggerFinish:
			//完成不做任何处理
		case BodyTriggerRetrigger:
			//目前直接调度，可调整，以后做到task中统一调整
			//因为万一delay很久。中间正好程序重新或是其它，就丢了
			//所以有必要使用task机制重新调度
			time.AfterFunc(body.Delay, func() {
				Trigger.Touch(ctx.Path, ctx.Value)
			})
		default:
			//默认，也没有什么好处理的
		}
	}
}



//路由执行，处理
func (trigger *triggerModule) handlerExecute(ctx *TriggerCtx) {

	//解析路由，拿到actions
	if ctx.Config == nil {

		//找不到路由
		//ctx.handler(trigger.handlerFound)
	} else {





		//先走，filter拦截器
		/*
		for _,n := range trigger.Service.filterNames {
			trigger.handler(trigger.Service.FilterActions(n)...)
		}
		*/

		//验证，参数，数据处理
		//验证处理，数据处理， 可以考虑走中间件
		/*
		if _,ok := trigger.Config["args"]; ok {
			trigger.handler(triggerArgs)
		}
		if _,ok := trigger.Config["auth"]; ok {
			trigger.handler(triggerAuth)
		}
		if _,ok := trigger.Config["item"]; ok {
			trigger.handler(triggerItem)
		}
		*/




		//最终都由分支处理
		ctx.handler(trigger.handlerBranch)
	}

	ctx.Next()
}


//触发器处理：处理分支
func (trigger *triggerModule) handlerBranch(ctx *TriggerCtx) {

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
			case func(*TriggerCtx)bool:
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
		ctx.handler(triggerArgs)
	}
	if _,ok := ctx.Config["auth"]; ok {
		ctx.handler(triggerAuth)
	}
	if _,ok := ctx.Config["item"]; ok {
		ctx.handler(triggerItem)
	}
	*/


	//action之前Execute拦截器
	for _,v := range trigger.contexts {
		ctx.handler(v.ExecuteFilter)
	}

	//把action加入调用列表
	if actionConfig,ok := ctx.Config[KeyMapAction]; ok {
		switch actions:=actionConfig.(type) {
		case func(*TriggerCtx):
			ctx.handler(actions)
		case []func(*TriggerCtx):
			for _,action := range actions {
				ctx.handler(action)
			}
		case TriggerCall:
			ctx.handler(actions)
		case []TriggerCall:
			ctx.handler(actions...)
		default:
		}
	}

	ctx.Next()
}




/*
	触发器模块方法 end
*/






/*
	触发器上下文处理方法 begin
*/



//添加执行线
func (trigger *TriggerCtx) handler(handlers ...TriggerCall) {
	for _,handler := range handlers {
		trigger.nexts = append(trigger.nexts, handler)
	}
}
//清空执行线
func (trigger *TriggerCtx) cleanup() {
	trigger.next = -1
	trigger.nexts = make([]TriggerCall, 0)
}

/* 执行下一个 */
func (trigger *TriggerCtx) Next() {
	trigger.next++
	if len(trigger.nexts) > trigger.next {
		next := trigger.nexts[trigger.next]
		if next != nil {
			next(trigger)
		} else {
			trigger.Next()
		}
	} else {
		//没有了，不再执行
	}
}


//触发器响应
//完成操作
func (ctx *TriggerCtx) Finish() {
	ctx.Body = BodyTriggerFinish{}
}

//触发器响应
//重新触发
func (ctx *TriggerCtx) Retrigger(delays ...time.Duration) {
	if len(delays) > 0 {
		//延时重新触发
		ctx.Body = BodyTriggerRetrigger{ Delay: delays[0] }
	} else {
		//立即重新触发
		ctx.Body = BodyTriggerRetrigger{ Delay: time.Second*0 }
	}
}

/*
	触发器上下文方法 end
*/

