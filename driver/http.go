/*
	HTTP接口
	2016-10-22  未定稿
*/


package driver

import (
	. "github.com/nogio/noggo/base"
	"net/http"
)

type (

	//HTTP驱动
	HttpDriver interface {
		Connect(config Map) (HttpConnect,error)
	}

	HttpAccept func(req *http.Request, res http.ResponseWriter)

	//HTTP连接
	HttpConnect interface {
		//打开驱动连接
		Open() error
		//关闭驱动连接
		Close() error

		//注册回调
		Accept(call HttpAccept) error

		//开始
		Start(addr string) error
		//开始TLS
		StartTLS(addr string, certFile, keyFile string) error

	}
)
