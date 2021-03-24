## Storage

#### Brief

Storage package define basic methods to implement object functionality

#### Support

- fs : file system, local storage
- qs : qing yun storage, cloud object storage
- cos : tencent storage, cloud object storage

#### Method

- Save to object

Put(typ, key string, r io.Reader, contentLength int64) error

- put to object by a source object

PutByPath(key string, src string) error

- Get stream from object

FileStream(typ, key string) (io.ReadCloser, *FileInfo, error)

- Get object

Get(typ, key string, wa io.WriterAt) error

- get object to a dest object

GetToPath(key string, dest string) error

- Get object's size, last modify and mode

Stat(typ, key string) (*FileInfo, error)

- Delete object

Del(typ, key string) error

- Get the size of object

Size(typ, key string) (int64, error)

- Is object existed

IsExist(typ, key string) (bool, error)

#### Road map

- 20210324

modify usage from main and bk to only main storage, besides complete demo

- 20200220

add example