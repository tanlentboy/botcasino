package expiretimer

import (
	"container/heap"
	"fmt"
	"sync"
	"time"

	"github.com/zhangpanyi/botcasino/config"
	"github.com/zhangpanyi/botcasino/envelopes/groupchat"
	"github.com/zhangpanyi/botcasino/models"
	"github.com/zhangpanyi/botcasino/storage"

	"github.com/zhangpanyi/basebot/logger"
	"github.com/zhangpanyi/basebot/telegram/methods"
	"github.com/zhangpanyi/basebot/telegram/updater"
)

var once sync.Once
var globalExpireTimer *expireTimer

// StartTimerForOnce 启动定时器
func StartTimerForOnce(bot *methods.BotExt, pool *updater.Pool) {
	once.Do(func() {
		// 获取最后过期红包
		handler := storage.RedEnvelopeStorage{}
		id, err := handler.GetLastExpired()
		if err != nil && err != storage.ErrNoBucket {
			logger.Panic(err)
		}

		// 遍历未过期列表
		h := make(expireHeap, 0)
		err = handler.ForeachRedEnvelopes(id+1, func(data *storage.RedEnvelope) {
			heap.Push(&h, expire{ID: data.ID, Timestamp: data.Timestamp})
		})
		if err != nil && err != storage.ErrNoBucket {
			logger.Panic(err)
		}

		// 初始化过期定时器
		globalExpireTimer = &expireTimer{
			h:    h,
			bot:  bot,
			pool: pool,
		}
		go globalExpireTimer.loop()
	})
}

// GetBot 获取机器人
func GetBot() *methods.BotExt {
	return globalExpireTimer.bot
}

// AddRedEnvelope 添加红包
func AddRedEnvelope(id uint64, timestamp int64) {
	globalExpireTimer.lock.Lock()
	defer globalExpireTimer.lock.Unlock()
	heap.Push(&globalExpireTimer.h, expire{ID: id, Timestamp: timestamp})
}

// 过期定时器
type expireTimer struct {
	h    expireHeap
	bot  *methods.BotExt
	pool *updater.Pool
	lock sync.RWMutex
}

// 处理过期红包
func (t *expireTimer) handleRedEnvelopeExpire() {
	now := time.Now().Unix()
	dynamicCfg := config.GetDynamic()

	var id uint64
	t.lock.RLock()
	for t.h.Len() > 0 {
		data := t.h.Front()
		t.lock.RUnlock()

		// 判断是否过期
		if now-data.Timestamp < dynamicCfg.RedEnvelopeExpire {
			return
		}

		// 获取过期信息
		t.lock.Lock()
		e := heap.Pop(&t.h).(expire)
		t.lock.Unlock()

		id = e.ID
		logger.Infof("On red envelope expired, %v", e.Timestamp)
		t.pool.Async(func() {
			t.handleRedEnvelopeExpireAsync(e.ID)
		})
		t.lock.RLock()
	}
	t.lock.RUnlock()

	// 更新过期红包
	if id != 0 {
		handler := storage.RedEnvelopeStorage{}
		if err := handler.SetLastExpired(id); err != nil {
			logger.Warnf("Failed to set last expired  of red envelope, %v", err)
		}
	}
}

// 异步处理过期红包
func (t *expireTimer) handleRedEnvelopeExpireAsync(id uint64) {
	// 设置红包过期
	handler := storage.RedEnvelopeStorage{}
	if handler.IsExpired(id) {
		return
	}
	err := handler.SetExpired(id)
	if err != nil {
		logger.Infof("Failed to set expired of red envelope, %v", err)
		return
	}

	// 获取红包信息
	redEnvelope, received, err := handler.GetRedEnvelope(id)
	if err != nil {
		logger.Warnf("Failed to set expired of red envelope, not found red envelope, %d, %v", id, err)
		return
	}
	if received == redEnvelope.Number {
		return
	}

	// 计算红包余额
	balance := redEnvelope.Amount - redEnvelope.Received
	if !redEnvelope.Lucky {
		balance = redEnvelope.Amount*redEnvelope.Number - redEnvelope.Received
	}

	// 返还红包余额
	assetHandler := storage.AssetStorage{}
	err = assetHandler.UnfreezeAsset(redEnvelope.SenderID, redEnvelope.Asset, balance)
	if err != nil {
		logger.Errorf("Failed to return red envelope asset of expired, %v", err)
	} else {
		logger.Errorf("Return red envelope asset of expired, UserID=%d, Asset=%s, Amount=%d",
			redEnvelope.SenderID, redEnvelope.Asset, balance)
		desc := fmt.Sprintf("您发放的红包(id: *%d*)过期无人领取, 退还余额*%.2f* *%s*", redEnvelope.ID,
			float64(balance)/100.0, redEnvelope.Asset)
		models.InsertHistory(redEnvelope.SenderID, desc)
	}

	// 更新聊天信息
	groupchat.UpdateRedEnvelope(t.bot, redEnvelope, received)
}

// 事件循环
func (t *expireTimer) loop() {
	tickTimer := time.NewTimer(time.Second)
	for {
		select {
		case <-tickTimer.C:
			t.handleRedEnvelopeExpire()
			tickTimer.Reset(time.Second)
		}
	}
}
