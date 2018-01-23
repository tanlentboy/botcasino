package history

import (
	"errors"
	"sync"

	"github.com/zhangpanyi/basebot/telegram/types"
)

// DefaultCapacity 默认容量
const DefaultCapacity = 16

// History 历史记录
type History struct {
	size     int
	capacity int
	history  []*types.Update
	mutex    sync.RWMutex
}

// NewHistory 创建记录
func NewHistory() *History {
	return &History{
		capacity: DefaultCapacity,
		history:  make([]*types.Update, DefaultCapacity),
	}
}

// Clear 清空历史
func (r *History) Clear() *History {
	r.mutex.RLock()
	defer func() {
		r.mutex.RUnlock()
	}()
	r.size = 0
	return r
}

// Empty 是否为空
func (r *History) Empty() bool {
	r.mutex.RLock()
	defer func() {
		r.mutex.RUnlock()
	}()
	return r.size == 0
}

// Push 插入数据
func (r *History) Push(update *types.Update) *History {
	r.mutex.Lock()
	defer func() {
		r.mutex.Unlock()
	}()
	if r.size == r.capacity {
		r.capacity <<= 1
		history := make([]*types.Update, r.capacity)
		copy(history, r.history)
		r.history = history
	}
	r.history[r.size] = update
	r.size++
	return r
}

// Pop 擦除数据
func (r *History) Pop() *History {
	r.mutex.Lock()
	defer func() {
		r.mutex.Unlock()
	}()
	if r.size > 0 {
		r.size--
	}
	return r
}

// Back 最后元素
func (r *History) Back() (*types.Update, error) {
	r.mutex.RLock()
	defer func() {
		r.mutex.RUnlock()
	}()
	if r.size > 0 {
		return r.history[r.size-1], nil
	}
	return nil, errors.New("not found")
}

// Foreach 遍历数据
func (r *History) Foreach(callback func(int, *types.Update) bool) {
	r.mutex.RLock()
	history := make([]*types.Update, r.size)
	copy(history, r.history[0:r.size])
	r.mutex.RUnlock()

	for i := r.size - 1; i >= 0; i-- {
		if !callback(r.size-1-i, history[i]) {
			break
		}
	}
}

// LastCallbackQuery 最后的查询回调
func (r *History) LastCallbackQuery() (*types.CallbackQuery, error) {
	r.mutex.RLock()
	defer func() {
		r.mutex.RUnlock()
	}()

	for i := r.size - 1; i >= 0; i-- {
		if r.history[i].CallbackQuery != nil {
			return r.history[i].CallbackQuery, nil
		}
	}
	return nil, errors.New("not found")
}
