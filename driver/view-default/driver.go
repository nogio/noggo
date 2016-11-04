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
		config Map
		root string
	}
)



//返回驱动
func Driver(roots ...string) (noggo.ViewDriver) {
	root := "views"
	if len(roots) > 0 {
		root = roots[0]
	}

	return &DefaultViewDriver{
		root: root,
	}
}











//连接驱动
func (driver *DefaultViewDriver) Connect(config Map) (noggo.ViewConnect,error) {
	return &DefaultViewConnect{
		config: config, root: driver.root,
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
	view := newDefaultView(connect.root, parse)
	return view.Parse()
}



