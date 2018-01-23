package pusher

import (
	"github.com/zhangpanyi/basebot/telegram/methods"
)

// To 推送消息
func To(sender *methods.BotExt, receiver int64, text string,
	markdownMode bool, markup *methods.InlineKeyboardMarkup) {

	globalPusher.push(sender, receiver, text, markdownMode, markup)
}
