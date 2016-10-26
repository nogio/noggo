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
		names       map[string]int  //map[name]line
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
	msg = &Msg{ subs: map[string]chan Map{} }
}






//返回驱动
func Driver() (driver.QueueDriver) {
	return &DefaultQueueDriver{}
}




//连接
func (drv *DefaultQueueDriver) Connect(config Map) (driver.QueueConnect,error) {
	return &DefaultQueueConnect{
		config: config, names: map[string]int{}, datas: map[string]QueueData{},
	}, nil
}




//打开连接
func (con *DefaultQueueConnect) Open() error {
	return nil
}
//关闭连接
func (con *DefaultQueueConnect) Close() error {
	return nil
}





//注册回调
func (con *DefaultQueueConnect) Accept(handler driver.QueueHandler) error {
	con.handler = handler
	return nil
}
//订阅者注册队列
func (con *DefaultQueueConnect) Register(name string, line int) error {
	con.mutex.Lock()
	defer con.mutex.Unlock()

	//注意，这里应该处理唯一的问题，如果已经存在某name，就跳过了得
	//con.names = append(con.names, name)
	con.names[name] = line

	return nil
}

//开始消费者
func (con *DefaultQueueConnect) StartConsumer() error {
	//订阅消息
	for name,line := range con.names {
		for i:=0;i<line;i++ {
			//用局部变量，以免name在循环后被修改，可能有问题
			//事件注册时，也有同样的问题
			queueName := name
			go con.consuming(queueName)
		}
	}

	return nil
}
//实际的订阅消息方法
func (con *DefaultQueueConnect) consuming(name string) {

	cc := msg.Sub(name)
	for {
		//从管道获取消息
		v := <- cc

		//新建队列
		id := NewMd5Id()
		queue := QueueData{
			Name: name, Value: v,
		}

		//保存队列
		con.mutex.Lock()
		con.datas[id] = queue
		con.mutex.Unlock()

		//调用队列
		con.execute(id, queue.Name, queue.Value)

	}
}




//执行统一到这里
func (con *DefaultQueueConnect) execute(id string, name string, value Map) {
	req := &driver.QueueRequest{ Id: id, Name: name, Value: value }
	res := &DefaultQueueResponse{ con }
	con.handler(req, res)
}










//开始生产者
func (con *DefaultQueueConnect) StartProducer() error {
	//发布者貌似不需要干什么？
	return nil
}
//生产者发布消息
func (con *DefaultQueueConnect) Publish(name string, value Map) error {
	msg.Pub(name, value)
	return nil
}
//生产者发布延时消息
func (con *DefaultQueueConnect) DeferredPublish(name string, delay time.Duration, value Map) error {
	time.AfterFunc(delay, func() {
		msg.Pub(name, value)
	})
	return nil
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