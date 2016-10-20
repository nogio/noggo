package queue_default

import (
	. "github.com/nogio/noggo/base"
	"sync"
)
type (
	Msg struct {
		mutex   sync.Mutex
		subs    map[string][]MsgFunc
	}
	MsgFunc func(Map)
)


//订阅消息
func (msg *Msg) Sub(name string, call MsgFunc) error {
	msg.mutex.Lock()
	defer msg.mutex.Unlock()

	if _,ok := msg.subs[name]; ok == false {
		msg.subs[name] = []MsgFunc{}
	}

	//加入调用列表
	msg.subs[name] = append(msg.subs[name], call)


	return nil
}

//发布消息
func (msg *Msg) Pub(name string, value Map) error {
	msg.mutex.Lock()
	defer msg.mutex.Unlock()

	if calls,ok := msg.subs[name]; ok {
		for _,call := range calls {
			go call(value)
		}
	}

	return nil
}