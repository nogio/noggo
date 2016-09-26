package driver

import (
	. "github.com/nogio/noggo/base"
	"time"
)

type (

	//会话值
	SessionValue struct {
		Value	Map
		Expiry	time.Time
	}

	//会话驱动
	SessionDriver interface {
		Connect(config Map) (SessionConnect)
	}
	//会话连接
	SessionConnect interface {
		//打开连接
		Open()
		//关闭连接
		Close()
		//生成session唯一id方法
		Id() string
		//创建或查询会话
		Create(id string, expiry int64) Map
		//更新会话数据
		Update(id string, value Map, expiry int64) bool
		//删除会话
		Remove(id string) bool
		//回收会话，系统会每一段时间自动调用此方法
		Recycle(expiry int64) bool
	}
)
