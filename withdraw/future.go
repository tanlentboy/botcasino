package withdraw

// 转账信息
type Transfer struct {
	UserID  int64  // 用户ID
	To      string // 帐户名
	AssetID string // 资产ID
	Amount  uint32 // 资产数量(分)
	Fee     uint32 // 手续费(分)
	OrderID int64  // 订单ID
}

// Future 任务
type Future struct {
	ch       chan error
	OrderID  int64
	Transfer *Transfer
}

// 创建任务
func newFuture(orderID int64, transfer *Transfer) *Future {
	ch := make(chan error)
	return &Future{ch: ch, OrderID: orderID, Transfer: transfer}
}

// GetResult 获取结果
func (f *Future) GetResult() error {
	err := <-f.ch
	return err
}

// 设置结果
func (f *Future) setResult(err error) {
	f.ch <- err
}
