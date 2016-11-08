package cache_redis


import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo"
	"errors"
)

type (
	RedisCacheConfig struct {
		Server      string      //服务器地址，ip:端口
		Password    string      //服务器auth密码
		Database    string      //数据库
		Expiry      int64       //默认超时时间，单位为秒，0不过期
		Prefix      string      //key前缀

		Idle        int         //最大空闲连接
		Active      int         //最大激活连接，同时最大并发
		Timeout     int64       //连接超时时间，单位为秒
	}

	RedisCacheDriver struct {}
)

//返回驱动
func Driver() (noggo.CacheDriver) {
	return &RedisCacheDriver{}
}






//驱动连接
func (drv *RedisCacheDriver) Connect(config Map) (noggo.CacheConnect,error) {
	//检查config
	if config == nil {
		return nil, errors.New("配置不可为空")
	}


	//获取配置信息
	cfg := &RedisCacheConfig{
		Database: "0", Expiry: 60*60*24,
		Idle: 2, Active: 5, Timeout: 240,
	}

	//服务器
	if v,ok := config["server"].(string); ok {
		cfg.Server = v
	} else {
		return nil, errors.New("缺少[server]配置或不是有效的值")
	}
	//auth密码
	if v,ok := config["password"].(string); ok {
		cfg.Password = v
	} else {
		return nil, errors.New("缺少[password]配置或不是有效的值")
	}

	//数据库，redis的0-16号
	if v,ok := config["database"].(string); ok {
		cfg.Database = v
	}
	//缓存key前缀
	if v,ok := config["prefix"].(string); ok {
		cfg.Prefix = v
	}
	//默认过期时间
	if v,ok := config["expiry"].(float64); ok {
		cfg.Expiry = int64(v)
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

	return &RedisCacheConnect{
		config: cfg,
	},nil
}
