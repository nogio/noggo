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
		handler    driver.TaskHandler

		names       []string
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
func (session *DefaultTaskDriver) Connect(config Map) (driver.TaskConnect,error) {
	return &DefaultTaskConnect{
		config: config, datas: map[string]TaskData{},
	},nil
}


//打开连接
func (connect *DefaultTaskConnect) Open() error {
	return nil
}

//关闭连接
func (connect *DefaultTaskConnect) Close() error {
	return nil
}



//订阅者注册事件
func (con *DefaultTaskConnect) Accept(name string) error {
	con.mutex.Lock()
	defer con.mutex.Unlock()

	//注意，这里应该处理唯一的问题，如果已经存在某name，就跳过了得
	con.names = append(con.names, name)

	return nil
}


//开始
func (con *DefaultTaskConnect) Start(handler driver.TaskHandler) error {
	con.handler = handler
	return nil
}



//发起任务
func (connect *DefaultTaskConnect) After(name string, delay time.Duration, value Map) error {

	if connect.handler == nil {
		return errors.New("无任务处理器")
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
	connect.handler(req, res)
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