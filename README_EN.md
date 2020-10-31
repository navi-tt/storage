##Storage

#### Brief
Storage package define basic methods to implement object functionality, besides, with an sync.once

#### Support

- fs : file system, local storage
- qs : qing yun storage, cloud object storage
- cos : tencent storage, cloud object storage

#### Method

- Save to object

Put(key string, r io.Reader, contentLength int64) error

- Get stream from object

FileStream(key string) (io.ReadCloser, *FileInfo, error)

- Get object

Get(key string, wa io.WriterAt) error

- Get object's size, last modify and mode

Stat(key string) (*FileInfo, error)

- Delete object
	
Del(key string) error

- Get the size of object

Size(key string) (int64, error)

- Is object existed
	
IsExist(key string) (bool, error)