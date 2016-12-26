package storage_default


import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo"
	"errors"
)

type (
	DefaultStorageDriver struct {
	}
)

//返回驱动
func Driver() (noggo.StorageDriver) {
	return &DefaultStorageDriver{}
}



//驱动连接
func (drv *DefaultStorageDriver) Connect(config Map) (noggo.StorageConnect,error) {

	if config == nil {
		return nil,errors.New("无效的存储配置")
	}
	path := ""
	if vv,ok := config["path"].(string); ok {
		path == vv
	} else {
		return nil,errors.New("无效的存储[path]配置")
	}


	return &DefaultStorageConnect{
		config: config, path: path,
	},nil
}
