package noggo


import (
	"time"
	. "github.com/nogio/noggo/base"
)


type (
	//触发器函数
	TriggerFunc func(*TriggerContext)
	TriggerMatch func(*TriggerContext) bool
	TriggerAcceptFunc func()



	//触发器模块
	triggerModule struct {
												 //路由器连接
		routerConfig	*routerConfig
		routerConnect	RouterConnect
												 //会话连接
		sessionConfig	*sessionConfig
		sessionConnect	SessionConnect


												 //路由
		routes 		map[string]Map			//路由定义
		routeNames	[]string				//路由名称原始顺序，因为map是无序的
		routeUris 	map[string]string		//记录所有uris指定	map[uri]name

												 //拦截器们
		requestFilters, executeFilters, responseFilters map[string]TriggerFunc
		requestFilterNames, executeFilterNames, responseFilterNames []string

												 //处理器们
		foundHandlers, errorHandlers, failedHandlers, deniedHandlers map[string]TriggerFunc
		foundHandlerNames, errorHandlerNames, failedHandlerNames, deniedHandlerNames []string
	}

	//触发器上下文
	TriggerContext struct {
		Module	*triggerModule
								 //执行线
		nexts []TriggerFunc		//方法列表
		next int				//下一个索引

								 //基础
		Id	string			//Session Id  会话时使用
		Session Map			//存储Session值
		Sign	*Sign		//签名功能，基于session

								 //请求相关
		Method	string		//请求的method， 继承之web请求， 暂时无用
		Path	string		//请求的路径，演变自web， 暂时等于trigger的名称
		Lang	string		//当前上下文的语言，默认应为default


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
		Code	int			//返回的状态码
		Type	string		//响应类型
		Body	Any			//响应内容

		Wrong	*Error		//错误信息
	}
)



/*
	触发器模块方法 begin
*/




//触发器初始化
func (module *triggerModule) init() {
	module.initRouter()
	module.initSession()
}
//初始化路由驱动
func (module *triggerModule) initRouter() {
	if Config.Trigger.Router == nil {
		//使用默认的由路连接
		module.routerConfig = Router.routerConfig
		module.routerConnect = Router.routerConnect
	} else {
		//使用自定义的由路连接
		module.routerConfig = Config.Trigger.Router
		module.routerConnect = Router.connect(module.routerConfig)

		if module.routerConnect == nil {
			panic("触发器连接路由服务失败")
		} else {
			err := module.routerConnect.Open()
			if err != nil {
				panic("触发器打开路由服务失败 " + err.Error())
			}
		}
	}
}

//初始化会话驱动
func (module *triggerModule) initSession() {
	if Config.Trigger.Session == nil {
		//使用默认的会话连接
		module.sessionConfig = Session.sessionConfig
		module.sessionConnect = Session.sessionConnect
	} else {
		//使用自定义的会话连接
		module.sessionConfig = Config.Trigger.Session
		module.sessionConnect = Session.connect(module.sessionConfig)


		if module.sessionConnect == nil {
			panic("触发器连接会话服务失败")
		} else {
			//打开会话连接
			err := module.sessionConnect.Open()
			if err != nil {
				panic("触发器打开会话服务失败 " + err.Error())
			}
		}
	}
}




//触发器退出
func (module *triggerModule) exit() {
	module.exitRouter()
	module.exitSession()
}
//触发器退出，路由器
func (module *triggerModule) exitRouter() {
	//关闭路由
	if module.routerConnect != nil {
		module.routerConnect.Close()
		module.routerConnect = nil
	}
}
//触发器退出，会话
func (module *triggerModule) exitSession() {
	//关闭会话
	if module.sessionConnect != nil {
		module.sessionConnect.Close()
		module.sessionConnect = nil
	}
}













//注册路由
func (module *triggerModule) Route(name string, config Map) {
	//保存配置
	module.routes[name] = config
	module.routeNames = append(module.routeNames, name)

	//处理uri
	module.routeUris[name] = name
	if v,ok := config[KeyMapUri]; ok {

		switch uris := v.(type) {
		case string:
			module.routeUris[uris] = name
		case []string:
			for _,uri := range uris {
				module.routeUris[uri] = name
			}
		}
	}
}








/* 注册拦截器 begin */
func (module *triggerModule) RequestFilter(name string, call TriggerFunc) {

	if module.requestFilters == nil {
		module.requestFilters = make(map[string]TriggerFunc)
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
func (module *triggerModule) ExecuteFilter(name string, call TriggerFunc) {

	if module.executeFilters == nil {
		module.executeFilters = make(map[string]TriggerFunc)
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
func (module *triggerModule) ResponseFilter(name string, call TriggerFunc) {

	if module.responseFilters == nil {
		module.responseFilters = make(map[string]TriggerFunc)
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
func (module *triggerModule) FoundHandler(name string, call TriggerFunc) {

	if module.foundHandlers == nil {
		module.foundHandlers = make(map[string]TriggerFunc)
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
func (module *triggerModule) ErrorHandler(name string, call TriggerFunc) {

	if module.errorHandlers == nil {
		module.errorHandlers = make(map[string]TriggerFunc)
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
func (module *triggerModule) FailedHandler(name string, call TriggerFunc) {

	if module.failedHandlers == nil {
		module.failedHandlers = make(map[string]TriggerFunc)
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
func (module *triggerModule) DeniedHandler(name string, call TriggerFunc) {

	if module.deniedHandlers == nil {
		module.deniedHandlers = make(map[string]TriggerFunc)
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










//创建Trigger上下文
func (module *triggerModule) newTriggerContext(method, path string, value Map) (*TriggerContext) {
	return &TriggerContext{
		Module: module,

		next: -1, nexts: []TriggerFunc{},

		Method: method, Path: path,
		Name: "", Config: Map{}, Branchs:[]Map{},

		Params: Map{}, Value: value, Locals: Map{},
		Args: Map{}, Items: Map{}, Auths: Map{},
	}
}



//触发器Trigger  请求开始
func (module *triggerModule) serveTrigger(method, path string, value Map) {
	ctx := module.newTriggerContext(method, path, value)

	//请求处理
	//filter中的request
	//用数组保证原始注册顺序
	for _,name := range module.requestFilterNames {
		ctx.handler(module.requestFilters[name])
	}
	ctx.handler(module.contextRequest)

	//响应处理
	ctx.handler(module.contextResponse)
	//filter中的response
	//用数组保证原始注册顺序
	for _,name := range module.responseFilterNames {
		ctx.handler(module.responseFilters[name])
	}

	//开始执行
	ctx.handler(module.contextExecute)
	ctx.Next()
}














//发起触发
func (module *triggerModule) Touch(path string, args ...Map) {

	value := Map{}
	if len(args) > 0 {
		value = args[0]
	}

	go module.serveTrigger("", path, value)
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
func (module *triggerModule) contextRequest(ctx *TriggerContext) {

	//路由解析
	//目前暂不支持driver
	//直接使用name相等就匹配
	if name,ok := module.routeUris[ctx.Path]; ok {
		ctx.Name = name
		ctx.Config = module.routes[name]
	} else {
		ctx.Config = nil
	}


	//请求处理
	//主要是SessionId处理、处理传过来的值或表单
	ctx.Id = ctx.Name	//使用name做为id，以便在同一个触发器之下共享session

	//会话处理
	ctx.Session = module.sessionConnect.Create(ctx.Id, module.sessionConfig.Expiry)
	ctx.Sign = &Sign{ ctx.Session }
	ctx.Next()
	module.sessionConnect.Update(ctx.Id, ctx.Session, module.sessionConfig.Expiry)
}

//处理响应
func (module *triggerModule) contextResponse(ctx *TriggerContext) {
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
func (module *triggerModule) contextExecute(ctx *TriggerContext) {

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


//触发器处理：处理分支
func (module *triggerModule) contextBranch(ctx *TriggerContext) {

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
func (module *triggerModule) contextArgs(ctx *TriggerContext) {

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
func (module *triggerModule) contextAuth(ctx *TriggerContext) {

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
		ctx.Auths = saveMap
	}

	ctx.Next()
}
//Entity实体处理
func (module *triggerModule) contextItem(ctx *TriggerContext) {
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

		ctx.Items = saveMap
	}
	ctx.Next()
}

































//路由执行，found
func (module *triggerModule) contextFound(ctx *TriggerContext) {
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
	for _,name := range module.foundHandlerNames {
		ctx.handler(module.foundHandlers[name])
	}

	//最后是默认found中间件
	ctx.handler(module.foundDefaultHandler)

	ctx.Next()
}


//路由执行，error
func (module *triggerModule) contextError(ctx *TriggerContext) {
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
	for _,name := range module.errorHandlerNames {
		ctx.handler(module.errorHandlers[name])
	}

	//最后是默认error中间件
	ctx.handler(module.errorDefaultHandler)

	ctx.Next()
}


//路由执行，failed
func (module *triggerModule) contextFailed(ctx *TriggerContext) {
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
	for _,name := range module.failedHandlerNames {
		ctx.handler(module.failedHandlers[name])
	}

	//最后是默认failed中间件
	ctx.handler(module.failedDefaultHandler)

	ctx.Next()
}



//路由执行，denied
func (module *triggerModule) contextDenied(ctx *TriggerContext) {
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
	for _,name := range module.deniedHandlerNames {
		ctx.handler(module.deniedHandlers[name])
	}

	//最后是默认denied中间件
	ctx.handler(module.deniedDefaultHandler)

	ctx.Next()
}



/*
	触发器模块方法 end
*/






/* 默认响应器 begin */
func (module *triggerModule) finishDefaultResponder(ctx *TriggerContext) {
	//触发器中，这些好像不需要处理
}
func (module *triggerModule) retryDefaultResponder(ctx *TriggerContext) {
	//触发器中，这些好像不需要处理
}
/* 默认响应器 end */




/* 默认处理器 begin */
func (module *triggerModule) foundDefaultHandler(ctx *TriggerContext) {
	//触发器中，这些好像不需要处理
}
func (module *triggerModule) errorDefaultHandler(ctx *TriggerContext) {
	//触发器中，这些好像不需要处理
}
func (module *triggerModule) failedDefaultHandler(ctx *TriggerContext) {
	//触发器中，这些好像不需要处理
}
func (module *triggerModule) deniedDefaultHandler(ctx *TriggerContext) {
	//触发器中，这些好像不需要处理
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
	ctx.Module.contextFound(ctx)
}
//返回错误
func (ctx *TriggerContext) Error(err *Error) {
	ctx.Wrong = err
	ctx.Module.contextError(ctx)
}

//失败, 就是参数处理失败为主
func (ctx *TriggerContext) Failed(err *Error) {
	ctx.Wrong = err
	ctx.Module.contextFailed(ctx)
}
//拒绝,主要是 auth
func (ctx *TriggerContext) Denied(err *Error) {
	ctx.Wrong = err
	ctx.Module.contextFailed(ctx)
}
/* 上下文处理器 end */







/* 上下文响应器 begin */
//完成操作
func (ctx *TriggerContext) Finish() {
	ctx.Body = BodyTriggerFinish{}
}
//重新触发
func (ctx *TriggerContext) Retrigger(delays ...time.Duration) {
	if len(delays) > 0 {
		//延时重新触发
		ctx.Body = BodyTriggerRetrigger{ Delay: delays[0] }
	} else {
		//立即重新触发
		ctx.Body = BodyTriggerRetrigger{ Delay: time.Second*0 }
	}
}
/* 上下文响应器 end */











/*
	触发器上下文方法 end
*/

