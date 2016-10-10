/*
	trigger 触发器模块
	触发器模块，是一个全局模块
	用于进程内的一些触发，比如数据被创建，修改，删除等等
	并且触发器不需要三方驱动
*/

package noggo


import (
	"time"
	"sync"
	. "github.com/nogio/noggo/base"
)


type (
	//触发器函数
	TriggerFunc func(*TriggerContext)

	//响应完成
	triggerBodyFinish struct {
	}
	//响应重新触发
	triggerBodyRetrigger struct {
		Delay time.Duration
	}

	//触发器模块
	triggerGlobal struct {
		mutex sync.Mutex

		//会话配置与连接
		sessionConfig	*sessionConfig
		sessionConnect	SessionConnect


		//路由
		routes 		map[string]Map			//路由定义
		routeNames	[]string				//路由名称原始顺序，因为map是无序的

		//拦截器们
		requestFilters, executeFilters, responseFilters map[string]TriggerFunc
		requestFilterNames, executeFilterNames, responseFilterNames []string

		//处理器们
		foundHandlers, errorHandlers, failedHandlers, deniedHandlers map[string]TriggerFunc
		foundHandlerNames, errorHandlerNames, failedHandlerNames, deniedHandlerNames []string
	}

	//触发器上下文
	TriggerContext struct {
		Global	*triggerGlobal

		//执行线
		nexts []TriggerFunc		//方法列表
		next int				//下一个索引

		//基础
		Id	string			//Session Id  会话时使用
		Session Map			//存储Session值
		Sign	*Sign		//签名功能，基于session

		//配置相关
		Name string			//解析路由后得到的name
		Config Map			//解析后得到的路由配置
		Branchs []Map		//解析后得到的路由分支配置

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



/*
	触发器模块方法 begin
*/




//触发器初始化
func (global *triggerGlobal) init() {
	global.initSession()
}

//初始化会话驱动
func (global *triggerGlobal) initSession() {
	if Config.Trigger.Session != nil {
		//使用自定的
		global.sessionConfig = Config.Trigger.Session
	} else {
		//如果触发器中会话配置为空，使用默认的会话配置
		global.sessionConfig = Config.Session
	}

	//连接会话
	global.sessionConnect = Session.connect(global.sessionConfig)

	if global.sessionConnect == nil {
		panic("触发器：连接会话失败")
	} else {
		//打开会话连接
		err := global.sessionConnect.Open()
		if err != nil {
			panic("触发器：打开会话失败 " + err.Error())
		}
	}
}




//触发器退出
func (global *triggerGlobal) exit() {
	global.exitSession()
}
//触发器退出，会话
func (global *triggerGlobal) exitSession() {
	//关闭会话
	if global.sessionConnect != nil {
		global.sessionConnect.Close()
		global.sessionConnect = nil
	}
}













//触发器：注册路由
func (global *triggerGlobal) Route(name string, config Map) {
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
func (global *triggerGlobal) RequestFilter(name string, call TriggerFunc) {
	global.mutex.Lock()
	defer global.mutex.Unlock()

	if global.requestFilters == nil {
		global.requestFilters = make(map[string]TriggerFunc)
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
func (global *triggerGlobal) ExecuteFilter(name string, call TriggerFunc) {
	global.mutex.Lock()
	defer global.mutex.Unlock()

	if global.executeFilters == nil {
		global.executeFilters = make(map[string]TriggerFunc)
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
func (global *triggerGlobal) ResponseFilter(name string, call TriggerFunc) {
	global.mutex.Lock()
	defer global.mutex.Unlock()

	if global.responseFilters == nil {
		global.responseFilters = make(map[string]TriggerFunc)
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
func (global *triggerGlobal) FoundHandler(name string, call TriggerFunc) {
	global.mutex.Lock()
	defer global.mutex.Unlock()

	if global.foundHandlers == nil {
		global.foundHandlers = make(map[string]TriggerFunc)
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
func (global *triggerGlobal) ErrorHandler(name string, call TriggerFunc) {
	global.mutex.Lock()
	defer global.mutex.Unlock()

	if global.errorHandlers == nil {
		global.errorHandlers = make(map[string]TriggerFunc)
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
func (global *triggerGlobal) FailedHandler(name string, call TriggerFunc) {
	global.mutex.Lock()
	defer global.mutex.Unlock()

	if global.failedHandlers == nil {
		global.failedHandlers = make(map[string]TriggerFunc)
	}
	if global.failedHandlerNames == nil {
		global.failedHandlerNames = make([]string, 0)
	}

	//如果没有注册个此name，才加入数组
	if _,ok := global.failedHandlers[name]; ok == false {
		global.failedHandlerNames = append(global.failedHandlerNames, name)
	}
	//函数直接写， 因为可以使用同名替换现有的
	global.failedHandlers[name] = call
}
func (global *triggerGlobal) DeniedHandler(name string, call TriggerFunc) {
	global.mutex.Lock()
	defer global.mutex.Unlock()

	if global.deniedHandlers == nil {
		global.deniedHandlers = make(map[string]TriggerFunc)
	}
	if global.deniedHandlerNames == nil {
		global.deniedHandlerNames = make([]string, 0)
	}

	//如果没有注册个此name，才加入数组
	if _,ok := global.deniedHandlers[name]; ok == false {
		global.deniedHandlerNames = append(global.deniedHandlerNames, name)
	}
	//函数直接写， 因为可以使用同名替换现有的
	global.deniedHandlers[name] = call
}
/* 注册处理器 end */










//创建Trigger上下文
func (global *triggerGlobal) newTriggerContext(name string, value Map) (*TriggerContext) {
	return &TriggerContext{
		Global: global,
		next: -1, nexts: []TriggerFunc{},

		Name: name, Config: Map{}, Branchs:[]Map{},

		Value: value, Local: Map{}, Item: Map{}, Auth: Map{}, Args: Map{},
	}
}



//触发器Trigger  请求开始
func (global *triggerGlobal) serveTrigger(name string, value Map) {
	ctx := global.newTriggerContext(name, value)

	//请求处理
	ctx.handler(global.contextRequest)
	//filter中的request
	//用数组保证原始注册顺序
	for _,name := range global.requestFilterNames {
		ctx.handler(global.requestFilters[name])
	}

	//响应处理
	ctx.handler(global.contextResponse)
	//filter中的response
	//用数组保证原始注册顺序
	for _,name := range global.responseFilterNames {
		ctx.handler(global.responseFilters[name])
	}

	//开始执行
	ctx.handler(global.contextExecute)
	ctx.Next()
}














//触发器：触发
func (global *triggerGlobal) Touch(name string, args ...Map) {

	value := Map{}
	if len(args) > 0 {
		value = args[0]
	}

	go global.serveTrigger(name, value)
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
func (global *triggerGlobal) contextRequest(ctx *TriggerContext) {

	//触发器不需要路由解析，直接new的时候就有name了
	if config,ok := global.routes[ctx.Name]; ok {
		ctx.Config = config
	} else {
		ctx.Config = nil
	}

	//请求处理
	//主要是SessionId处理、处理传过来的值或表单
	ctx.Id = ctx.Name	//使用name做为id，以便在同一个触发器之下共享session

	//会话处理
	err,m := global.sessionConnect.Query(ctx.Id, global.sessionConfig.Expiry)
	if err == nil {
		ctx.Session = m
	} else {
		ctx.Session = Map{}
	}
	ctx.Sign = &Sign{ ctx.Session }
	ctx.Next()
	global.sessionConnect.Update(ctx.Id, ctx.Session, global.sessionConfig.Expiry)
}

//处理响应
func (global *triggerGlobal) contextResponse(ctx *TriggerContext) {
	ctx.Next()


	if ctx.Body == nil {
		//没有响应，应该走到found流程
		global.contextFound(ctx)
	}


	switch ctx.Body.(type) {
	case triggerBodyFinish:
		global.finishResponder(ctx)
	case triggerBodyRetrigger:
		global.retriggerResponder(ctx)
	default:
		global.defaultResponder(ctx)
	}
}



//路由执行，处理
func (global *triggerGlobal) contextExecute(ctx *TriggerContext) {

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
		if _,ok := ctx.Config[KeyMapAuth]; ok {
			ctx.handler(global.contextAuth)
		}
		if _,ok := ctx.Config[KeyMapItem]; ok {
			ctx.handler(global.contextItem)
		}

		//最终都由分支处理
		ctx.handler(global.contextBranch)
	}

	ctx.Next()
}


//触发器处理：处理分支
func (global *triggerGlobal) contextBranch(ctx *TriggerContext) {

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
			case func(*TriggerContext)bool:
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
	//触发器路由不支持多method，非http
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
		ctx.handler(global.contextArgs)
	}
	if _,ok := ctx.Config[KeyMapAuth]; ok {
		ctx.handler(global.contextAuth)
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
		case func(*TriggerContext):
			ctx.handler(actions)
		case []func(*TriggerContext):
			for _,action := range actions {
				ctx.handler(action)
			}
		case TriggerFunc:
			ctx.handler(actions)
		case []TriggerFunc:
			ctx.handler(actions...)
		default:
		}
	}

	ctx.Next()
}

















//自带中间件，参数处理
func (global *triggerGlobal) contextArgs(ctx *TriggerContext) {

	//argn表示参数都可为空
	argn := false
	if v,ok := ctx.Config["argn"].(bool); ok {
		argn = v
	}

	//所有值都会放在 global.Value 中
	err := Mapping.Parse([]string{}, ctx.Config["args"].(Map), ctx.Value, ctx.Args, argn)
	if err != nil {
		ctx.Failed(err)
	} else {
		ctx.Next()
	}
}



//Auth验证处理
func (global *triggerGlobal) contextAuth(ctx *TriggerContext) {

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
func (global *triggerGlobal) contextItem(ctx *TriggerContext) {
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

































//路由执行，found
func (global *triggerGlobal) contextFound(ctx *TriggerContext) {
	//清理执行线
	ctx.cleanup()

	//如果路由配置中有found，就自定义处理
	if v,ok := ctx.Config[KeyMapFound]; ok {
		switch c := v.(type) {
		case TriggerFunc: {
			ctx.handler(c)
		}
		case []TriggerFunc: {
			for _,v := range c {
				ctx.handler(v)
			}
		}
		case func(*TriggerContext): {
			ctx.handler(c)
		}
		case []func(*TriggerContext): {
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
func (global *triggerGlobal) contextError(ctx *TriggerContext) {
	//清理执行线
	ctx.cleanup()

	//如果路由配置中有found，就自定义处理
	if v,ok := ctx.Config[KeyMapError]; ok {
		switch c := v.(type) {
		case TriggerFunc: {
			ctx.handler(c)
		}
		case []TriggerFunc: {
			for _,v := range c {
				ctx.handler(v)
			}
		}
		case func(*TriggerContext): {
			ctx.handler(c)
		}
		case []func(*TriggerContext): {
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


//路由执行，failed
func (global *triggerGlobal) contextFailed(ctx *TriggerContext) {
	//清理执行线
	ctx.cleanup()

	//如果路由配置中有found，就自定义处理
	if v,ok := ctx.Config[KeyMapFailed]; ok {
		switch c := v.(type) {
		case TriggerFunc: {
			ctx.handler(c)
		}
		case []TriggerFunc: {
			for _,v := range c {
				ctx.handler(v)
			}
		}
		case func(*TriggerContext): {
			ctx.handler(c)
		}
		case []func(*TriggerContext): {
			for _,v := range c {
				ctx.handler(v)
			}
		}
		default:
		}
	}


	//handler中的failed
	//用数组保证原始注册顺序
	for _,name := range global.failedHandlerNames {
		ctx.handler(global.failedHandlers[name])
	}

	//最后是默认failed中间件
	ctx.handler(global.failedDefaultHandler)

	ctx.Next()
}



//路由执行，denied
func (global *triggerGlobal) contextDenied(ctx *TriggerContext) {
	//清理执行线
	ctx.cleanup()

	//如果路由配置中有found，就自定义处理
	if v,ok := ctx.Config[KeyMapDenied]; ok {
		switch c := v.(type) {
		case TriggerFunc: {
			ctx.handler(c)
		}
		case []TriggerFunc: {
			for _,v := range c {
				ctx.handler(v)
			}
		}
		case func(*TriggerContext): {
			ctx.handler(c)
		}
		case []func(*TriggerContext): {
			for _,v := range c {
				ctx.handler(v)
			}
		}
		default:
		}
	}

	//handler中的denied
	//用数组保证原始注册顺序
	for _,name := range global.deniedHandlerNames {
		ctx.handler(global.deniedHandlers[name])
	}

	//最后是默认denied中间件
	ctx.handler(global.deniedDefaultHandler)

	ctx.Next()
}



/*
	触发器模块方法 end
*/






/* 默认响应器 begin */
func (global *triggerGlobal) finishResponder(ctx *TriggerContext) {
	//完成就完成了。 不做任何处理
	//因为目前，触发器不需要给调用者响应信息
}

//目前直接调度，可调整，以后做到task中统一调整
//因为万一delay很久。中间正好程序重新或是其它，就丢了
//所以有必要使用task机制重新调度
func (global *triggerGlobal) retriggerResponder(ctx *TriggerContext) {
	body := ctx.Body.(triggerBodyRetrigger)

	time.AfterFunc(body.Delay, func() {
		global.Touch(ctx.Name, ctx.Value)
	})
}
func (global *triggerGlobal) defaultResponder(ctx *TriggerContext) {
	//触发器中，这些好像不需要处理
	//因为目前，触发器不需要给调用者响应信息
}
/* 默认响应器 end */




/* 默认处理器 begin */
func (global *triggerGlobal) foundDefaultHandler(ctx *TriggerContext) {
	//触发器中，这些好像不需要处理
	//因为目前，触发器不需要给调用者响应信息
}
func (global *triggerGlobal) errorDefaultHandler(ctx *TriggerContext) {
	//触发器中，这些好像不需要处理
	//因为目前，触发器不需要给调用者响应信息
}
func (global *triggerGlobal) failedDefaultHandler(ctx *TriggerContext) {
	//触发器中，这些好像不需要处理
	//因为目前，触发器不需要给调用者响应信息
}
func (global *triggerGlobal) deniedDefaultHandler(ctx *TriggerContext) {
	//触发器中，这些好像不需要处理
	//因为目前，触发器不需要给调用者响应信息
}
/* 默认处理器 end */









































/*
	触发器上下文处理方法 begin
*/



//添加执行线
func (ctx *TriggerContext) handler(handlers ...TriggerFunc) {
	for _,handler := range handlers {
		ctx.nexts = append(ctx.nexts, handler)
	}
}
//清空执行线
func (ctx *TriggerContext) cleanup() {
	ctx.next = -1
	ctx.nexts = make([]TriggerFunc, 0)
}

/* 执行下一个 */
func (ctx *TriggerContext) Next() {
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
func (ctx *TriggerContext) Found() {
	ctx.Global.contextFound(ctx)
}
//返回错误
func (ctx *TriggerContext) Error(err *Error) {
	ctx.Wrong = err
	ctx.Global.contextError(ctx)
}

//失败, 就是参数处理失败为主
func (ctx *TriggerContext) Failed(err *Error) {
	ctx.Wrong = err
	ctx.Global.contextFailed(ctx)
}
//拒绝,主要是 auth
func (ctx *TriggerContext) Denied(err *Error) {
	ctx.Wrong = err
	ctx.Global.contextFailed(ctx)
}
/* 上下文处理器 end */







/* 上下文响应器 begin */
//完成操作
func (ctx *TriggerContext) Finish() {
	ctx.Body = triggerBodyFinish{}
}
//重新触发
func (ctx *TriggerContext) Retrigger(delays ...time.Duration) {
	if len(delays) > 0 {
		//延时重新触发
		ctx.Body = triggerBodyRetrigger{ Delay: delays[0] }
	} else {
		//立即重新触发
		ctx.Body = triggerBodyRetrigger{ Delay: time.Second*0 }
	}
}
/* 上下文响应器 end */











/*
	触发器上下文方法 end
*/

