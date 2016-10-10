package noggo

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
)

var (
	Nodes map[string]*Noggo

	//当前位置
	//此变量比较重要，是用来注册对象时，表明当前的位置，如当前节点，当前数据库，等等
	Current string

	//全局配置
	Config configConfig

	//常量模块
	Const *constGlobal
	//Map模块
	Mapping	*mappingGlobal

	//日志模块
	Logger *loggerGlobal


	//会话模块
	Session *sessionGlobal

	//触发器模块
	Trigger *triggerGlobal

	//任务模块
	Task *taskGlobal



	/*


	//路由器模块
	Router *routerModule


	//计划模块
	Plan *planModule
	//任务模块
	Task *taskModule
	//http模块
	Http *httpModule
	*/
)



//框架初驱化
func init() {
	Nodes = map[string]*Noggo{}
	//当前位置为空
	Current = ""
	//读取配置文件
	Config = readJsonConfig()

	//常量模块
	Const = &constGlobal{}
	//Mapping模块
	Mapping = &mappingGlobal{}
	//日志模块
	Logger = &loggerGlobal{}
	//会话模块
	Session = &sessionGlobal{}
	//触发器模块
	Trigger = &triggerGlobal{}
	//任务模块
	Task = &taskGlobal{}
}


func readJsonConfig() (configConfig) {
	m := configConfig{}

	bytes, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(fmt.Sprintf("读取配置文件出错:%v", err))
	} else {
		err := json.Unmarshal(bytes, &m)
		if err != nil {
			panic(fmt.Sprintf("解析配置文件出错:%v", err))
		}
	}

	return m
}
