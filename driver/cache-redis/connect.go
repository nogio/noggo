package cache_redis

import (
	"github.com/nogio/noggo"
	"github.com/nogio/noggo/depend/redigo/redis"
	"time"
)

type (

	//数据库连接
	RedisCacheConnect struct {
		config  *RedisCacheConfig
		connect *redis.Pool
	}
)

//打开连接
func (rcc *RedisCacheConnect) Open() error {
	rcc.connect = &redis.Pool{
		MaxIdle: rcc.config.Idle, MaxActive: rcc.config.Active,
		IdleTimeout: time.Duration(rcc.config.Timeout) * time.Second,
		Dial: func () (redis.Conn, error) {
			c, err := redis.Dial("tcp", rcc.config.Server)
			if err != nil {
				return nil, err
			}
			//如果有验证
			if rcc.config.Password != "" {
				if _, err := c.Do("AUTH", rcc.config.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			//如果指定库
			if rcc.config.Database != "" {
				if _, err := c.Do("SELECT", rcc.config.Database); err != nil {
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
	conn := rcc.connect.Get()
	defer conn.Close()
	if err := conn.Err(); err != nil {
		return err
	} else {
		return nil
	}
}
//关闭连接
func (rcc *RedisCacheConnect) Close() error {
	if rcc.connect != nil {
		return rcc.connect.Close()
	}
	return nil
}




func (rcc *RedisCacheConnect) Base(name string) (noggo.CacheBase,error) {
	return &RedisCacheBase{name, rcc},nil
}