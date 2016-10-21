package queue_default

import (
	. "github.com/nogio/noggo/base"
	"sync"
)
type (
	Msg struct {
		mutex   sync.Mutex
		subs    map[string]chan Map
	}
	MsgFunc func(Map)
)


//订阅消息
func (msg *Msg) Sub(name string) (chan Map) {
	msg.mutex.Lock()
	defer msg.mutex.Unlock()

	if cc,ok := msg.subs[name]; ok {
		return cc
	} else {

		cc := make(chan Map)
		msg.subs[name] = cc

		return cc
	}
}

//发布消息
func (msg *Msg) Pub(name string, value Map) error {
	msg.mutex.Lock()
	defer msg.mutex.Unlock()

	if cc,ok := msg.subs[name]; ok {
		cc <- value
	}

	return nil
}