package pusher

import (
	"github.com/zhangpanyi/basebot/logger"
	"github.com/zhangpanyi/basebot/telegram/methods"
)

// 报文内容
type telegram struct {
	sender       *methods.BotExt               // 发送者
	receiver     int64                         // 接收者
	text         string                        // 文本
	markdownMode bool                          // MarkDown渲染
	markup       *methods.InlineKeyboardMarkup // 内联键盘标记
}

// 发送消息
func (msg *telegram) send() {
	_, err := msg.sender.SendMessage(msg.receiver, msg.text, msg.markdownMode, msg.markup)
	if err != nil {
		logger.Warnf("Failed to push message, %v", err)
	}
}
