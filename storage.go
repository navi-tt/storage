package storage

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
	"unicode/utf8"
)

//当前支持三种存储模式
// fs 本地存储
// qs 青云S3接口存储
// cos 腾讯云COS共享存储

var (
	once sync.Once
	inst Storage
)

func Register(s Storage) {
	once.Do(func() {
		inst = s
	})
}

var (
	ErrObjectNotFound   = errors.New("object not found")
	ErrObjectKeyInvalid = errors.New("invalid object key")
)

type Storage interface {
	//把一个文件当做对象，读写即为，Get和Put

	// 保存data至某个文件
	Put(key string, r io.Reader, contentLength int64) error

	// 获取语音流
	FileStream(key string) (io.ReadCloser, *FileInfo, error)

	// 获取数据
	Get(key string, wa io.WriterAt) error

	// 获取文件信息  大小，修改时间，权限
	Stat(key string) (*FileInfo, error)

	// 删除文件
	Del(key string) error

	// 获取文件大小
	Size(key string) (int64, error)

	// 判断文件是否存在
	IsExist(key string) (bool, error)
}

func PutByPath(key string, path string) error {
	fd, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fd.Close()

	fi, err := fd.Stat()
	if err != nil {
		return err
	}

	return Put(key, fd, fi.Size())
}

func Put(key string, r io.Reader, contentLength int64) error {
	return inst.Put(key, r, contentLength)
}

func GetToPath(key string, path string) error {
	dir, _ := filepath.Split(path)
	_ = os.MkdirAll(dir, 0666)
	fd, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer fd.Close()
	return Get(key, fd)
}

func Get(key string, wa io.WriterAt) error {
	return inst.Get(key, wa)
}

func FileStream(key string) (io.ReadCloser, *FileInfo, error) {
	return inst.FileStream(key)
}

func Size(key string) (int64, error) {
	return inst.Size(key)
}

func IsExist(key string) (bool, error) {
	return inst.IsExist(key)
}

func Del(key string) error {
	return inst.Del(key)
}

// 长度在1-1023字节之间
// 第一个字符不能是 '\'
// 必须是UTF-8编码
// 不能包含 ' '、'\t'、'\r'或者'\n'等字符

func ValidKey(key string) bool {
	if len(key) == 0 || key[0] == '\\' || len(key) >= 1024 {
		return false
	}

	for i := range key {
		if key[i] == ' ' || key[i] == '\t' || key[i] == '\r' || key[i] == '\n' {
			return false
		}
	}

	return utf8.ValidString(key)
}

type FileInfo struct {
	Size    int64
	ModTime time.Time
	Mode    os.FileMode
}
