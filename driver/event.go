/*
	事件模块
	2016-10-22  定稿
*/

package driver

import (
	"github.com/nogio/noggo"
	"github.com/nogio/noggo/driver/event-default"
)





//默认事件引擎
func EventDefault() (noggo.EventDriver) {
	return event_default.Driver()
}
