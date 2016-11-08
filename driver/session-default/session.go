/*
	内存会话驱动
	有个BUG，Value值为引用，当一个session同时请求的时候，session.Update时会冲突
	此BUG待处理
*/


package session_default



import (
	. "github.com/nogio/noggo/base"
	"sync"
	"time"
	"github.com/nogio/noggo"
	"errors"
)



type (
	//会话驱动
	DefaultDriver struct {}
	//会话连接
	DefaultConnect struct {
		config Map
		mutex sync.Mutex
		sessions map[string]DefaultSessionValue

	}

	DefaultSessionValue struct {
		Value	Map
		Expiry	time.Time
	}
)



//返回驱动
func Driver() (noggo.SessionDriver) {
	return &DefaultDriver{}
}











//连接
func (session *DefaultDriver) Connect(config Map) (noggo.SessionConnect,error) {
	return &DefaultConnect{
		config: config, sessions: map[string]DefaultSessionValue{},
	},nil
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
func (session *DefaultConnect) Query(id string) (Map,error) {
	session.mutex.Lock()
	defer session.mutex.Unlock()

	if v,ok := session.sessions[id]; ok {
		//简单复制一份
		m := Map{}
		for k,v := range v.Value {
			m[k] = v
		}
		return m,nil
	} else {
		return nil,errors.New("无会话")
	}

/*else {
		v := DefaultSessionValue{
			Value: Map{}, Expiry: time.Now().Add(time.Second*time.Duration(expiry)),
		}
		session.sessions[id] = v
		return v.Value,nil
	}
	*/
}



//更新会话
func (session *DefaultConnect) Update(id string, value Map, exps ...int64) error {
	session.mutex.Lock()
	defer session.mutex.Unlock()

	expiry := int64(3600)
	if len(exps) > 0 {
		expiry = exps[0]
	}

	session.sessions[id] = DefaultSessionValue{
		Value: value, Expiry: time.Now().Add(time.Second*time.Duration(expiry)),
	}

	return nil
}


//删除会话
func (session *DefaultConnect) Remove(id string) error {
	session.mutex.Lock()
	defer session.mutex.Unlock()

	delete(session.sessions, id)

	return nil
}



//回收会话
//自动回收过期的会话
func (session *DefaultConnect) Recycle(expiry int64) error {
	session.mutex.Lock()
	defer session.mutex.Unlock()

	for k,v := range session.sessions {
		if v.Expiry.Unix() < time.Now().Unix() {
			delete(session.sessions, k)
		}
	}

	return nil
}