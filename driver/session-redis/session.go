package session_redis



import (
	"github.com/nogio/noggo"
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo/depend/redigo/redis"
	"sync"
	"time"
	"errors"
	"encoding/json"
)



type (
	RedisSessionConfig struct {
		Server      string      //服务器地址，ip:端口
		Password    string      //服务器auth密码
		Expiry      int64       //默认超时时间，单位为秒
		Prefix      string      //key前缀
		Database    string      //数据库

		Idle        int         //最大空闲连接
		Active      int         //最大激活连接，同时最大并发
		Timeout     int64         //超时时间，单位为秒
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
	if config["server"]==nil || config["password"]==nil || config["prefix"]==nil {
		return nil, errors.New("配置缺少必要的信息")
	}

	//获取配置信息
	cfg := &RedisSessionConfig{
		Idle: 2, Active: 5, Timeout: 240,
	}
	if v,ok := config["server"].(string); ok {
		cfg.Server = v
	} else {
		return nil, errors.New("缺少[server]配置或不是有效的值")
	}
	if v,ok := config["password"].(string); ok {
		cfg.Password = v
	} else {
		return nil, errors.New("缺少[password]配置或不是有效的值")
	}
	if v,ok := config["prefix"].(string); ok {
		cfg.Prefix = v
	} else {
		return nil, errors.New("缺少[prefix]配置或不是有效的值")
	}

	if v,ok := config["expiry"].(float64); ok {
		cfg.Expiry = int64(v)
	}
	//数据库，redis的0-16号
	if v,ok := config["database"].(string); ok {
		cfg.Database = v
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

	return &RedisSessionConnect{
		config: cfg,
	},nil
}












//打开连接
func (rsc *RedisSessionConnect) Open() error {
	rsc.connect = &redis.Pool{
		MaxIdle: rsc.config.Idle, MaxActive: rsc.config.Active,
		IdleTimeout: time.Duration(rsc.config.Timeout) * time.Second,
		Dial: func () (redis.Conn, error) {
			c, err := redis.Dial("tcp", rsc.config.Server)
			if err != nil {
				return nil, err
			}
			//如果有验证
			if rsc.config.Password != "" {
				if _, err := c.Do("AUTH", rsc.config.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			//如果指定库
			if rsc.config.Database != "" {
				if _, err := c.Do("SELECT", rsc.config.Database); err != nil {
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
	conn := rsc.connect.Get()
	defer conn.Close()
	if err := conn.Err(); err != nil {
		return err
	} else {
		return nil
	}
}

//关闭连接
func (rsc *RedisSessionConnect) Close() error {
	if rsc.connect != nil {
		return rsc.connect.Close()
	}
	return nil
}




//查询会话，
func (rsc *RedisSessionConnect) Query(id string) (Map,error) {
	if rsc.connect == nil {
		return nil,errors.New("连接失败")
	} else {
		conn := rsc.connect.Get()
		defer conn.Close()

		key := rsc.config.Prefix + id

		val,err := redis.String(conn.Do("GET", key))
		if err != nil {
			return nil,err
		} else {

			m := Map{}

			err := json.Unmarshal([]byte(val), &m)
			if err != nil {
				return nil,err
			} else {
				//终于成功了
				return m,nil
			}
		}
	}
}



//更新会话
func (rsc *RedisSessionConnect) Update(id string, value Map, exps ...int64) error {
	if rsc.connect == nil {
		return errors.New("连接失败")
	} else {
		conn := rsc.connect.Get()
		defer conn.Close()

		//带前缀
		key := rsc.config.Prefix + id

		//JSON解析
		val,err := json.Marshal(value)
		if err != nil {
			return err
		} else {

			expiry := rsc.config.Expiry
			if len(exps) > 0 {
				expiry = exps[0]
			}

			args := []interface{}{
				key, string(val),
			}
			if expiry > 0 {
				args = append(args, "EX", expiry)
			}

			_,err := conn.Do("SET", args...)
			if err != nil {
				return err
			} else {
				//成功
				return nil
			}
		}
	}
}


//删除会话
func (rsc *RedisSessionConnect) Remove(id string) error {
	if rsc.connect == nil {
		return errors.New("连接失败")
	} else {
		conn := rsc.connect.Get()
		defer conn.Close()

		//key要加上前缀
		key := rsc.config.Prefix + id

		_,err := conn.Do("DEL", key)
		if err != nil {
			return err
		} else {
			//成功
			return nil
		}
	}
}
