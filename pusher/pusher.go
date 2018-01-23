package pusher

import (
	"container/list"
	"sync"

	"github.com/zhangpanyi/basebot/telegram/methods"
	"github.com/zhangpanyi/basebot/telegram/updater"
)

var once sync.Once
var globalPusher *telegramPusher

// CreatePusherForOnce 创建推送器
func CreatePusherForOnce(pool *updater.Pool) {
	once.Do(func() {
		globalPusher = &telegramPusher{
			pool:  pool,
			queue: list.New(),
			cond:  sync.NewCond(&sync.Mutex{}),
		}
		go globalPusher.loop()
	})
}

// 推送器
type telegramPusher struct {
	queue *list.List
	cond  *sync.Cond
	pool  *updater.Pool
}

// 推送消息
func (m *telegramPusher) push(sender *methods.BotExt, receiver int64, text string,
	markdownMode bool, markup *methods.InlineKeyboardMarkup) {

	// 构造消息结构
	msg := telegram{
		sender:       sender,
		receiver:     receiver,
		text:         text,
		markdownMode: markdownMode,
		markup:       markup,
	}

	// 添加到推送队列
	m.cond.L.Lock()
	isempty := m.queue.Len() == 0
	m.queue.PushBack(&msg)
	if isempty && m.queue.Len() == 1 {
		m.cond.Signal()
	}
	m.cond.L.Unlock()
}

// 事件循环
func (m *telegramPusher) loop() {
	for {
		m.cond.L.Lock()
		for m.queue.Len() == 0 {
			m.cond.Wait()
		}
		for m.queue.Len() > 0 {
			element := m.queue.Front()
			msg, ok := element.Value.(*telegram)
			if ok {
				m.pool.Async(msg.send)
			}
			m.queue.Remove(element)
		}
		m.cond.L.Unlock()
	}
}
