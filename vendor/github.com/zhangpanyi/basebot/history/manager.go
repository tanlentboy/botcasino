package history

import (
	"errors"

	"github.com/zhangpanyi/basebot/hashdb"
)

// Manager 记录管理器
type Manager struct {
	db *hashdb.HashDatabase // 数据库
}

// NewManager 创建记录管理器
func NewManager(bucketNum uint32) (*Manager, error) {
	db, err := hashdb.Create(bucketNum)
	if err != nil {
		return nil, err
	}
	return &Manager{db: db}, nil
}

// Del 删除记录
func (m *Manager) Del(userID uint32) {
	m.db.Erase(userID)
}

// Get 获取记录
func (m *Manager) Get(userID uint32) (*History, error) {
	data, err := m.db.Find(userID)
	if err != nil {
		if err != hashdb.ErrDataStoreNotFound {
			return nil, err
		}
		m.db.Insert(userID, NewHistory())
	}

	data, err = m.db.Find(userID)
	if err != nil {
		return nil, err
	}

	r, ok := data.(*History)
	if !ok {
		return nil, errors.New("bad cast")
	}
	return r, nil
}
