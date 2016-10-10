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


type (
	//计划驱动
	PlanDriver interface {
		Connect(config Map) (PlanConnect)
	}
	//计划连接
	PlanConnect interface {
		//打开连接
		Open() error
		//关闭连接
		Close() error

		//注册计划
		Accept(id,name,time string, call func()) error
		//删除计划
		Remove(id string) error
		//清空计划
		Clear() error

		//开始计划
		Start() error
		//停止计划
		Stop() error
	}
	//计划全局容器
	planGlobal	struct {
		mutex sync.Mutex
		//驱动
		drivers map[string]PlanDriver
	}

)



//计划全局初始化
func (module *planGlobal) init() {
	//计划全局无需处理任何东西
}

//计划全局退出
func (module *planGlobal) exit() {
	//计划全局无需处理任何东西
}


//计划：连接驱动
func (module *planGlobal) connect(config *planConfig) (PlanConnect) {
	if planDriver,ok := module.drivers[config.Driver]; ok {
		return planDriver.Connect(config.Config)
	} else {
		panic("计划：不支持的驱动 " + config.Driver)
	}
}







//---------------------------------------------------------------------------------------------------------------//












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

	//计划上下文
	PlanContext struct {
		Module	*planModule

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