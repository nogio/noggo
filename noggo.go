package noggo


import (
	. "github.com/nogio/noggo/base"
)

type (
	Noggo struct {
		running bool

		//节点名称和唯一标识
		Id		string
		Name	string
		Port	string
		Config	*nodeConfig

		session *sessionModule

		Plan	*planModule
		Event   *eventModule
		Queue   *queueModule
		Http	*httpModule
	}
)



//创建新节点
func New(names ...string) (*Noggo) {

	name := ConstNodeGlobal
	if len(names) > 0 {
		name = names[0]
	}


	//如果已经实例过了, 直接返回
	if nodes[name] != nil {
		return nodes[name]
	}


	node := &Noggo{}

	if config,ok := Config.Node[name]; ok {
		node.Id = config.Id
		node.Name = name
		node.Port = config.Port
		node.Config = config
	} else {
		node.Id = name
		node.Name = name
		node.Port = ":8080"
		node.Config = &nodeConfig{
			Charset:"", Cookie:"noggo", Domain:"",
		}
	}

	//模块们
	node.session = newSessionModule(node)
	node.Plan = newPlanModule(node)
	node.Event = newEventModule(node)
	node.Queue = newQueueModule(node)
	node.Http = newHttpModule(node)


	//加入节点列表
	nodes[name] = node

	return node
}




//启动节点
func (node *Noggo) Run(ports ...string) {
	if len(ports) > 0 {
		node.Port = ports[0]
	}


	//如果还没初妈化， 先初始化
	if initialized == false {
		Init()
	}

	if node.running == false {

		node.session.run()
		node.Plan.run()
		node.Event.run()
		node.Queue.run()
		node.Http.run()

		node.running = true

		//如果是直接运行， 就监听退出信号
		if node.Name == ConstNodeGlobal {
			Logger.Info("noggo", "is running at", node.Port)
			Exit()
		} else {
			Logger.Info("node", node.Name, node.Id, "is running at", node.Port)
		}
	}
}
//结束节点
func (node *Noggo) End() {

	if node.Name == ConstNodeGlobal {
		Logger.Info("noggo", "is ending")
	} else {
		Logger.Info("node", node.Name, node.Id, "is ending")
	}

	node.Http.end()
	node.Queue.end()
	node.Event.end()
	node.Plan.end()
	node.session.end()
}



/* 注册处理器 end */




//---------------------------------------------------------- 语法糖 begin ----------------------------------------------------------


//注册中间件
//用Any类似，是方便，为HTTP，PLAN，等等不同的模块直接注册中间件
func (node *Noggo) Use(call Any) {
	switch v := call.(type) {
	case HttpFunc:
		node.Http.Use(v)
	case func(*HttpContext):
		node.Http.Use(v)

	case PlanFunc:
		node.Plan.Use(v)
	case func(*PlanContext):
		node.Plan.Use(v)

	case EventFunc:
		node.Event.Use(v)
	case func(*EventContext):
		node.Event.Use(v)
	}

}




//添加路由
func (node *Noggo) Add(name string, call Any) {
	switch v := call.(type) {

	case TriggerFunc:
		Trigger.Add(name, v)
	case func(*TriggerContext):
		Trigger.Add(name, v)

	case TaskFunc:
		Task.Add(name, v)
	case func(*TaskContext):
		Task.Add(name, v)



	case HttpFunc:
		node.Http.Any(name, v)
	case func(*HttpContext):
		node.Http.Any(name, v)

	case PlanFunc:
		node.Plan.Add(name, v)
	case func(*PlanContext):
		node.Plan.Add(name, v)

	case EventFunc:
		node.Event.Add(name, v)
	case func(*EventContext):
		node.Event.Add(name, v)
	}
}


//注册Any方法
func (node *Noggo) Any(path string, call HttpFunc) {
	node.Http.Any(path, call)
}
//注册get方法
func (node *Noggo) Get(path string, call HttpFunc) {
	node.Http.Get(path, call)
}
//注册post方法
func (node *Noggo) Post(path string, call HttpFunc) {
	node.Http.Post(path, call)
}
//注册put方法
func (node *Noggo) Put(path string, call HttpFunc) {
	node.Http.Put(path, call)
}
//注册delete方法
func (node *Noggo) Delete(path string, call HttpFunc) {
	node.Http.Put(path, call)
}
//---------------------------------------------------------- 语法糖 end ----------------------------------------------------------



