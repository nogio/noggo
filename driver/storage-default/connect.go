package storage_default

import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo"
	"sync"
)

type (
	//数据库连接
	DefaultStorageConnect struct {
		config  Map
		path    string
	}
)

//打开连接
func (conn *DefaultStorageConnect) Open() error {
	return nil
}
//关闭连接
func (conn *DefaultStorageConnect) Close() error {
	return nil
}




func (conn *DefaultStorageConnect) Base(name string) (noggo.StorageBase,error) {
	return &DefaultStorageBase{name, conn},nil
}