package hashdb

const backlog int = 1024

// 桶结构
type bucket struct {
	id    uint32                      // 桶id
	tabel map[interface{}]interface{} // 数据表
	event chan requestEvent           // 事件通道
}

// 请求类型
type requestType int

const (
	_                 requestType = iota
	requestTypeInsert             // 插入
	requestTypeErase              // 擦除
	requestTypeUpdate             // 更新
	requestTypeFind               // 查找
	requestTypeClear              // 清除
)

// 请求事件
type requestEvent struct {
	t        requestType // 请求类型
	userData interface{} // 用户数据
}

// 插入请求
type insertRequest struct {
	key      uint32      // 键
	value    interface{} // 用户数据
	affected chan uint32 // 影响数量
}

// 擦除请求
type eraseRequest struct {
	key      uint32      // 键
	affected chan uint32 // 影响数量
}

// 更新请求
type updateRequest struct {
	key      uint32                        // 键
	update   func(interface{}) interface{} // 更新策略
	affected chan uint32                   // 影响数量
}

// 响应结果
type response struct {
	value interface{} // 用户数据
	err   error       // 错误信息
}

// 查找请求
type findRequest struct {
	key uint32        // 键
	out chan response // 查找结果
}

// 新建桶
func newBucket(id uint32) *bucket {
	event := make(chan requestEvent, backlog)
	table := make(map[interface{}]interface{})
	b := bucket{id: id, tabel: table, event: event}
	go b.loop()
	return &b
}

// 插入数据
func (b *bucket) insert(key uint32, value interface{}) uint32 {
	ch := make(chan uint32)
	userData := insertRequest{key: key, value: value, affected: ch}
	b.event <- requestEvent{t: requestTypeInsert, userData: &userData}
	affected := <-ch
	return affected
}

// 擦除数据
func (b *bucket) erase(key uint32) uint32 {
	ch := make(chan uint32)
	userData := eraseRequest{key: key, affected: ch}
	b.event <- requestEvent{t: requestTypeErase, userData: &userData}
	affected := <-ch
	return affected
}

// 更新数据
func (b *bucket) update(key uint32, fn func(interface{}) interface{}) uint32 {
	if fn == nil {
		return 0
	}

	ch := make(chan uint32)
	userData := updateRequest{key: key, update: fn, affected: ch}
	b.event <- requestEvent{t: requestTypeUpdate, userData: &userData}
	affected := <-ch
	return affected
}

// 查找数据
func (b *bucket) find(key uint32) (interface{}, error) {
	ch := make(chan response)
	userData := findRequest{key: key, out: ch}
	b.event <- requestEvent{t: requestTypeFind, userData: &userData}
	result := <-ch
	return result.value, result.err
}

// 清理数据
func (b *bucket) clear() {
	b.event <- requestEvent{t: requestTypeClear}
}

// 事件循环
func (b *bucket) loop() {
	for {
		select {
		case event := <-b.event:
			switch event.t {
			// 插入
			case requestTypeInsert:
				b.handleInsertInGoroutine(event.userData.(*insertRequest))
			// 擦除
			case requestTypeErase:
				b.handleEraseInGoroutine(event.userData.(*eraseRequest))
			// 更新
			case requestTypeUpdate:
				b.handleUpdateInGoroutine(event.userData.(*updateRequest))
			// 查找
			case requestTypeFind:
				b.handleFindInGoroutine(event.userData.(*findRequest))
			// 清理
			case requestTypeClear:
				b.handleClearInGoroutine()
			}
		}
	}
}

// 处理插入
func (b *bucket) handleInsertInGoroutine(req *insertRequest) {
	_, ok := b.tabel[req.key]
	b.tabel[req.key] = req.value
	if ok {
		req.affected <- 0
	} else {
		req.affected <- 1
	}
}

// 处理擦除
func (b *bucket) handleEraseInGoroutine(req *eraseRequest) {
	if _, ok := b.tabel[req.key]; ok {
		delete(b.tabel, req.key)
		req.affected <- 1
	} else {
		req.affected <- 0
	}
}

// 处理更新
func (b *bucket) handleUpdateInGoroutine(req *updateRequest) {
	var newValue interface{}
	old, ok := b.tabel[req.key]
	if !ok {
		newValue = req.update(nil)
	} else {
		newValue = req.update(old)
	}

	if newValue != nil {
		b.tabel[req.key] = newValue
		req.affected <- 1
	} else {
		delete(b.tabel, req.key)
		req.affected <- 0
	}
}

// 处理查找
func (b *bucket) handleFindInGoroutine(req *findRequest) {
	value, ok := b.tabel[req.key]
	if ok {
		req.out <- response{value: value}
	} else {
		req.out <- response{err: ErrDataStoreNotFound}
	}
}

// 处理清理
func (b *bucket) handleClearInGoroutine() {
	b.tabel = make(map[interface{}]interface{})
}
