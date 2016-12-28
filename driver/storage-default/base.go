package storage_default


import (
	. "github.com/nogio/noggo/base"
	"errors"
)

type (
	DefaultStorageBase struct {
		name    string
		conn    *DefaultStorageConnect
	}
)
//关闭库
func (base *DefaultStorageBase) Close() (error) {
	return nil
}

func (base *DefaultStorageBase) Create(args ...Map) (string,error) {
	return "",errors.New("未实现")
}
func (base *DefaultStorageBase) Remove(id string) (error) {
	return errors.New("未实现")
}
func (base *DefaultStorageBase) Recover(id string) (error) {
	return errors.New("未实现")
}

func (base *DefaultStorageBase) Count() (int64,error) {
	return int64(0),errors.New("未实现")
}
func (base *DefaultStorageBase) Space() (int64,error) {
	return int64(0),errors.New("未实现")
}



func (base *DefaultStorageBase) Upload(id string, path string) (error) {
	return errors.New("未实现")
}
func (base *DefaultStorageBase) UploadBytes(id string, bytes []byte, offset int64) (error) {
	return errors.New("未实现")
}
func (base *DefaultStorageBase) Download(id string, path string) (error) {
	return errors.New("未实现")
}
func (base *DefaultStorageBase) DownloadBytes(id string, offset int64, limits ...int64) ([]byte, error) {
	return nil,errors.New("未实现")
}


func (base *DefaultStorageBase) PublicUrl(id string) (string,error) {
	return "",errors.New("未实现")
}
func (base *DefaultStorageBase) PrivateUrl(id string) (string,error) {
	return "",errors.New("未实现")
}
func (base *DefaultStorageBase) PreviewUrl(id string, width int, height int) (string,error) {
	return "",errors.New("未实现")
}
