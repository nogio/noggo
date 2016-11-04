/*
	HTTP接口
	2016-10-22  定稿
*/


package driver


import (
	"github.com/nogio/noggo"
	"github.com/nogio/noggo/driver/http-default"
)




//默认HTTP引擎
func HttpDefault() (noggo.HttpDriver) {
	return http_default.Driver()
}
