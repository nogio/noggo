/*
	内存会话驱动
	有个BUG，Value值为引用，当一个session同时请求的时候，session.Update时会冲突
	此BUG待处理
*/


package session_default



import (
	"github.com/nogio/noggo"
	. "github.com/nogio/noggo/base"
	"sync"
	"time"
	"errors"
	"github.com/nogio/noggo/depend/redigo/redis"
)



type (
	RedisSessionConfig struct {
		Server      string      //服务器地址，ip:端口
		Password    string      //服务器auth密码
		Expiry      int64       //默认超时时间，单位为秒
		Prefix      string      //key前缀

		Idle        int         //最大空闲连接
		Active      int         //最大激活连接，同时最大并发
		Timeout     int         //超时时间，单位为秒
	}
	//会话驱动
	RedisSessionDriver struct {}
	//会话连接
	RedisSessionConnect struct {
		config *RedisSessionConfig
		connect *redis.Pool
		mutex sync.Mutex
	}
)



//返回驱动
func Driver() (noggo.SessionDriver) {
	return &RedisSessionDriver{}
}











//连接
func (driver *RedisSessionDriver) Connect(config Map) (noggo.SessionConnect,error) {

	//检查config
	if config == nil {
		return nil, errors.New("配置不可为空")
	}
	if config["server"]==nil || config["password"] || config["expiry"]==nil || config["prefix"]==nil {
		return nil, errors.New("配置缺少必要的信息")
	}

	//获取配置信息
	cfg := &RedisSessionConfig{
		Idle: 2, Active: 5, Timeout: 240,
	}
	if v,ok := config["string"].(string); ok {
		cfg.Server = v
	} else {
		return nil, errors.New("缺少[server]配置或不是有效的值")
	}
	if v,ok := config["password"].(string); ok {
		cfg.Password = v
	} else {
		return nil, errors.New("缺少[password]配置或不是有效的值")
	}
	if v,ok := config["expiry"].(float64); ok {
		cfg.Expiry = int64(v)
	} else {
		return nil, errors.New("缺少[expiry]配置或不是有效的值")
	}
	if v,ok := config["prefix"].(string); ok {
		cfg.Prefix = v
	} else {
		return nil, errors.New("缺少[prefix]配置或不是有效的值")
	}


	if v,ok := config["idle"].(float64); ok {
		cfg.Idle = int(v)
	}
	if v,ok := config["active"].(float64); ok {
		cfg.Active = int(v)
	}
	if v,ok := config["timeout"].(float64); ok {
		cfg.Timeout = int(v)
	}

	return &RedisSessionConnect{
		config: cfg,
	},nil
}












//打开连接
func (connect *RedisSessionConnect) Open() error {
	connect.connect = redis.Pool{
		MaxIdle: connect.config.Idle, MaxActive: connect.config.Active,
		IdleTimeout: connect.config.Timeout * time.Second,
		Dial: func () (redis.Conn, error) {
			c, err := redis.Dial("tcp", connect.config.Server)
			if err != nil {
				return nil, err
			}
			if connect.config.Password != "" {
				if _, err := c.Do("AUTH", connect.config.Password); err != nil {
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
	conn := connect.connect.Get()
	defer conn.Close()
	if err := conn.Err(); err != nil {
		return err
	} else {
		return nil
	}
}

//关闭连接
func (connect *RedisSessionConnect) Close() error {
	return nil
}




//查询会话，
func (session *RedisSessionConnect) Query(id string, expiry int64) (Map,error) {
	return nil,nil
}



//更新会话
func (session *RedisSessionConnect) Update(id string, value Map, expiry int64) error {
	return nil
}


//删除会话
func (session *RedisSessionConnect) Remove(id string) error {
	return nil
}
