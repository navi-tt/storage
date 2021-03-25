## Storage

#### 简述

Storage包通过interface定义存储的基本方法。

#### 支持

- fs: 文件系统，本地存储
- qs：青云存储，云对象存储
- cos：腾讯云存储，云对象存储

#### 接口定义

- 保存data至某个对象

Put(key string, r io.Reader, contentLength int64) error

- 从某个对象中读取数据，并存储到另一个对象

PutByPath(key string, src string) error

- 获取文件流

FileStream(key string) (io.ReadCloser, *FileInfo, error)

- 获取对象

Get(key string, wa io.WriterAt) error

- 获取对象到指定路径

GetToPath(key string, dest string) error

- 获取文件信息 大小，修改时间，权限

Stat(key string) (*FileInfo, error)

- 删除对象

Del(key string) error

- 获取对象大小

Size(key string) (int64, error)

- 判断对象是否存在

IsExist(key string) (bool, error)

#### Road map

- 20210324

修改使用逻辑，改成只用一个主要存储，及完善demo

- 20210220

添加使用范例