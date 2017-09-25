package event_redis


import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo"
	"errors"
	"time"
	"sync"
	"github.com/nogio/noggo/depend/redigo/redis"
	"encoding/json"
	"strings"
)


type (
	//配置文件
	redisEventConfig struct {
		Server      string      //服务器地址，ip:端口
		Password    string      //服务器auth密码
		Database    string      //数据库
		Prefix      string      //key前缀

		Idle        int         //最大空闲连接
		Active      int         //最大激活连接，同时最大并发
		Timeout     int64       //连接超时时间，单位为秒
	}

	//驱动
	redisEventDriver struct {}
	//连接
	redisEventConnect struct {
		config  redisEventConfig
		handler noggo.EventHandler

		mutex   sync.Mutex
		names   []string
		//data是内存版的事件，用来保存所有调用事件列表的
		//为了可以reevent而做，更新了response接口，不需要这个了
		//datas   map[string]redisEventData


		connect *redis.Pool
	}

	//响应对象
	redisEventResponse struct {
		bonder *redisEventConnect
	}

	redisEventData struct{
		Name    string
		//Time    string
		Value   Map
	}

)




//返回驱动
func Driver() (noggo.EventDriver) {
	return &redisEventDriver{}
}




//连接
func (drv *redisEventDriver) Connect(config Map) (noggo.EventConnect,error) {

	//检查config
	if config == nil {
		return nil, errors.New("配置不可为空")
	}


	//获取配置信息
	cfg := redisEventConfig{
		Database: "0", Idle: 2, Active: 5, Timeout: 240,
	}

	//服务器
	if v,ok := config["server"].(string); ok {
		cfg.Server = v
	} else {
		return nil, errors.New("缺少[server]配置或不是有效的值")
	}
	//auth密码，密码不是必填的
	if v,ok := config["password"].(string); ok {
		cfg.Password = v
	}

	//数据库，redis的0-16号
	if v,ok := config["database"].(string); ok {
		cfg.Database = v
	}
	//缓存key前缀
	if v,ok := config["prefix"].(string); ok {
		cfg.Prefix = v
	}


	if v,ok := config["idle"].(float64); ok {
		cfg.Idle = int(v)
	}
	if v,ok := config["active"].(float64); ok {
		cfg.Active = int(v)
	}
	if v,ok := config["timeout"].(float64); ok {
		cfg.Timeout = int64(v)
	}

	return &redisEventConnect{
		config: cfg,
	},nil
}













//打开连接
func (bonder *redisEventConnect) Open() error {

	bonder.connect = &redis.Pool{
		MaxIdle: bonder.config.Idle, MaxActive: bonder.config.Active,
		IdleTimeout: time.Duration(bonder.config.Timeout) * time.Second,
		Dial: func () (redis.Conn, error) {
			c, err := redis.Dial("tcp", bonder.config.Server)
			if err != nil {
				return nil, err
			}
			//如果有验证
			if bonder.config.Password != "" {
				if _, err := c.Do("AUTH", bonder.config.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			//如果指定库
			if bonder.config.Database != "" {
				if _, err := c.Do("SELECT", bonder.config.Database); err != nil {
					c.Close()
					return nil, err
				}
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}

	//打开一个试一下
	conn := bonder.connect.Get(); defer conn.Close()
	return conn.Err()
}
//关闭连接
func (bonder *redisEventConnect) Close() error {
	if bonder.connect != nil {
		return bonder.connect.Close()
	}
	return nil
}






//订阅者注册回调
func (bonder *redisEventConnect) Accept(handler noggo.EventHandler) error {
	bonder.mutex.Lock()
	defer bonder.mutex.Unlock()

	//保存回调
	bonder.handler = handler

	return nil
}


//订阅者注册事件
func (bonder *redisEventConnect) Register(name string) error {
	bonder.mutex.Lock()
	defer bonder.mutex.Unlock()

	//注意，这里应该处理唯一的问题，如果已经存在某name，就跳过了得
	bonder.names = append(bonder.names, name)

	return nil
}





//开始订阅者
func (bonder *redisEventConnect) StartSubscriber() error {

	//订阅消息
	/*
	for _,name := range bonder.names {
		//一定要用一个局部变量
		//因为name循环的， 在加调的时候name已经改变
		//回调时候就会有问题
		eventName := name

		conn,err := bonder.connect.Dial()
		if err != nil {
			return err
		}

		go bonder.subscribing(conn, eventName)
	}
	*/

	//改成一个连接监听所有
	go bonder.subscribing()
	return nil
}




//这里是调用的，要一直循环啊，连接一直不关的样子
func (bonder *redisEventConnect) subscribing() {
	names := []interface{}{}
	for _,name := range bonder.names {
		names = append(names, bonder.config.Prefix+name)
	}


	//多个同时调用会串事件。。。待处理

	conn,err := bonder.connect.Dial()
	if err == nil {
		defer conn.Close()

		psc := redis.PubSubConn{Conn: conn}
		psc.Subscribe(names...) //一次订阅多个

		for {
			switch msg := psc.Receive().(type) {
			case redis.Message:
				//fmt.Printf("%s: message: %s\n", v.Channel, v.Data)
				go bonder.gotmsg(msg)
			case redis.Subscription:
				//fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
			case error:
				//3秒后重新连接
				//time.Sleep(time.Second * 3)
				//goto reconnection
			}
		}
	}
}



//执行统一到这里
func (bonder *redisEventConnect) gotmsg(msg redis.Message) {
	id := NewMd5Id()
	//name := strings.TrimLeft(msg.Channel, bonder.config.Prefix)
	name := strings.Replace(msg.Channel, bonder.config.Prefix, "", 1)
	value := Map{}
	json.Unmarshal(msg.Data, &value)
	bonder.execute(id, name, value)
}






//执行统一到这里
func (bonder *redisEventConnect) execute(id string, name string, value Map) {
	req := &noggo.EventRequest{ Id: id, Name: name, Value: value }
	res := &redisEventResponse{ bonder }
	bonder.handler(req, res)
}











//开始发布者
func (bonder *redisEventConnect) StartPublisher() error {
	//好像不用干什么，发布消息的时候，直接把从池里拉一条连接比较好
	//或者可以保存一条连接专门用来发消息（但可能有并发问题，还是用上一行的方案）
	return nil
}
//发布者发布消息
func (bonder *redisEventConnect) Publish(name string, value Map) error {

	if bonder.connect == nil {
		return errors.New("无效失败")
	}

	//再转成json
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	//获取连接
	conn, err := bonder.connect.Dial()
	if err != nil {
		return err
	}
	defer conn.Close()

	//写入
	realname := bonder.config.Prefix + name
	_,err = conn.Do("PUBLISH", realname, string(bytes))
	return err
}
//发布者延时发布消息
//注意：redis并不支持延迟发消息
//所以，在内存延迟是不可靠的，当出现意外情况，此消息可能丢失
func (bonder *redisEventConnect) DeferredPublish(name string, delay time.Duration, value Map) error {
	time.AfterFunc(delay, func() {
		bonder.Publish(name, value)
	})
	return nil
}





























//完成事件，从列表中清理
func (bonder *redisEventConnect) finish(id string) error {
	//redis版事件，不需要做任何处理
	return nil
}
//重开事件
//注意：redis并不支持延迟消息，此处在进程中延迟是不可靠的
func (bonder *redisEventConnect) reevent(id string, name string, value Map, delay time.Duration) error {
	time.AfterFunc(delay, func() {
		bonder.execute(id, name, value)
	})
	return nil
}











//响应接口对象



//完成事件，从列表中清理
func (res *redisEventResponse) Finish(id string) error {
	return res.bonder.finish(id)
}
//重开事件
func (res *redisEventResponse) Reevent(id string, name string, value Map, delay time.Duration) error {
	return res.bonder.reevent(id, name, value, delay)
}