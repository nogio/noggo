package driver

import (
	. "github.com/nogio/noggo/base"
	"net/http"
)

type (
	HttpAcceptFunc func(res http.ResponseWriter, req *http.Request)

	//HTTP驱动
	HttpDriver interface {
		Connect(config Map) (HttpConnect)
	}
	//HTTP连接
	HttpConnect interface {
		//打开
		Open() error
		//关闭
		Close() error
		//注册
		Accept(call HttpAcceptFunc) error

		//开始
		Start(addr string) error
		//开始SSL
		StartTLS(addr string, certFile, keyFile string) error
	}
)
