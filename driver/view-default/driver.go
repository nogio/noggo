package view_default



import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo"
)



type (
	//驱动
	DefaultViewDriver struct {}
	//会话连接
	DefaultViewConnect struct {
		config Map
		helpers map[string]Any
	}
)



//返回驱动
func Driver() *DefaultViewDriver {
	return &DefaultViewDriver{}
}











//连接会话驱动
func (driver *DefaultViewDriver) Connect(config Map) (noggo.ViewConnect) {
	return  &DefaultViewConnect{
		config: config, helpers: map[string]Any{},
	}
}












//打开连接
func (connect *DefaultViewConnect) Open() error {
	return nil
}

//关闭连接
func (connect *DefaultViewConnect) Close() error {
	return nil
}


//注册helper
func (connect *DefaultViewConnect) Helper(name string, helper Any) error {
	connect.helpers[name] = helper
	return nil
}




//解析VIEW
func (connect *DefaultViewConnect) Parse(ctx *noggo.HttpContext, name string, model Map, data Map) (error,string) {
	view := newDefaultView(ctx, data)
	return view.Parse(name, model)
}



