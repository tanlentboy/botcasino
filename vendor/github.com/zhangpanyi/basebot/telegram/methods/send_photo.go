package methods

import (
	"encoding/json"
	"strconv"

	"github.com/zhangpanyi/basebot/telegram/types"
)

// 发送照片
type sendPhoto struct {
	ChatID           int64       `json:"chat_id"`                       // 聊天ID
	Photo            string      `json:"photo"`                         // 文件ID
	Caption          string      `json:"caption"`                       // 照片标题
	ReplyToMessageID int32       `json:"reply_to_message_id,omitempty"` // 回复消息ID
	ReplyMarkup      interface{} `json:"reply_markup,omitempty"`        // 回复标记
}

// SendPhoto 发送照片
func (bot *BotExt) SendPhoto(chatID int64, caption string, fileID string,
	markup *InlineKeyboardMarkup) (*types.Message, error) {

	request := sendPhoto{
		ChatID:  chatID,
		Photo:   fileID,
		Caption: caption,
	}

	if markup != nil {
		request.ReplyMarkup = markup
	}

	res, err := bot.Call("sendPhoto", request)
	if err != nil {
		return nil, err
	}

	resonpe := SendMessageResonpe{}
	err = json.Unmarshal(res, &resonpe)
	if err != nil {
		return nil, err
	}

	return resonpe.Result, nil
}

// ReplyPhoto 回复照片
func (bot *BotExt) ReplyPhoto(message *types.Message, caption string, fileID string,
	markup *InlineKeyboardMarkup) (*types.Message, error) {

	request := sendPhoto{
		ChatID:           message.Chat.ID,
		Photo:            fileID,
		Caption:          caption,
		ReplyToMessageID: message.MessageID,
	}

	if markup != nil {
		request.ReplyMarkup = markup
	}

	res, err := bot.Call("sendPhoto", request)
	if err != nil {
		return nil, err
	}

	resonpe := SendMessageResonpe{}
	err = json.Unmarshal(res, &resonpe)
	if err != nil {
		return nil, err
	}

	return resonpe.Result, nil
}

// SendPhotoFile 发送照片文件
func (bot *BotExt) SendPhotoFile(chatID int64, caption string, file []byte,
	markup *InlineKeyboardMarkup) (*types.Message, error) {

	// 生成请求内容
	formdata := make([]Field, 0)
	formdata = append(formdata, Field{Name: "chat_id", Text: strconv.FormatInt(chatID, 10)})
	formdata = append(formdata, Field{Name: "photo", File: file, FileName: "photo.jpeg"})
	if len(caption) > 0 {
		formdata = append(formdata, Field{Name: "caption", Text: caption})
	}
	if markup != nil {
		if jsb, err := markup.ToJSON(); err == nil {
			formdata = append(formdata, Field{Name: "reply_markup", Text: string(jsb)})
		}
	}

	res, err := bot.Upload("sendPhoto", formdata)
	if err != nil {
		return nil, err
	}

	resonpe := SendMessageResonpe{}
	err = json.Unmarshal(res, &resonpe)
	if err != nil {
		return nil, err
	}

	return resonpe.Result, nil
}
