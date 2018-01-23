package hashdb

import "errors"

var (
	// ErrCreatingDatabase 创建数据库错误
	ErrCreatingDatabase = errors.New(`error while creating Database`)
	// ErrDataStoreNotFound 没有找到数据记录
	ErrDataStoreNotFound = errors.New(`the requested data store was not found`)
)
