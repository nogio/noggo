/*
	事件模块
	2016-10-22  定稿
*/

package driver

import (
	. "github.com/nogio/noggo/base"
	"time"
)

type (
	//事件驱动
	EventDriver interface {
		Connect(Map) (EventConnect,error)
	}
	//事件处理器
	EventHandler func(*EventRequest, EventResponse)

	//事件连接
	EventConnect interface {
		//打开连接
		Open() error
		//关闭连接
		Close() error



		//订阅者注册事件
		Accept(EventHandler) error
		//订阅者注册事件
		Register(name string) error
		//开始订阅者
		StartSubscriber() error


		//开始发布者
		StartPublisher() error
		//发布消息
		Publish(name string, value Map) error
		//发布延时消息
		DeferredPublish(name string, delay time.Duration, value Map) error
	}


	//事件请求实体
	EventRequest struct {
		Id string
		Name string
		Value Map
	}
	//事件响应接口
	EventResponse interface {
		//完成
		Finish(id string) error
		//重新开始
		Reevent(id string, delay time.Duration) error
	}
)
