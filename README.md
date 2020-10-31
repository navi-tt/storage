##Storage

#### 简述

Storage包通过interface定义存储的基本方法，加上单例实现一个存储包

#### 支持
- fs: 文件系统，本地存储
- qs：青云存储，云对象存储
- cos：腾讯云存储，云对象存储

#### 接口定义

- 保存data至某个对象

Put(key string, r io.Reader, contentLength int64) error

- 获取文件流

FileStream(key string) (io.ReadCloser, *FileInfo, error)

- 获取对象

Get(key string, wa io.WriterAt) error

- 获取文件信息  大小，修改时间，权限

Stat(key string) (*FileInfo, error)

- 删除对象
	
Del(key string) error

- 获取对象大小

Size(key string) (int64, error)

- 判断对象是否存在
	
IsExist(key string) (bool, error)