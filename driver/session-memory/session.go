package session_memory



import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo"
	"sync"
	"time"
)



type (
	//驱动
	Session struct {
		sessions map[string]SessionValue
		sessionsMutex sync.RWMutex
	}

	//值
	SessionValue struct {
		Value Map
		Expiry time.Time
	}

)


func NewSession() *Session {
	return &Session{
		map[string]SessionValue{},
	}
}






//打开会话连接
func (session *Session) Open() {
	session.sessions = map[string]SessionValue{}
}

//关闭会话连接
func (session *Session) Close() {
	//关闭不用做什么处理
}


//生成会话编号
func (session *Session) Id() string{
	return noggo.NewMd5Id()
}



//创建会话
func (session *Session) Create(id string, expiry int64) Map {
	session.sessionsMutex.Lock()
	defer session.sessionsMutex.Unlock()

	if v,ok := session.sessions[id]; ok {
		return v.Value
	} else {
		v := SessionValue{
			v, time.Now().Add(time.Second*expiry),
		}
		session.sessions[id] = v
		return v.Value
	}
}



//更新会话
func (session *Session) Update(id string, value Map, expiry int64) bool {
	session.sessionsMutex.Lock()
	defer session.sessionsMutex.Unlock()

	session.sessions[id] = SessionValue{
		value, time.Now().Add(time.Second*expiry),
	}

	return true
}


//删除会话
func (session *Session) Remove(id string) bool {
	session.sessionsMutex.Lock()
	defer session.sessionsMutex.Unlock()

	delete(session.sessions, id)

	return true
}



//回收会话
func (session *Session) Recycle(expiry int64) bool {
	session.sessionsMutex.Lock()
	defer session.sessionsMutex.Unlock()

	for k,v := range session.sessions {
		if v.Expiry < time.Now() {
			delete(session.sessions, k)
		}
	}

	return true
}