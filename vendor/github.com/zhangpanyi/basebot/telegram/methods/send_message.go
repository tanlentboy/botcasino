package methods

import (
	"encoding/json"

	"github.com/zhangpanyi/basebot/telegram/types"
)

// KeyboardButton 键盘按钮
type KeyboardButton struct {
	Text            string `json:"text"`                       // 按钮文本标签
	RequestContact  bool   `json:"request_contact,omitempty"`  // 发送联系人信息
	RequestLocation bool   `json:"request_location,omitempty"` // 发送位置信息
}

// ReplyKeyboardMarkup 回复键盘
type ReplyKeyboardMarkup struct {
	Keyboard        [][]*KeyboardButton `json:"keyboard"`                    // 按钮列表
	ResizeKeyboard  bool                `json:"resize_keyboard,omitempty"`   // 自动设置大小
	OneTimeKeyboard bool                `json:"one_time_keyboard,omitempty"` // 使用键盘时隐藏自定义键盘
	Selective       bool                `json:"selective,omitempty"`         // 特定用户可见
}

// ToJSON 转换为
func (markup *ReplyKeyboardMarkup) ToJSON() ([]byte, error) {
	jsb, err := json.Marshal(markup)
	if err != nil {
		return nil, err
	}
	return jsb, nil
}

// MakeReplyKeyboardMarkup 生成回复键盘
func MakeReplyKeyboardMarkup(menus []KeyboardButton, columns ...int) *ReplyKeyboardMarkup {
	offset := 0
	replyKeyboardMarkup := ReplyKeyboardMarkup{
		Keyboard:        make([][]*KeyboardButton, 0),
		ResizeKeyboard:  true,
		OneTimeKeyboard: true,
	}

	for i := 0; i < len(columns); i++ {
		col := columns[i]
		if len(menus)-col-offset < 0 {
			return &replyKeyboardMarkup
		}

		replyKeyboard := make([]*KeyboardButton, 0, col)
		for j := offset; j < len(menus) && j < offset+col; j++ {
			button := menus[j]
			replyKeyboard = append(replyKeyboard, &button)
		}
		replyKeyboardMarkup.Keyboard = append(replyKeyboardMarkup.Keyboard, replyKeyboard)
		offset += col
	}
	return &replyKeyboardMarkup
}

// MakeReplyKeyboardMarkupAuto 生成回复键盘(自动排版)
func MakeReplyKeyboardMarkupAuto(menus []KeyboardButton, columns uint) *ReplyKeyboardMarkup {
	if len(menus) <= 0 || columns <= 0 {
		return nil
	}

	rows := len(menus) / int(columns)
	if len(menus)%int(columns) != 0 {
		rows++
	}

	replyKeyboardMarkup := ReplyKeyboardMarkup{
		Keyboard:        make([][]*KeyboardButton, rows),
		ResizeKeyboard:  true,
		OneTimeKeyboard: true,
	}
	for i := 0; i < len(menus); i++ {
		idx := i / int(columns)
		if replyKeyboardMarkup.Keyboard[idx] == nil {
			replyKeyboardMarkup.Keyboard[idx] = make([]*KeyboardButton, 0)
		}

		button := menus[i]
		replyKeyboardMarkup.Keyboard[idx] = append(replyKeyboardMarkup.Keyboard[idx], &button)
	}
	return &replyKeyboardMarkup
}

// ReplyKeyboardRemove 删除回复键盘
type ReplyKeyboardRemove struct {
	RemoveKeyboard bool `json:"remove_keyboard"`     // 删除键盘
	Selective      bool `json:"selective,omitempty"` // 指定人可见
}

// ToJSON 转换为
func (markup *ReplyKeyboardRemove) ToJSON() ([]byte, error) {
	jsb, err := json.Marshal(markup)
	if err != nil {
		return nil, err
	}
	return jsb, nil
}

// InlineKeyboardButton 内联键盘按钮
type InlineKeyboardButton struct {
	Text                         string `json:"text"`                                       // 按钮文本标签
	URL                          string `json:"url,omitempty"`                              // 打开地址
	CallbackData                 string `json:"callback_data,omitempty"`                    // 回调数据
	SwitchInlineQuery            string `json:"switch_inline_query,omitempty"`              // 切换内联查询
	SwitchInlineQueryCurrentChat string `json:"switch_inline_query_current_chat,omitempty"` // 切换内联查询到当前聊天
}

// InlineKeyboardMarkup 内联键盘
type InlineKeyboardMarkup struct {
	InlineKeyboard [][]*InlineKeyboardButton `json:"inline_keyboard"` // 按钮列表
}

// ToJSON 转换为
func (markup *InlineKeyboardMarkup) ToJSON() ([]byte, error) {
	jsb, err := json.Marshal(markup)
	if err != nil {
		return nil, err
	}
	return jsb, nil
}

// Merge 合并内联键盘
func (markup *InlineKeyboardMarkup) Merge(args ...*InlineKeyboardMarkup) *InlineKeyboardMarkup {
	for _, item := range args {
		for _, rows := range item.InlineKeyboard {
			markup.InlineKeyboard = append(markup.InlineKeyboard, rows)
		}
	}
	return markup
}

// MakeInlineKeyboardMarkup 生成回调键盘
func MakeInlineKeyboardMarkup(menus []InlineKeyboardButton, columns ...int) *InlineKeyboardMarkup {
	offset := 0
	inlineKeyboardMarkup := InlineKeyboardMarkup{
		InlineKeyboard: make([][]*InlineKeyboardButton, 0),
	}

	for i := 0; i < len(columns); i++ {
		col := columns[i]
		if len(menus)-col-offset < 0 {
			return &inlineKeyboardMarkup
		}

		inlineKeyboard := make([]*InlineKeyboardButton, 0, col)
		for j := offset; j < len(menus) && j < offset+col; j++ {
			button := menus[j]
			inlineKeyboard = append(inlineKeyboard, &button)
		}
		inlineKeyboardMarkup.InlineKeyboard = append(inlineKeyboardMarkup.InlineKeyboard, inlineKeyboard)
		offset += col
	}
	return &inlineKeyboardMarkup
}

// MakeInlineKeyboardMarkupAuto 生成回调内联键盘(自动排版)
func MakeInlineKeyboardMarkupAuto(menus []InlineKeyboardButton, columns uint) *InlineKeyboardMarkup {
	if len(menus) <= 0 || columns <= 0 {
		return nil
	}

	rows := len(menus) / int(columns)
	if len(menus)%int(columns) != 0 {
		rows++
	}

	inlineKeyboardMarkup := InlineKeyboardMarkup{
		InlineKeyboard: make([][]*InlineKeyboardButton, rows),
	}
	for i := 0; i < len(menus); i++ {
		idx := i / int(columns)
		if inlineKeyboardMarkup.InlineKeyboard[idx] == nil {
			inlineKeyboardMarkup.InlineKeyboard[idx] = make([]*InlineKeyboardButton, 0)
		}

		button := menus[i]
		if len(menus[i].SwitchInlineQuery) > 0 {
			button.SwitchInlineQuery = menus[i].SwitchInlineQuery
		}
		if len(menus[i].SwitchInlineQueryCurrentChat) > 0 {
			button.SwitchInlineQueryCurrentChat = menus[i].SwitchInlineQueryCurrentChat
		}
		inlineKeyboardMarkup.InlineKeyboard[idx] = append(inlineKeyboardMarkup.InlineKeyboard[idx], &button)
	}
	return &inlineKeyboardMarkup
}

const (
	// ParseModeHTML 解析HTML
	ParseModeHTML = "HTML"
	// ParseModeMarkdown 解析Markdown
	ParseModeMarkdown = "Markdown"
)

// 发送消息
type sendMessage struct {
	ChatID                int64       `json:"chat_id"`                            // 聊天ID
	Text                  string      `json:"text"`                               // 消息文本
	ParseMode             string      `json:"parse_mode,omitempty"`               // 解析模式
	ReplyToMessageID      int32       `json:"reply_to_message_id,omitempty"`      // 回复消息ID
	DisableWebPagePreview bool        `json:"disable_web_page_preview,omitempty"` // 禁用网页预览
	ReplyMarkup           interface{} `json:"reply_markup,omitempty"`             // 回复标记
}

// 发送消息
func (bot *BotExt) sendMessage(request *sendMessage) (*types.Message, error) {
	res, err := bot.Call("sendMessage", request)
	if err != nil {
		return nil, err
	}

	if request.ReplyMarkup != nil {
		switch real := request.ReplyMarkup.(type) {
		case *InlineKeyboardMarkup:
			if real != nil {
				request.ReplyMarkup = real
			}
		case *ReplyKeyboardMarkup:
			if real != nil {
				request.ReplyMarkup = real
			}
		case *ReplyKeyboardRemove:
			if real != nil {
				request.ReplyMarkup = real
			}
		default:
			request.ReplyMarkup = nil
		}
	}

	resonpe := SendMessageResonpe{}
	err = json.Unmarshal(res, &resonpe)
	if err != nil {
		return nil, err
	}
	return resonpe.Result, nil
}

// SendMessage 发送消息
func (bot *BotExt) SendMessage(chatID int64, text string, mdMode bool,
	markup interface{}) (*types.Message, error) {
	request := sendMessage{
		ChatID:      chatID,
		Text:        text,
		ReplyMarkup: markup,
	}
	if mdMode {
		request.ParseMode = ParseModeMarkdown
	}
	return bot.sendMessage(&request)
}

// ReplyMessage 回复消息
func (bot *BotExt) ReplyMessage(message *types.Message, text string, mdMode bool,
	markup interface{}) (*types.Message, error) {
	request := sendMessage{
		ChatID:           message.Chat.ID,
		Text:             text,
		ReplyToMessageID: message.MessageID,
		ReplyMarkup:      markup,
	}
	if mdMode {
		request.ParseMode = ParseModeMarkdown
	}
	return bot.sendMessage(&request)
}

// SendMessageDisableWebPagePreview 发送消息并禁用网页预览
func (bot *BotExt) SendMessageDisableWebPagePreview(chatID int64, text string, mdMode bool,
	markup interface{}) (*types.Message, error) {
	request := sendMessage{
		ChatID: chatID,
		Text:   text,
		DisableWebPagePreview: true,
		ReplyMarkup:           markup,
	}
	if mdMode {
		request.ParseMode = ParseModeMarkdown
	}
	return bot.sendMessage(&request)
}

// ReplyMessageDisableWebPagePreview 回复消息并禁用网页预览
func (bot *BotExt) ReplyMessageDisableWebPagePreview(message *types.Message, text string, mdMode bool,
	markup *InlineKeyboardMarkup) (*types.Message, error) {
	request := sendMessage{
		ChatID:                message.Chat.ID,
		Text:                  text,
		ReplyToMessageID:      message.MessageID,
		DisableWebPagePreview: true,
		ReplyMarkup:           markup,
	}
	if mdMode {
		request.ParseMode = ParseModeMarkdown
	}
	return bot.sendMessage(&request)
}
