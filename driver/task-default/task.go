package task_default


import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo/driver"
	"time"
	"errors"
	"github.com/anacrolix/sync"
)




type (
	//会话驱动
	DefaultTaskDriver struct {}
	//会话连接
	DefaultTaskConnect struct {
		mutex   sync.Mutex
		config      Map
		callback    driver.TaskCallback
		datas       map[string]driver.TaskData
	}
)



//返回驱动
func Driver() *DefaultTaskDriver {
	return &DefaultTaskDriver{}
}











//连接任务驱动
func (session *DefaultTaskDriver) Connect(config Map) (driver.TaskConnect) {
	return  &DefaultTaskConnect{
		config: config, datas: map[string]driver.TaskData{},
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




//注册回调
func (connect *DefaultTaskConnect) Accept(callback driver.TaskCallback) error {
	connect.callback = callback
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





//发起任务
func (connect *DefaultTaskConnect) After(name string, delay time.Duration, value Map) error {

	if connect.callback == nil {
		return errors.New("未注册回调")
	}

	connect.mutex.Lock()
	defer connect.mutex.Unlock()

	//新建任务
	id := NewMd5Id()
	task := driver.TaskData{
		Name: name, Delay: delay, Value: value,
	}

	//保存任务
	connect.datas[id] = task

	//直接回调
	time.AfterFunc(delay, func() {
		connect.callback(id,task.Name, task.Delay, task.Value)
	})

	return nil
}





//完成任务，从列表中清理
func (connect *DefaultTaskConnect) Finish(id string) error {
	connect.mutex.Lock()
	defer connect.mutex.Unlock()

	//从列表中删除
	//三方驱动， 应该做一些持久化的处理
	//更新掉此任务已经完成
	delete(connect.datas, id)

	return nil
}
//重开任务
func (connect *DefaultTaskConnect) Retask(id string, delay time.Duration) error {
	if task,ok := connect.datas[id]; ok {

		//更新一下任务信息
		task.Delay = delay
		connect.datas[id] = task

		time.AfterFunc(delay, func() {
			connect.callback(id, task.Name, task.Delay, task.Value)
		})
	}

	return errors.New("任务不存在")
}