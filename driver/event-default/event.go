package event_default

/*
	默认事件驱动
	直接基于内存的事件，这样只在单进程有效，不支持跨进程事件
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
	DefaultEventDriver struct {
	}
	//连接
	DefaultEventConnect struct {
		config      Map
		handler     driver.EventHandler

		mutex       sync.Mutex
		names       []string
		datas      map[string]EventData
	}

	//响应对象
	DefaultEventResponse struct {
		con *DefaultEventConnect
	}


	EventData struct{
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
func Driver() (driver.EventDriver) {
	return &DefaultEventDriver{}
}




//连接
func (drv *DefaultEventDriver) Connect(config Map) (driver.EventConnect,error) {
	return &DefaultEventConnect{
		config: config, names: []string{}, datas: map[string]EventData{},
	}, nil
}




//打开连接
func (drv *DefaultEventConnect) Open() error {
	return nil
}
//关闭连接
func (drv *DefaultEventConnect) Close() error {
	return nil
}






//订阅者注册回调
func (con *DefaultEventConnect) Accept(handler driver.EventHandler) error {
	con.mutex.Lock()
	defer con.mutex.Unlock()

	//保存回调
	con.handler = handler

	return nil
}
//订阅者注册事件
func (con *DefaultEventConnect) Register(name string) error {
	con.mutex.Lock()
	defer con.mutex.Unlock()

	//注意，这里应该处理唯一的问题，如果已经存在某name，就跳过了得
	con.names = append(con.names, name)

	return nil
}

//开始订阅者
func (con *DefaultEventConnect) StartSubscriber() error {
	//订阅消息
	for _,name := range con.names {

		//一定要用一个局部变量
		//因为name循环的， 在加调的时候name已经改变
		//回调时候就会有问题
		eventName := name

		msg.Sub(name, func(value Map) {

			//新建计划
			id := NewMd5Id()
			event := EventData{
				Name: eventName, Value: value,
			}

			//保存计划
			con.mutex.Lock()
			con.datas[id] = event
			con.mutex.Unlock()

			//调用事件
			con.execute(id, event.Name, event.Value)
		})
	}
	return nil
}


//执行统一到这里
func (con *DefaultEventConnect) execute(id string, name string, value Map) {
	req := &driver.EventRequest{ Id: id, Name: name, Value: value }
	res := &DefaultEventResponse{ con }
	con.handler(req, res)
}











//开始发布者
func (con *DefaultEventConnect) StartPublisher() error {
	//内存版发布者貌似不需要干什么
	return nil
}
//发布者发布消息
func (con *DefaultEventConnect) Publish(name string, value Map) error {
	msg.Pub(name, value)
	return nil
}
//发布者延时发布消息
func (con *DefaultEventConnect) DeferredPublish(name string, delay time.Duration, value Map) error {
	time.AfterFunc(delay, func() {
		msg.Pub(name, value)
	})
	return nil
}























//完成事件，从列表中清理
func (con *DefaultEventConnect) finish(id string) error {
	con.mutex.Lock()
	defer con.mutex.Unlock()

	delete(con.datas, id)
	return nil
}
//重开事件
func (con *DefaultEventConnect) reevent(id string, delay time.Duration) error {
	if event,ok := con.datas[id]; ok {
		time.AfterFunc(delay, func() {
			con.execute(id, event.Name, event.Value)
		})
	}

	return errors.New("计划不存在")
}











//响应接口对象



//完成事件，从列表中清理
func (res *DefaultEventResponse) Finish(id string) error {
	return res.con.finish(id)
}
//重开事件
func (res *DefaultEventResponse) Reevent(id string, delay time.Duration) error {
	return res.con.reevent(id, delay)
}