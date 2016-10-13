package task_default


import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo"
	"time"
	"errors"
)



type (
	//会话驱动
	DefaultTaskDriver struct {}
	//会话连接
	DefaultTaskConnect struct {
		config Map
		tasks map[string]noggo.TaskAcceptFunc
	}
)



//返回驱动
func Driver() *DefaultTaskDriver {
	return &DefaultTaskDriver{}
}











//连接任务驱动
func (session *DefaultTaskDriver) Connect(config Map) (noggo.TaskConnect) {
	return  &DefaultTaskConnect{
		config: config, tasks: map[string]noggo.TaskAcceptFunc{},
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
func (connect *DefaultTaskConnect) Accept(name string, call noggo.TaskAcceptFunc) error {
	connect.tasks[name] = call
	return nil
}



//开始
func (connect *DefaultTaskConnect) Start() error {
	return nil
}

//结束
func (connect *DefaultTaskConnect) Stop() error {
	return nil
}





//触发任务
func (connect *DefaultTaskConnect) Touch(id string, name string, delay time.Duration, value Map) error {

	if call,ok := connect.tasks[name]; ok {
		time.AfterFunc(delay, func() {
			call(id,name,delay,value)
		})

		return nil

	} else {
		return errors.New("不支持的任务")
	}
}





//完成任务，从列表中清理
func (connect *DefaultTaskConnect) Finish(id string) error {
	//不做任务处理
	//三方驱动， 应该做一些持久化的处理
	return nil
}