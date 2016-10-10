/*
	内存会话驱动
	有个BUG，Value值为引用，当一个session同时请求的时候，session.Update时会冲突
	此BUG待处理
*/


package session_default



import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo"
	"sync"
	"time"
)



type (
	//会话驱动
	DefaultDriver struct {}
	//会话连接
	DefaultConnect struct {
		config Map
		sessions map[string]noggo.SessionValue
		sessionsMutex sync.Mutex
	}
)



//返回驱动
func Driver() *DefaultDriver {
	return &DefaultDriver{}
}











//连接会话驱动
func (session *DefaultDriver) Connect(config Map) (noggo.SessionConnect) {
	return  &DefaultConnect{
		config: config, sessions: map[string]noggo.SessionValue{},
	}
}












//打开连接
func (session *DefaultConnect) Open() error {
	return nil
}

//关闭连接
func (session *DefaultConnect) Close() error {
	return nil
}




//查询会话，
func (session *DefaultConnect) Query(id string, expiry int64) (error,Map) {
	session.sessionsMutex.Lock()
	defer session.sessionsMutex.Unlock()

	if v,ok := session.sessions[id]; ok {
		return nil, v.Value
	} else {
		v := noggo.SessionValue{
			Value: Map{}, Expiry: time.Now().Add(time.Second*time.Duration(expiry)),
		}
		session.sessions[id] = v
		return nil, v.Value
	}
}



//更新会话
func (session *DefaultConnect) Update(id string, value Map, expiry int64) error {
	session.sessionsMutex.Lock()
	defer session.sessionsMutex.Unlock()

	session.sessions[id] = noggo.SessionValue{
		Value: value, Expiry: time.Now().Add(time.Second*time.Duration(expiry)),
	}

	return nil
}


//删除会话
func (session *DefaultConnect) Remove(id string) error {
	session.sessionsMutex.Lock()
	defer session.sessionsMutex.Unlock()

	delete(session.sessions, id)

	return nil
}



//回收会话
func (session *DefaultConnect) Recycle(expiry int64) error {
	session.sessionsMutex.Lock()
	defer session.sessionsMutex.Unlock()

	for k,v := range session.sessions {
		if v.Expiry.Unix() < time.Now().Unix() {
			delete(session.sessions, k)
		}
	}

	return nil
}