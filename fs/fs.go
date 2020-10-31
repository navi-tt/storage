package fs

import (
	"github.com/navi-tt/storage"
	"io"
	"os"
	"path/filepath"
)

// 本地文件系统存储服务
type FS struct {
	BaseDir string
}

func Init(cfg *FS) error {
	storage.Register(&fs{
		baseDir: cfg.BaseDir,
	})
	return nil
}

type fs struct {
	baseDir string
}

func (f *fs) open(key string) (*os.File, error) {
	if !storage.ValidKey(key) {
		return nil, storage.ErrObjectKeyInvalid
	}

	path := f.pathJoinBaseDir(key)

	fd, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, storage.ErrObjectNotFound
		}
		return nil, err
	}

	return fd, nil
}

func (f *fs) Put(key string, r io.Reader, contentLength int64) error {
	if !storage.ValidKey(key) {
		return storage.ErrObjectKeyInvalid
	}

	path := f.path(key)
	p, _ := filepath.Split(path)
	if err := os.MkdirAll(p, 0666); err != nil {
		return err
	}

	fd, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer fd.Close()

	_, err = io.Copy(fd, r)

	return err
}

func (f *fs) Get(key string, wa io.WriterAt) error {
	if !storage.ValidKey(key) {
		return storage.ErrObjectKeyInvalid
	}

	path := f.path(key)

	fd, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return storage.ErrObjectNotFound
		}
		return err
	}
	defer fd.Close()

	return storage.Copy(wa, fd)
}

func (f *fs) FileStream(key string) (io.ReadCloser, *storage.FileInfo, error) {
	if !storage.ValidKey(key) {
		return nil, nil, storage.ErrObjectKeyInvalid
	}

	path := f.path(key)

	fd, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, storage.ErrObjectNotFound
		}
		return nil, nil, err
	}

	stat, err := fd.Stat()
	if err != nil {
		return nil, nil, err
	}

	return fd, &storage.FileInfo{
		ModTime: stat.ModTime(),
		Size:    stat.Size(),
		Mode:    stat.Mode(),
	}, nil

}

func (f *fs) Stat(key string) (*storage.FileInfo, error) {
	fd, err := f.open(key)
	if err != nil {
		return nil, err
	}

	stat, err := fd.Stat()
	if err != nil {
		return nil, err
	}

	return &storage.FileInfo{
		ModTime: stat.ModTime(),
		Size:    stat.Size(),
		Mode:    stat.Mode(),
	}, nil

}

func (f *fs) Del(key string) error {
	if !storage.ValidKey(key) {
		return storage.ErrObjectKeyInvalid
	}
	path := f.path(key)
	err := os.Remove(path)
	if err != nil {
		if os.IsNotExist(err) {
			return storage.ErrObjectNotFound
		}
		return err
	}

	return nil
}

func (f *fs) Size(key string) (int64, error) {
	if !storage.ValidKey(key) {
		return 0, storage.ErrObjectKeyInvalid
	}

	path := f.path(key)
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, storage.ErrObjectNotFound
		}
		return 0, err
	}

	return info.Size(), nil
}

func (f *fs) IsExist(key string) (bool, error) {
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
	return filepath.Join(f.baseDir, key)
}

func (f *fs) pathJoinBaseDir(key string) string {
	return filepath.Join(f.baseDir, key)
}
