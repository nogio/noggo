package session_memory



import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo"
	"sync"
	"time"
)



type (
	//会话驱动
	MemoryDriver struct {}
	//会话连接
	MemoryConnect struct {
		config Map
		sessions map[string]noggo.SessionValue
		sessionsMutex sync.Mutex
	}
)



//返回驱动
func Driver() *MemoryDriver {
	return &MemoryDriver{}
}











//连接会话驱动
func (session *MemoryDriver) Connect(config Map) (noggo.SessionConnect) {
	return  &MemoryConnect{
		config: config, sessions: map[string]noggo.SessionValue{},
	}
}












//打开连接
func (session *MemoryConnect) Open() error {
	return nil
}

//关闭连接
func (session *MemoryConnect) Close() {

}






//生成ID
func (session *MemoryConnect) Id() string {
	return noggo.NewMd5Id()
}


//创建会话
func (session *MemoryConnect) Create(id string, expiry int64) Map {
	session.sessionsMutex.Lock()
	defer session.sessionsMutex.Unlock()

	if v,ok := session.sessions[id]; ok {
		return v.Value
	} else {
		v := noggo.SessionValue{
			Value: Map{}, Expiry: time.Now().Add(time.Second*time.Duration(expiry)),
		}
		session.sessions[id] = v
		return v.Value
	}
}



//更新会话
func (session *MemoryConnect) Update(id string, value Map, expiry int64) bool {
	session.sessionsMutex.Lock()
	defer session.sessionsMutex.Unlock()

	session.sessions[id] = noggo.SessionValue{
		Value: value, Expiry: time.Now().Add(time.Second*time.Duration(expiry)),
	}

	return true
}


//删除会话
func (session *MemoryConnect) Remove(id string) bool {
	session.sessionsMutex.Lock()
	defer session.sessionsMutex.Unlock()

	delete(session.sessions, id)

	return true
}



//回收会话
func (session *MemoryConnect) Recycle(expiry int64) bool {
	session.sessionsMutex.Lock()
	defer session.sessionsMutex.Unlock()

	for k,v := range session.sessions {
		if v.Expiry.Unix() < time.Now().Unix() {
			delete(session.sessions, k)
		}
	}

	return true
}