package methods

import (
	"github.com/zhangpanyi/basebot/telegram/types"
)

// 应答查询回调
type answerCallbackQuery struct {
	CallbackQueryID string `json:"callback_query_id"`    // 查询回调的唯一ID
	Text            string `json:"text,omitempty"`       // 通知文本
	ShowAlert       bool   `json:"show_alert,omitempty"` // 显示警告
	URL             string `json:"url,omitempty"`        // 打开URL
	CacheTime       int32  `json:"cache_time,omitempty"` // 缓存时间
}

// AnswerCallbackQuery 应答查询回调
func (bot *BotExt) AnswerCallbackQuery(query *types.CallbackQuery, text string, alert bool, url string, cacheTime int32) error {
	request := answerCallbackQuery{
		CallbackQueryID: query.ID,
		Text:            text,
		ShowAlert:       alert,
		URL:             url,
		CacheTime:       cacheTime,
	}
	_, err := bot.Call("answerCallbackQuery", &request)
	return err
}
