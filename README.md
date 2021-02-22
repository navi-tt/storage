## Storage

#### 简述

Storage包通过interface定义存储的基本方法。\
允许初始化多个storage，场景用于一个主存储，一个副存储用于备份。

#### 支持
- fs: 文件系统，本地存储
- qs：青云存储，云对象存储
- cos：腾讯云存储，云对象存储

#### 接口定义

- 保存data至某个对象

Put(t, key string, r io.Reader, contentLength int64) error

- 获取文件流

FileStream(t, key string) (io.ReadCloser, *FileInfo, error)

- 获取对象

Get(t, key string, wa io.WriterAt) error

- 获取文件信息  大小，修改时间，权限

Stat(t, key string) (*FileInfo, error)

- 删除对象
	
Del(t, key string) error

- 获取对象大小

Size(t, key string) (int64, error)

- 判断对象是否存在
	
IsExist(t, key string) (bool, error)

#### Road map

- 20200220

添加使用范例