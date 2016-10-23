package noggo


import (
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

