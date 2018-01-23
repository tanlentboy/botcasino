package updater

import (
	"github.com/Jeffail/tunny"
)

// Pool 工作队列
type Pool struct {
	pool *tunny.WorkPool
}

// NewPool 创建工作队列
func NewPool(numWorkers int) (*Pool, error) {
	pool, err := tunny.CreatePoolGeneric(numWorkers).Open()
	if err != nil {
		return nil, err
	}
	return &Pool{pool: pool}, nil
}

// Async 添加异步任务
func (pool *Pool) Async(jobData interface{}) {
	pool.pool.SendWorkAsync(jobData, nil)
}
