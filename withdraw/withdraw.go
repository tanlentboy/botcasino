package withdraw

import (
	"github.com/zhangpanyi/basebot/logger"
	"sync"
)

var once sync.Once
var service *withdrawService

// RunWithdrawServiceForOnce 运行提现服务
func RunWithdrawServiceForOnce(numWorkers int) {
	once.Do(func() {
		service = newWithdrawService(numWorkers)
	})
}

// AddFuture 添加任务
func AddFuture(userID int64, to, assetID string, amount uint32, fee uint32) (*Future, error) {
	return service.getWorker().addFuture(userID, to, assetID, amount, fee)
}

// RestoreFuture 恢复任务
func RestoreFuture(orderID int64) (*Future, error) {
	return service.getWorker().restoreFuture(orderID)
}

// 提现服务
type withdrawService struct {
	seq     int32
	mutex   sync.Mutex
	workers []*work
}

// 创建提现服务
func newWithdrawService(numWorkers int) *withdrawService {
	if numWorkers <= 0 {
		logger.Panicf("Invalid number of workers, %d", numWorkers)
	}

	workers := make([]*work, 0, numWorkers)
	for i := 0; i < numWorkers; i++ {
		workers = append(workers, newWork())
	}
	return &withdrawService{workers: workers}
}

// 获取工作者
func (s *withdrawService) getWorker() *work {
	var seq int32
	s.mutex.Lock()
	if s.seq < int32(len(s.workers)) {
		seq = s.seq
	} else {
		seq = 0
	}
	s.seq++
	s.mutex.Unlock()
	return s.workers[seq]
}
