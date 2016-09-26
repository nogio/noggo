package noggo

import (
	. "github.com/nogio/noggo/base"
)

/*
	trigger 触发器模块
*/

type triggerModule struct {

	//路由
	routes map[string]Map			//路由定义
	routeNames []string				//路由名称原始顺序，因为map是无序的
	routeUris map[string]string		//记录所有uris指定
}


/*
	trigger 触发器
	Route  注册路由
*/
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
	//ctx.handler(triggerRequest)
	//ctx.handler(triggerSession)

	//这里做use的拦截器
	//for _,n := range ctx.Service.useNames {
	//	ctx.handler(ctx.Service.UseActions(n)...)
	//}

	//ctx.handler(triggerResponse)
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