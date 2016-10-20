package queue_default

/*
	默认队列驱动
	直接基于内存的队列，这样只在单进程有效，不支持跨进程队列
*/


import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo/driver"
	"errors"
	"time"
	"sync"
)



type (
	//驱动
	DefaultQueueDriver struct {
	}
	//连接
	DefaultQueueConnect struct {
		config      Map
		handler     driver.QueueHandler

		mutex       sync.Mutex
		names       []string
		datas      map[string]QueueData
	}

	//响应对象
	DefaultQueueResponse struct {
		con *DefaultQueueConnect
	}


	QueueData struct{
		Name    string
		Time    string
		Value   Map
	}
)


var (
	msg *Msg
)

func init() {
	msg = &Msg{ subs: map[string][]MsgFunc{} }
}






//返回驱动
func Driver() (driver.QueueDriver) {
	return &DefaultQueueDriver{}
}




//连接
func (drv *DefaultQueueDriver) Connect(config Map) (driver.QueueConnect,error) {
	return &DefaultQueueConnect{
		config: config, names: []string{}, datas: map[string]QueueData{},
	}, nil
}




//打开连接
func (drv *DefaultQueueConnect) Open() error {
	return nil
}
//关闭连接
func (drv *DefaultQueueConnect) Close() error {
	return nil
}






//订阅者注册队列
func (con *DefaultQueueConnect) Accept(name string) error {
	con.mutex.Lock()
	defer con.mutex.Unlock()

	//注意，这里应该处理唯一的问题，如果已经存在某name，就跳过了得
	con.names = append(con.names, name)

	return nil
}


//发布者发布消息
func (con *DefaultQueueConnect) Publish(name string, value Map) error {
	msg.Pub(name, value)
	return nil
}








//开始订阅者
func (con *DefaultQueueConnect) Subscriber(handler driver.QueueHandler) error {
	con.handler = handler

	//订阅消息
	for _,name := range con.names {
		msg.Sub(name, func(value Map) {

			//新建计划
			id := NewMd5Id()
			queue := QueueData{
				Name: name, Value: Map{},
			}

			//保存计划
			con.mutex.Lock()
			con.datas[id] = queue
			con.mutex.Unlock()

			//调用队列
			con.execute(id, queue.Name, queue.Value)
		})
	}

	
	return nil
}
//开始发布者
func (con *DefaultQueueConnect) Publisher() error {

	//发布者貌似不需要干什么？

	return nil
}








//执行统一到这里
func (con *DefaultQueueConnect) execute(id string, name string, value Map) {
	req := &driver.QueueRequest{ Id: id, Name: name, Value: value }
	res := &DefaultQueueResponse{ con }
	con.handler(req, res)
}
















//完成队列，从列表中清理
func (con *DefaultQueueConnect) finish(id string) error {
	con.mutex.Lock()
	defer con.mutex.Unlock()

	delete(con.datas, id)
	return nil
}
//重开队列
func (con *DefaultQueueConnect) requeue(id string, delay time.Duration) error {
	if queue,ok := con.datas[id]; ok {
		time.AfterFunc(delay, func() {
			con.execute(id, queue.Name, queue.Value)
		})
	}

	return errors.New("计划不存在")
}











//响应接口对象



//完成队列，从列表中清理
func (res *DefaultQueueResponse) Finish(id string) error {
	return res.con.finish(id)
}
//重开队列
func (res *DefaultQueueResponse) Requeue(id string, delay time.Duration) error {
	return res.con.requeue(id, delay)
}