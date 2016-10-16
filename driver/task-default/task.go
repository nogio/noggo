package task_default


import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo/driver"
	"time"
	"errors"
	"sync"
)




type (
	//驱动
	DefaultTaskDriver struct {}
	//连接
	DefaultTaskConnect struct {
		mutex   sync.Mutex
		config      Map
		callback    driver.TaskAccept

		datas       map[string]TaskData
	}
	//响应对象
	DefaultTaskResponse struct {
		connect *DefaultTaskConnect
	}


	TaskData struct{
		Name    string
		Delay   time.Duration
		Value   Map
	}
)



//返回驱动
func Driver() (driver.TaskDriver) {
	return &DefaultTaskDriver{}
}






//连接驱动
func (session *DefaultTaskDriver) Connect(config Map) (error,driver.TaskConnect) {
	return nil, &DefaultTaskConnect{
		config: config, datas: map[string]TaskData{},
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
func (connect *DefaultTaskConnect) Accept(callback driver.TaskAccept) error {
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


	//新建任务
	id := NewMd5Id()
	task := TaskData{
		Name: name, Delay: delay, Value: value,
	}

	//保存任务
	connect.mutex.Lock()
	connect.datas[id] = task
	connect.mutex.Unlock()

	//直接回调
	time.AfterFunc(delay, func() {
		connect.execute(id, task.Name, task.Delay, task.Value)
	})

	return nil
}



//执行统一到这里
func (connect *DefaultTaskConnect) execute(id string, name string, delay time.Duration, value Map) {
	req := &driver.TaskRequest{ Id: id, Name: name, Delay: delay, Value: value }
	res := &DefaultTaskResponse{ connect }
	connect.callback(req, res)
}






//完成任务，从列表中清理
func (connect *DefaultTaskConnect) finish(id string) error {
	//从数据中删除
	connect.mutex.Lock()
	delete(connect.datas, id)
	connect.mutex.Unlock()

	return nil
}
//重开任务
func (connect *DefaultTaskConnect) retask(id string, delay time.Duration) error {
	if task,ok := connect.datas[id]; ok {

		//更新一下任务信息
		connect.mutex.Lock()
		task.Delay = delay
		connect.datas[id] = task
		connect.mutex.Unlock()

		time.AfterFunc(delay, func() {
			connect.execute(id, task.Name, task.Delay, task.Value)
		})
	}

	return errors.New("任务不存在")
}





























//完成任务，从列表中清理
func (res *DefaultTaskResponse) Finish(id string) error {
	return res.connect.finish(id)
}
//重开任务
func (res *DefaultTaskResponse) Retask(id string, delay time.Duration) error {
	return res.connect.retask(id, delay)
}