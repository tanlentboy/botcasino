package updater

import (
	"container/list"
	"log"
	"time"

	"github.com/zhangpanyi/basebot/telegram/methods"
	"github.com/zhangpanyi/basebot/telegram/types"
)

// Queue 消息队列
type Queue struct {
	workers []*worker
}

// NewQueue 创建消息队列
func NewQueue(numWorkers uint32) *Queue {
	var i uint32
	queue := Queue{workers: make([]*worker, 0, numWorkers)}
	for i = 0; i < numWorkers; i++ {
		w := worker{}
		w.queue = newMultimap()
		w.ch = make(chan future, 128)
		queue.workers = append(queue.workers, &w)
		w.startPolling()
	}
	return &queue
}

// Put 添加更新
func (queue *Queue) Put(f future) {
	idx := f.bot.ID % int64(len(queue.workers))
	queue.workers[idx].put(f)
}

// 任务
type future struct {
	handler Handler
	bot     *methods.BotExt
	update  *types.Update
	pool    *Pool
}

// 异步调用
func (f *future) asyncCall() {
	f.pool.Async(func() {
		f.handler(f.bot, f.update)
	})
}

// 多值容器
type multimap struct {
	orders  *list.List            // 序列列表
	updates map[uint32]*list.List // 更新集合
}

// 创建multimap
func newMultimap() *multimap {
	return &multimap{
		orders:  list.New(),
		updates: make(map[uint32]*list.List),
	}
}

// 取出数据
func (m *multimap) popFront() (future, bool) {
	if m.orders.Len() == 0 {
		return future{}, false
	}

	// 获取key
	element := m.orders.Front()
	if element == nil || element.Value == nil {
		return future{}, false
	}

	key, ok := element.Value.(uint32)
	if !ok {
		return future{}, false
	}

	// 获取数据
	updates, ok := m.updates[key]
	if !ok {
		return future{}, false
	}

	defer func() {
		if updates.Len() == 0 {
			delete(m.updates, key)
			m.orders.Remove(element)
		}
	}()

	element = updates.Front()
	updates.Remove(element)
	if element == nil || element.Value == nil {
		return future{}, false
	}

	f, ok := element.Value.(future)
	if !ok {
		return future{}, false
	}

	return f, true
}

// 插入数据
func (m *multimap) insert(key uint32, f future) {
	updates, ok := m.updates[key]
	if !ok {
		updates = list.New()
		m.orders.PushBack(key)
		m.updates[key] = updates
	}
	updates.PushBack(f)
}

// 工作者
type worker struct {
	ch    chan future // 压入通道
	queue *multimap   // 工作队列
}

// 添加任务
func (w *worker) put(f future) {
	w.ch <- f
}

// 处理任务
func (w *worker) consume() {
	f, ok := w.queue.popFront()
	if !ok {
		time.Sleep(10 * time.Millisecond)
		return
	}

	defer func() {
		if err := recover(); err != nil {
			log.Printf("Telegram updater: panic, %v", err)
		}
	}()

	f.asyncCall()
}

// 开始轮询
func (w *worker) startPolling() {
	go func() {
		for {
			select {
			case data := <-w.ch:
				// 插入数据
				w.queue.insert(uint32(data.bot.ID), data)
			default:
				// 处理数据
				w.consume()
			}
		}
	}()
}
