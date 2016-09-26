package noggo


import (
	"time"
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo/driver"
	"fmt"
)


type (
	//触发器函数
	TriggerCall func(*TriggerCtx)
	TriggerMatch func(*TriggerCtx) bool

	//触发器模块
	triggerModule struct {
		//会话连接
		session			driver.SessionConnect
		sessionConfig	*sessionConfig

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
		Code	int			//返回的状态码
		Type	Type		//响应类型
		Body	Any			//响应内容
		Error	*Error		//响应错误
	}

)



/*
	触发器模块方法 begin
*/


//触发器初始化
func (trigger *triggerModule) init() {
	trigger.initSession()
}
//初始化会话驱动
func (trigger *triggerModule) initSession() {
	//先拿到默认的会话配置
	trigger.sessionConfig = Config.Session
	//如果触发器有单独定义会话配置，则使用
	if Config.Trigger.Session != nil {
		trigger.sessionConfig = Config.Trigger.Session
	}

	trigger.session = Session.connect(trigger.sessionConfig)
	if trigger.session != nil {
		trigger.session.Open()
	} else {
		panic("触发器连接会话服务失败")
	}
}

//触发器退出
func (trigger *triggerModule) exit() {
}


















//触发器，注册路由
func (trigger *triggerModule) Route(name string, config Map) {
	//保存配置
	trigger.routes[name] = config
	trigger.routeNames = append(trigger.routeNames, name)

	//处理uri
	trigger.routeUris[name] = name
	if v,ok := config[KeyConfigUri]; ok {

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

	//前置处理
	//ctx.handler(triggerRoute)
	ctx.handler(trigger.handlerRequest)
	//ctx.handler(triggerSession)

	//这里做use的拦截器
	//for _,n := range ctx.Service.useNames {
	//	ctx.handler(ctx.Service.UseActions(n)...)
	//}

	ctx.handler(trigger.handlerResponse)
	//ctx.handler(triggerExecute)

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

	fmt.Printf("handlerRequest, path=%v\n", ctx.Path)


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
	ctx.Session = trigger.session.Create(ctx.Id, trigger.sessionConfig.Expiry)
	ctx.Sign = &Sign{ ctx.Session }
	ctx.Next()
	trigger.session.Update(ctx.Id, ctx.Session, trigger.sessionConfig.Expiry)
}


func (trigger *triggerModule) handlerResponse(ctx *TriggerCtx) {
	fmt.Printf("session=%v\n", ctx.Session)

	ctx.Session["now1"] = time.Now()

	fmt.Printf("session=%v\n", ctx.Session)

	ctx.Session["now2"] = time.Now()

	fmt.Printf("session=%v\n", ctx.Session)
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
		//没有了，不再执行，Response会处理为404
	}
}


//触发器响应
//完成操作
func (ctx *TriggerCtx) Finish() {
}

//触发器响应
//重新触发
func (ctx *TriggerCtx) Retrigger(delays ...time.Duration) {
	if len(delays) > 0 {
		//延时重新触发


	} else {
		//立即重新触发

	}
}

/*
	触发器上下文方法 end
*/

