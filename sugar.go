package noggo


import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo/driver"
	"fmt"
)



type (
	sugarGlobal struct {
		https   map[string]Map
	}
)


//初始化，把any,get,post啥的，写入HTTP
func (global *sugarGlobal) init() {
	for uri,route := range global.https {
		Http.Route(NewMd5Id(), Map{
			"uri": uri, "route": route,
		})
	}
}
func (global *sugarGlobal) exit() {
	//没东西需要处理
}





//返回状态字定以及字串
func (global *sugarGlobal) States(args ...string) (Map) {
	m := Map{}

	if len(args) > 0 {
		for _,n := range args {
			m[fmt.Sprintf("%d", Const.StateCode(n))] = Const.LangString(n)
		}
	} else {
		for k,v := range Const.states {
			m[fmt.Sprintf("%d", v)] = Const.LangString(k)
		}
	}

	return m
}












//注册驱动
//方便注册各种驱动
func Driver(name string, drv Any) {
	switch v := drv.(type) {
	case driver.LoggerDriver:
		Logger.Driver(name, v)
	case driver.SessionDriver:
		Session.Driver(name, v)

		//触发器没有驱动
	case driver.TaskDriver:
		Task.Driver(name, v)

	case driver.PlanDriver:
		Plan.Driver(name, v)
	case driver.EventDriver:
		Event.Driver(name, v)
	case driver.QueueDriver:
		Queue.Driver(name, v)


	case driver.CacheDriver:
		Cache.Driver(name, v)
	case driver.DataDriver:
		Data.Driver(name, v)

	case driver.HttpDriver:
		Http.Driver(name, v)
	case driver.ViewDriver:
		View.Driver(name, v)

	default:
		panic("未支持的驱动")
	}
}
//语法大糖
func Drv(name string, drv Any) {
	Driver(name, drv)
}



//注册中间件
//方便注册各种中间件
func Middler(call Any) {
	switch v := call.(type) {

	case TriggerFunc:
		Trigger.Middler(NewMd5Id(), v)
	case func(*TriggerContext):
		Trigger.Middler(NewMd5Id(), v)

	case TaskFunc:
		Task.Middler(NewMd5Id(), v)
	case func(*TaskContext):
		Task.Middler(NewMd5Id(), v)

	case PlanFunc:
		Plan.Middler(NewMd5Id(), v)
	case func(*PlanContext):
		Plan.Middler(NewMd5Id(), v)



	case EventFunc:
		Event.Middler(NewMd5Id(), v)
	case func(*EventContext):
		Event.Middler(NewMd5Id(), v)
	case QueueFunc:
		Queue.Middler(NewMd5Id(), v)
	case func(*QueueContext):
		Queue.Middler(NewMd5Id(), v)


	case HttpFunc:
		Http.Middler(NewMd5Id(), v)
	case func(*HttpContext):
		Http.Middler(NewMd5Id(), v)


	default:
		panic("未支持的中间件")
	}
}
//语法大糖
func Use(call Any) {
	Middler(call)
}







//---------------------------------------------------------------------



//通用添加方法，用于添加task,plan,event,queue,http等
func Add(data string, call Any) {

	switch v := call.(type) {

	case TriggerFunc:
		addTrigger(data, v)
	case func(*TriggerContext):
		addTrigger(data, v)

	case TaskFunc:
		addTask(data, v)
	case func(*TaskContext):
		addTask(data, v)

	case PlanFunc:
		addPlan(data, v)
	case func(*PlanContext):
		addPlan(data, v)



	case EventFunc:
		addEvent(data, v)
	case func(*EventContext):
		addEvent(data, v)
	case QueueFunc:
		addQueue(data, v)
	case func(*QueueContext):
		addQueue(data, v)


	case HttpFunc:
		addHttp(data, v)
	case func(*HttpContext):
		addHttp(data, v)


	default:
		panic("未支持的路由")
	}
}
func addTrigger(name string, call TriggerFunc) {
	Trigger.Route(name, Map{
		"route": Map{
			"name": name, "text": name,
			"action": call,
		},
	})
}
func addTask(name string, call TaskFunc) {
	Task.Route(name, Map{
		"route": Map{
			"name": name, "text": name,
			"action": call,
		},
	})
}
func addPlan(time string, call PlanFunc) {
	Plan.Route(NewMd5Id(), Map{
		"time": time,
		"route": Map{
			"name": time, "text": time,
			"action": call,
		},
	})
}
func addEvent(name string, call EventFunc) {
	Event.Route(name, Map{
		"route": Map{
			"name": name, "text": name,
			"action": call,
		},
	})
}
func addQueue(name string, call QueueFunc) {
	Queue.Route(name, Map{
		"line": 1,
		"route": Map{
			"name": name, "text": name,
			"action": call,
		},
	})
}
func addHttp(name string, call HttpFunc) {
	Http.Route(name, Map{
		"line": 1,
		"route": Map{
			"name": name, "text": name,
			"action": call,
		},
	})
}




//---------------------------------------------------------

// Get, Post, Put, Patch, Delete

/*
//http Any
func Any(uri string, call *HttpContext) {
	Sugar.https[uri] = Map{
		"name": "", "text": "",
		"action": call,
	}
}
*/
//http Get
func Get(uri string, call HttpFunc) {
	if v,ok := Sugar.https[uri]; ok {
		v["get"] = Map{
			"name": "", "text": "",
			"action": call,
		}
	} else {
		Sugar.https[uri] = Map{
			"get": Map{
				"name": "", "text": "",
				"action": call,
			},
		}
	}
}
//http Post
func Post(uri string, call HttpFunc) {
	if v,ok := Sugar.https[uri]; ok {
		v["post"] = Map{
			"name": "", "text": "",
			"action": call,
		}
	} else {
		Sugar.https[uri] = Map{
			"post": Map{
				"name": "", "text": "",
				"action": call,
			},
		}
	}
}
//http Put
func Put(uri string, call HttpFunc) {
	if v,ok := Sugar.https[uri]; ok {
		v["put"] = Map{
			"name": "", "text": "",
			"action": call,
		}
	} else {
		Sugar.https[uri] = Map{
			"put": Map{
				"name": "", "text": "",
				"action": call,
			},
		}
	}
}
//http Patch
func Patch(uri string, call HttpFunc) {
	if v,ok := Sugar.https[uri]; ok {
		v["patch"] = Map{
			"name": "", "text": "",
			"action": call,
		}
	} else {
		Sugar.https[uri] = Map{
			"patch": Map{
				"name": "", "text": "",
				"action": call,
			},
		}
	}
}
//http Delete
func Delete(uri string, call HttpFunc) {
	if v,ok := Sugar.https[uri]; ok {
		v["delete"] = Map{
			"name": "", "text": "",
			"action": call,
		}
	} else {
		Sugar.https[uri] = Map{
			"delete": Map{
				"name": "", "text": "",
				"action": call,
			},
		}
	}
}
