package storage

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	"unicode/utf8"
)

//当前支持三种存储模式
// fs 本地存储
// qs 青云S3接口存储
// cos 腾讯云COS共享存储

const (
	COS = "COS"
	QS  = "QINGYUN"
	FS  = "FS"
)

var (
	//once sync.Once
	instMap map[string]Storage
)

func Register(t string, s Storage) {
	instMap[t] = s
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

func PutByPath(t string, key string, path string) error {
	fd, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fd.Close()

	fi, err := fd.Stat()
	if err != nil {
		return err
	}

	return Put(t, key, fd, fi.Size())
}

func Put(t string, key string, r io.Reader, contentLength int64) error {
	inst, ok := instMap[t]
	if !ok {
		return fmt.Errorf("unregister storage type : %s", t)
	}

	return inst.Put(key, r, contentLength)
}

func GetToPath(t, key string, path string) error {
	dir, _ := filepath.Split(path)
	_ = os.MkdirAll(dir, 0666)
	fd, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer fd.Close()
	return Get(t, key, fd)
}

func Get(t, key string, wa io.WriterAt) error {
	inst, ok := instMap[t]
	if !ok {
		return fmt.Errorf("unregister storage type : %s", t)
	}

	return inst.Get(key, wa)
}

func FileStream(t, key string) (io.ReadCloser, *FileInfo, error) {
	inst, ok := instMap[t]
	if !ok {
		return nil, nil, fmt.Errorf("unregister storage type : %s", t)
	}

	return inst.FileStream(key)
}

func Size(t, key string) (int64, error) {
	inst, ok := instMap[t]
	if !ok {
		return 0, fmt.Errorf("unregister storage type : %s", t)
	}

	return inst.Size(key)
}

func IsExist(t, key string) (bool, error) {
	inst, ok := instMap[t]
	if !ok {
		return false, fmt.Errorf("unregister storage type : %s", t)
	}

	return inst.IsExist(key)
}

func Del(t, key string) error {
	inst, ok := instMap[t]
	if !ok {
		return fmt.Errorf("unregister storage type : %s", t)
	}

	return inst.Del(key)
}

// 长度在1-1023字节之间
// 第一个字符不能是 '\'
// 必须是UTF-8编码
// 不能包含 ' '、'\t'、'\r'或者'\n'等字符

func ValidKey(key string) bool {
	system := strings.ToLower(runtime.GOOS)

	if system == "ubuntu" {
		if strings.Contains(key, "（") || strings.Contains(key, "）") {
			return false
		}
	}

	if len(key) == 0 || key[0] == '\\' || key[0] == '/' || len(key) >= 1024 {
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
