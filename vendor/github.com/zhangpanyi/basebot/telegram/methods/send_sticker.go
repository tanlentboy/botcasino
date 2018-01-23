package methods

import (
	"encoding/json"
	"strconv"

	"github.com/zhangpanyi/basebot/telegram/types"
)

// 发送贴纸
type sendSticker struct {
	ChatID           int64       `json:"chat_id"`                       // 聊天ID
	Sticker          string      `json:"sticker"`                       // 文件ID
	ReplyToMessageID int32       `json:"reply_to_message_id,omitempty"` // 回复消息ID
	ReplyMarkup      interface{} `json:"reply_markup,omitempty"`        // 回复标记
}

// SendSticker 发送贴纸
func (bot *BotExt) SendSticker(chatID int64, fileID string,
	markup *InlineKeyboardMarkup) (*types.Message, error) {

	request := sendSticker{
		ChatID:  chatID,
		Sticker: fileID,
	}

	if markup != nil {
		request.ReplyMarkup = markup
	}

	res, err := bot.Call("sendSticker", request)
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

// SendStickerFile 发送贴纸文件
func (bot *BotExt) SendStickerFile(chatID int64, file []byte,
	markup *InlineKeyboardMarkup) (*types.Message, error) {

	// 生成请求内容
	formdata := make([]Field, 0)
	formdata = append(formdata, Field{Name: "chat_id", Text: strconv.Itoa(int(chatID))})
	formdata = append(formdata, Field{Name: "sticker", File: file, FileName: "sticker.webp"})
	if markup != nil {
		if jsb, err := markup.ToJSON(); err == nil {
			formdata = append(formdata, Field{Name: "reply_markup", Text: string(jsb)})
		}
	}

	res, err := bot.Upload("sendSticker", formdata)
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
