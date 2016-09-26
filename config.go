package noggo



import (
	. "github.com/nogio/noggo/base"
)


type (
	//路由器配置
	routerConfig struct {
		//驱动
		Driver	string	`json:"driver"`
		//自定义配置
		Config	Map		`json:"config"`
	}
	//会话配置
	sessionConfig struct {
		//驱动
		Driver	string	`json:"driver"`
		//超时时间
		Expiry	int64	`json:"expiry"`
		//自定义配置
		Config	Map		`json:"config"`
	}


	//触发器配置
	triggerConfig struct {
		//会话配置
		Session	*sessionConfig	`json:"session"`
		//路由器配置
		Router	*routerConfig	`json:"router"`
	}







	//配置
	configConfig struct {
		//默认路由配置
		Router	*routerConfig	`json:"router"`
		//默认会话配置
		Session	*sessionConfig	`json:"session"`
		//触发器配置
		Trigger *triggerConfig	`json:"trigger"`
	}

)



func (config configConfig) LoadJsonConfig(file string) {

}