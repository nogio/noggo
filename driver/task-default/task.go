package task_default


import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo"
	"time"
)



type (
	//会话驱动
	DefaultTaskDriver struct {}
	//会话连接
	DefaultTaskConnect struct {
		config Map
	}
)



//返回驱动
func Driver() *DefaultTaskDriver {
	return &DefaultTaskDriver{}
}











//连接任务驱动
func (session *DefaultTaskDriver) Connect(config Map) (noggo.TaskConnect) {
	return  &DefaultTaskConnect{
		config: config,
	}
}












//打开连接
func (connect *DefaultTaskConnect) Open() error {
	return nil
}

//关闭连接
func (connect *DefaultTaskConnect) Close() error {
	return nil
}




//查询会话，
func (connect *DefaultTaskConnect) Accept(id string, name string, delay time.Duration, value Map, call noggo.TaskCall) error {
	time.AfterFunc(delay, func() {
		call(id,name,delay,value)
	})
	return nil
}



//完成任务，从列表中清理
func (connect *DefaultTaskConnect) Finish(id string) error {
	//不做任务处理
	//三方驱动， 应该做一些持久化的处理
	return nil
}