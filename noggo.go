package noggo

import "fmt"

type (
	Noggo struct {
		//节点名称和唯一标识
		Id		string
		Name	string
		Port	string
		Config	nodeConfig

		Plan	*planModule
		Http	*httpModule
	}
)



//创建新节点
func New(name string) (*Noggo) {

	//所有节点实例，都保存到全局变量中
	//已经存在了， 就直接返回， 每一个进程一个节点。 只启动一份
	//如果已经实例过了, 直接返回
	if nodes[name] != nil {
		return nodes[name]
	}

	if config,ok := Config.Node[name]; ok {

		node := &Noggo{
			Id: config.Id, Name: name, Port: config.Port, Config: config,
		}

		//计划
		node.Plan = newPlanModule(node)
		node.Http = newHttpModule(node)

		//加入节点列表
		nodes[name] = node

		return node

	} else {
		panic(fmt.Sprintf("节点: 不存在 %v", name))
	}
}




//启动节点
func (node *Noggo) Run() {
	node.Plan.run()
	node.Http.run()

	Logger.Info("node", node.Name, node.Id, "is running at", node.Port)
}
//结束节点
func (node *Noggo) End() {

	Logger.Info("node", node.Name, node.Id, "is ending")

	node.Http.end()
	node.Plan.end()
}