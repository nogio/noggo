package view_default



import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo"
)



type (
	//驱动
	DefaultViewDriver struct {
		root string
	}
	//会话连接
	DefaultViewConnect struct {
		config *DefaultViewConfig
	}
	DefaultViewConfig struct {
		Root    string
		Left    string
		Right   string
	}
)



//返回驱动
func Driver(roots ...string) (noggo.ViewDriver) {
	root := ""
	if len(roots) > 0 {
		root = roots[0]
	}

	return &DefaultViewDriver{
		root: root,
	}
}











//连接驱动
func (driver *DefaultViewDriver) Connect(config Map) (noggo.ViewConnect,error) {

	//获取配置信息
	cfg := &DefaultViewConfig{
		Root: "views", Left: "{%", Right: "%}",
	}

	//VIEW目录
	if v,ok := config["root"].(string); ok {
		cfg.Root = v
	}
	//左
	if v,ok := config["left"].(string); ok {
		cfg.Left = v
	}
	//右
	if v,ok := config["right"].(string); ok {
		cfg.Right = v
	}

	//driver中的 root优先
	if driver.root != "" {
		cfg.Root = driver.root
	}


	return &DefaultViewConnect{
		config: cfg,
	},nil
}












//打开连接
func (connect *DefaultViewConnect) Open() error {
	return nil
}

//关闭连接
func (connect *DefaultViewConnect) Close() error {
	return nil
}




//解析VIEW
func (connect *DefaultViewConnect) Parse(parse *noggo.ViewParse) (string,error) {
	view := newDefaultView(connect.config, parse)
	return view.Parse()
}



