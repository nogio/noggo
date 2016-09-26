package noggo

import (
	. "github.com/nogio/noggo/base"
	"time"
)





//触发器上下文
type TriggerCtx struct {
	nexts []TriggerCall		//上下文执行线
	next int				//下一个执行索引

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
