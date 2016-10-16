package driver

import (
	. "github.com/nogio/noggo/base"
)

type (
	//日志驱动
	LoggerDriver interface {
		Connect(config Map) (error,LoggerConnect)
	}
	//日志连接
	LoggerConnect interface {
		//打开连接
		Open() error
		//关闭连接
		Close() error

		//输出调试
		Debug(args ...interface{})
		//输出信息
		Info(args ...interface{})
		//输出错误
		Error(args ...interface{})
	}
)
