package handle

import (
	"github.com/navi-tt/storage"
	"io"
)

var (
	storageType = storage.FS
)

func Put(key string, r io.Reader, len int64) error {
	return storage.Put(storageType, key, r, len)
}

func GetToPath(key, path string) error {
	return storage.GetToPath(storageType, key, path)
}

func Del(key string) error {
	return storage.Del(storageType, key)
}
