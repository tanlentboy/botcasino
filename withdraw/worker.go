package withdraw

import (
	"container/list"
	"errors"
	"sync"

	"github.com/zhangpanyi/basebot/logger"
	"github.com/zhangpanyi/botcasino/models"
	"github.com/zhangpanyi/botcasino/remote"
)

// 工作者
type work struct {
	set   sync.Map
	queue *list.List
	cond  *sync.Cond
}

// 创建工作者
func newWork() *work {
	w := work{
		queue: list.New(),
		cond:  sync.NewCond(&sync.Mutex{}),
	}
	go w.loop()
	return &w
}

// 添加任务
func (w *work) addFuture(userID int64, to, assetID string, amount, fee uint32) (*Future, error) {
	// 生成订单
	orderID, err := models.InsertWithdraw(userID, to, assetID, amount, fee)
	if err != nil {
		logger.Errorf("Failed to create order, user_id=%d, to=%s, asset_id=%s, amount=%d, fee=%d, %v",
			userID, to, assetID, amount, fee, err)
		return nil, err
	}

	// 生成转账信息
	w.set.Store(orderID, nil)
	transfer := Transfer{
		UserID:  userID,
		To:      to,
		AssetID: assetID,
		Amount:  amount,
		Fee:     fee,
		OrderID: orderID,
	}
	future := newFuture(orderID, &transfer)

	// 添加到转账队列
	w.cond.L.Lock()
	defer w.cond.L.Unlock()
	w.queue.PushBack(future)
	if w.queue.Len() == 1 {
		w.cond.Signal()
	}
	return future, nil
}

// 恢复任务
func (w *work) restoreFuture(orderID int64) (*Future, error) {
	_, ok := w.set.Load(orderID)
	if ok {
		return nil, errors.New("already existed")
	}

	// 查询订单信息
	order, err := models.GetWithdrawOrder(orderID)
	if err != nil {
		return nil, err
	}

	// 检查订单状态
	if order.Status != models.WithdrawStatusWaiting {
		return nil, errors.New("order status error")
	}

	// 生成转账信息
	w.set.Store(orderID, nil)
	transfer := Transfer{
		UserID:  order.UserID,
		To:      order.To,
		AssetID: order.AssetID,
		Amount:  order.Amount,
		Fee:     order.Fee,
		OrderID: orderID,
	}
	future := newFuture(orderID, &transfer)

	// 添加到转账队列
	w.cond.L.Lock()
	defer w.cond.L.Unlock()
	w.queue.PushBack(future)
	if w.queue.Len() == 1 {
		w.cond.Signal()
	}
	return future, nil
}

// 事件循环
func (w *work) loop() {
	for {
		w.cond.L.Lock()
		for w.queue.Len() == 0 {
			w.cond.Wait()
		}
		element := w.queue.Front()
		w.queue.Remove(element)
		w.cond.L.Unlock()

		// 更新订单状态
		future := element.Value.(*Future)
		w.set.Delete(future.OrderID)
		err := models.UpdateWithdraw(future.OrderID, models.WithdrawStatusProcessing, nil)
		if err != nil {
			future.setResult(err)
			logger.Errorf("Failed to update order status as 'WithdrawStatusProcessing', order_id=%d, %v",
				future.OrderID, err)
			continue
		}

		// 执行转账操作
		err = remote.Transfer(future.OrderID, future.Transfer.To, future.Transfer.AssetID, future.Transfer.Amount)
		future.setResult(err)
		if err != nil {
			// 更新订单状态
			reason := err.Error()
			if err = models.UpdateWithdraw(future.OrderID, models.WithdrawStatusFailure, &reason); err != nil {
				logger.Errorf("Failed to update order status as 'WithdrawStatusFailure', order_id=%d, %v",
					future.OrderID, err)
				continue
			}
		}

		// 更新订单状态
		if err = models.UpdateWithdraw(future.OrderID, models.WithdrawStatusSuccessful, nil); err != nil {
			logger.Errorf("Failed to update order status as 'WithdrawStatusSuccessful', order_id=%d, %v",
				future.OrderID, err)
		}
	}
}
