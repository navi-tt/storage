package storage

import (
	"io"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
	"unicode/utf8"
)

type StorageType string

const (
	FS  StorageType = "fs"
	COS StorageType = "cos"
	QS  StorageType = "qs"
)

// 长度在1-1023字节之间
// 第一个字符不能是 '\'
// 必须是UTF-8编码
// 不能包含 ' '、'\t'、'\r'或者'\n'等字符

func ValidKey(key string) bool {
	system := strings.ToLower(runtime.GOOS)

	_, fn := path.Split(key)

	if system == "ubuntu" {
		if strings.Contains(fn, "（") || strings.Contains(fn, "）") {
			return false
		}
	}

	if len(fn) == 0 || fn[0] == '\\' || fn[0] == '/' || len(fn) >= 1024 {
		return false
	}

	for i := range fn {
		if fn[i] == ' ' || fn[i] == '\t' || fn[i] == '\r' || fn[i] == '\n' {
			return false
		}
	}

	return utf8.ValidString(fn)
}

type FileInfo struct {
	Size    int64
	ModTime time.Time
	Mode    os.FileMode
}

func OpenLocal(key string) (*os.File, os.FileInfo, error) {
	var (
		err  error
		fd   *os.File
		stat os.FileInfo
	)

	if !ValidKey(key) {
		return nil, nil, ErrObjectKeyInvalid
	}

	fd, err = os.Open(key)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, ErrObjectNotFound
		}

		if os.IsPermission(err) {
			return nil, nil, ErrObjectReadPermissionDeny
		}

		return nil, nil, err
	}

	stat, err = fd.Stat()
	if err != nil {
		return nil, nil, err
	}

	if stat.Size() == 0 {
		return nil, nil, ErrObjectEmptyContent
	}

	return fd, stat, nil
}

func Copy(wa io.WriterAt, r io.Reader) error {
	var (
		buf       = make([]byte, 1024)
		off int64 = 0
	)

	for {
		n, err1 := r.Read(buf)
		buf := buf[:n]

		n, err2 := wa.WriteAt(buf, off)
		if err2 != nil {
			return err2
		}

		off += int64(n)
		if err1 != nil {
			if err1 == io.EOF {
				break
			}
			return err1
		}
	}
	return nil
}
