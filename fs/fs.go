package fs

import (
	"encoding/json"
	"fmt"
	"github.com/navi-tt/storage"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func (f *fs) Init(cfg string) (storage.Storage, error) {
	fsConfig := &fs{}
	if err := json.Unmarshal([]byte(cfg), fsConfig); err != nil {
		return nil, err
	}

	fmt.Printf("[FS Init] config: \n %v \n", fsConfig)

	f.BaseDir = fsConfig.BaseDir
	return f, nil
}

type fs struct {
	BaseDir string
}

func (f *fs) PutByPath(key string, src string) error {
	fmt.Printf("[FS PUT BY PATH] object: %s \n", key)

	path := f.path(src)

	fd, fi, err := storage.OpenLocal(path)
	if err != nil {
		return err
	}
	defer fd.Close()

	return f.Put(key, fd, fi.Size())
}

func (f *fs) Put(key string, r io.Reader, contentLength int64) error {
	fmt.Printf("[FS PUT] object: %s \n", key)
	if !storage.ValidKey(key) {
		return storage.ErrObjectKeyInvalid
	}

	path := f.path(key)
	p, _ := filepath.Split(path)
	if err := os.MkdirAll(p, 0766); err != nil {
		fmt.Println(p)
		return err
	}

	fd, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0766)
	if err != nil {
		if os.IsPermission(err) {
			return storage.ErrObjectWritePermissionDeny
		}
		return err
	}
	defer fd.Close()

	_, err = io.Copy(fd, r)

	return err
}

func (f *fs) Get(key string, wa io.WriterAt) error {
	fmt.Printf("[FS GET] object: %s \n", key)
	if !storage.ValidKey(key) {
		return storage.ErrObjectKeyInvalid
	}

	path := f.path(key)

	fd, _, err := storage.OpenLocal(path)
	if err != nil {
		return err
	}
	defer fd.Close()

	return storage.Copy(wa, fd)
}

func (f *fs) GetToPath(key string, dest string) error {
	fmt.Printf("[FS GET TO PATH] object: %s \n", key)

	dir, _ := filepath.Split(dest)
	_ = os.MkdirAll(dir, 0766)
	fd, err := os.OpenFile(dest, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0766)
	if err != nil {
		return err
	}
	defer fd.Close()
	return f.Get(key, fd)
}

func (f *fs) FileStream(key string) (io.ReadCloser, *storage.FileInfo, error) {
	fmt.Printf("[FS FileStream] object: %s \n", key)
	if !storage.ValidKey(key) {
		return nil, nil, storage.ErrObjectKeyInvalid
	}

	path := f.path(key)

	fd, stat, err := storage.OpenLocal(path)
	if err != nil {
		return nil, nil, err
	}
	defer fd.Close()

	return fd, &storage.FileInfo{
		ModTime: stat.ModTime(),
		Size:    stat.Size(),
		Mode:    stat.Mode(),
	}, nil

}

func (f *fs) Stat(key string) (*storage.FileInfo, error) {
	fmt.Printf("[FS STAT] object: %s \n", key)

	path := f.path(key)

	fd, stat, err := storage.OpenLocal(path)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	return &storage.FileInfo{
		ModTime: stat.ModTime(),
		Size:    stat.Size(),
		Mode:    stat.Mode(),
	}, nil

}

func (f *fs) Del(key string) error {
	fmt.Printf("[FS DEL] object: %s \n", key)
	if !storage.ValidKey(key) {
		return storage.ErrObjectKeyInvalid
	}

	path := f.path(key)
	err := os.Remove(path)
	if err != nil {
		if os.IsNotExist(err) {
			return storage.ErrObjectNotFound
		}

		if os.IsPermission(err) {
			return storage.ErrObjectWritePermissionDeny
		}
		return err
	}

	return nil
}

func (f *fs) Size(key string) (int64, error) {
	fmt.Printf("[FS SIZE] object: %s \n", key)
	if !storage.ValidKey(key) {
		return 0, storage.ErrObjectKeyInvalid
	}

	path := f.path(key)
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, storage.ErrObjectNotFound
		}

		if os.IsPermission(err) {
			return 0, storage.ErrObjectReadPermissionDeny
		}
		return 0, err
	}

	return info.Size(), nil
}

func (f *fs) IsExist(key string) (bool, error) {
	fmt.Printf("[FS ISEXIST] object: %s \n", key)
	if !storage.ValidKey(key) {
		return false, storage.ErrObjectKeyInvalid
	}

	path := f.path(key)
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (f *fs) path(key string) string {
	if strings.TrimSpace(f.BaseDir) == "" {
		return key
	}

	return filepath.Join(f.BaseDir, key)
}

var f = &fs{}

func init() {
	storage.Register(storage.FS, f)
}
