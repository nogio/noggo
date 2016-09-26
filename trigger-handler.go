package noggo

import (
	. "github.com/nogio/noggo/base"
)


/*
	trigger 触发器处理器
	请求处理
	包含：route解析、request处理、session处理
*/
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
	ctx.Session = Map{}		//这里应该要去读取或是创建session
	ctx.Sign = &Sign{ ctx.Session }
	ctx.Next()

}
