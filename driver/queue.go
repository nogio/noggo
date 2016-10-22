/*
	队列接口
	2016-10-22  定稿
*/


package driver

import (
	. "github.com/nogio/noggo/base"
	"time"
)

type (
	//队列驱动
	QueueDriver interface {
		Connect(Map) (QueueConnect,error)
	}
	//队列处理器
	QueueHandler func(*QueueRequest, QueueResponse)

	//队列连接
	QueueConnect interface {
		//打开连接
		Open() error
		//关闭连接
		Close() error


		//注册回调
		Accept(QueueHandler) error
		//注册队列
		Register(name string, line int) error
		//开始消费者
		StartConsumer() error


		//开始生产者
		StartProducer() error
		//发布消息
		Publish(name string, value Map) error
		//发布延时消息
		DeferredPublish(name string, delay time.Duration, value Map) error
	}


	//队列请求实体
	QueueRequest struct {
		Id string
		Name string
		Value Map
	}
	//队列响应接口
	QueueResponse interface {
		//完成
		Finish(id string) error
		//重新开始
		Requeue(id string, delay time.Duration) error
	}
)
