/*
	会话接口
	2016-10-22  定稿
*/


package driver

import (
	"github.com/nogio/noggo"
	"github.com/nogio/noggo/driver/session-default"
)

//默认会话引擎
func SessionDefault() (noggo.SessionDriver) {
	return session_default.Driver()
}
