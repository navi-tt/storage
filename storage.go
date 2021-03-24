package storage

import (
	"errors"
	"fmt"
	"io"
)

//当前支持三种存储模式
// fs 本地存储
// qs 青云S3接口存储
// cos 腾讯云COS共享存储

var (
	ErrObjectNotFound            = errors.New("object not found")
	ErrObjectKeyInvalid          = errors.New("invalid object key")
	ErrObjectWritePermissionDeny = errors.New("no write permission")
	ErrObjectReadPermissionDeny  = errors.New("no read permission")
	ErrObjectEmptyContent        = errors.New("zero content")
	ErrStorageUnRegister         = errors.New("unregister storage")
)

type Storage interface {
	//把一个文件当做对象，读写即为，Get和Put

	Init(cfg string) (Storage, error)

	// 保存data至某个文件
	Put(key string, r io.Reader, contentLength int64) error

	// 根据某个文件写入到另一个文件里
	PutByPath(key string, /*要带上路径，相对路径或绝对路径都行*/src string) error

	// 获取语音流
	FileStream(key string) (io.ReadCloser, *FileInfo, error)

	// 获取数据
	Get(key string, wa io.WriterAt) error

	// 获取到某个文件
	GetToPath(key string, /*要带上路径，相对路径或绝对路径都行*/dest string) error

	// 获取文件信息  大小，修改时间，权限
	Stat(key string) (*FileInfo, error)

	// 删除文件
	Del(key string) error

	// 获取文件大小
	Size(key string) (int64, error)

	// 判断文件是否存在
	IsExist(key string) (bool, error)
}

var (
	instMap = make(map[string]Storage)
)

func Register(name string, inst Storage) {
	if inst == nil {
		panic("storage: Register storage is nil")
	}
	if _, dup := instMap[name]; !dup {
		instMap[name] = inst
	}
	instMap[name] = inst
}

func Init(name string, cfg string) (Storage, error) {

	s, ok := instMap[name]
	if !ok {
		return nil, fmt.Errorf("storage: unknown storage %q (forgotten import?)", name)
	}

	storageProvider, err := s.Init(cfg)
	if err != nil {
		return nil, err
	}

	return storageProvider, nil
}

func GetStorage(name string) (Storage, error) {
	s, ok := instMap[name]
	if !ok {
		return nil, fmt.Errorf("storage: unknown storage %q (forgotten import?)", name)
	}

	return s, nil
}
