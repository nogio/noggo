package noggo

import (
	. "github.com/nogio/noggo/base"
	"sync"
	"errors"
	"time"
)



// 存储驱动接口定义 begin
type (
	//存储驱动
	StorageDriver interface {
		Connect(config Map) (StorageConnect,error)
	}
	//存储连接
	StorageConnect interface {
		//打开连接
		Open() error
		//关闭连接
		Close() error
		//获取存储库对象
		Base(string) (StorageBase,error)
	}

	//存储库
	StorageBase interface {
		Close() error

		//创建，生成文件ID
		Create() (string,error)
		//删除文件
		Remove(id string) (error)
		//恢复文件，兼容设计，以支持部分存储系统逻辑删除
		Recover(id string) (error)


		//统计文件数量
		Count() (int64,error)
		//统计剩余空间，单位字节
		Space() (int64, error)


		//上传文件
		Upload(id string, path string) (error)
		//上传二进制文件
		UploadBytes(id string, bytes []byte, offset int64) (error)
		//下载文件
		Download(id string, path string) (error)
		//下载二进制文件内容，支持断点，limit不传表示到结尾，offset默认传0
		DownloadBytes(id string, offset int64, limits ...int64) ([]byte, error)


		//获取访问url
		TargetUrl(id string) (string,error)
		//获取预览（缩图）url
		PreviewUrl(id string, width, height int) (string,error)

	}
)
// 存储驱动接口定义 end









type (
	storageGlobal struct {
		mutex       sync.Mutex
		drivers     map[string]StorageDriver
		connects    map[string]StorageConnect
	}
)



//注册存储驱动
func (global *storageGlobal) Driver(name string, config StorageDriver) {
	global.mutex.Lock()
	defer global.mutex.Unlock()


	if global.drivers == nil {
		global.drivers = map[string]StorageDriver{}
	}

	if config == nil {
		panic("存储: 驱动不可为空")
	}
	global.drivers[name] = config
}


//连接驱动
func (global *storageGlobal) connect(config *storageConfig) (StorageConnect,error) {
	if storageDriver,ok := global.drivers[config.Driver]; ok {
		return storageDriver.Connect(config.Config)
	} else {
		panic("存储：不支持的驱动 " + config.Driver)
	}
}

//存储初始化
func (global *storageGlobal) init() {
	for name,config := range Config.Storage {
		con,err := global.connect(config)
		if err != nil {
			panic("存储：连接失败：" + err.Error())
		} else {
			err := con.Open()
			if err != nil {
				panic("存储：打开连接失败：" + err.Error())
			} else {
				//保存连接
				global.connects[name] = con
			}
		}
	}
}
//存储退出
func (global *storageGlobal) exit() {
	for _,con := range global.connects {
		con.Close()
	}
}







//返回存储Base对象
func (global *storageGlobal) Base(name string) (StorageBase) {
	if conn,ok := global.connects[name]; ok {
		db,err := conn.Base(name)
		if err == nil {
			return db
		}
	}
	return &noStorageBase{}
}











//---------------------------------------------------------------------------------

type (
	noStorageBase struct {}
)
func (base *noStorageBase) Close() (error) {
	return nil
}

func (base *noStorageBase) Create() (string,error) {
	return "",errors.New("无存储")
}
func (base *noStorageBase) Remove(id string) (error) {
	return errors.New("无存储")
}
func (base *noStorageBase) Recover(id string) (error) {
	return errors.New("无存储")
}



func (base *noStorageBase) Count() (int64,error) {
	return int64(0),errors.New("无存储")
}
func (base *noStorageBase) Space() (int64,error) {
	return int64(0),errors.New("无存储")
}




func (base *noStorageBase) Upload(id string, path string) (error) {
	return errors.New("无存储")
}
func (base *noStorageBase) UploadBytes(id string, bytes []byte, offset int64) (error) {
	return errors.New("无存储")
}
func (base *noStorageBase) Download(id string, path string) (error) {
	return errors.New("无存储")
}
func (base *noStorageBase) DownloadBytes(id string, offset int64, limits ...int64) ([]byte, error) {
	return nil,errors.New("无存储")
}


func (base *noStorageBase) TargetUrl(id string) (string,error) {
	return "",errors.New("无存储")
}
func (base *noStorageBase) PreviewUrl(id string, width, height int) (string,error) {
	return "",errors.New("无存储")
}