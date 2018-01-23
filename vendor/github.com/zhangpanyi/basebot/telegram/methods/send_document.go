package methods

import (
	"encoding/json"
	"strconv"

	"github.com/zhangpanyi/basebot/telegram/types"
)

// 发送文档
type sendDocument struct {
	ChatID           int64       `json:"chat_id"`                       // 聊天ID
	Document         string      `json:"document"`                      // 文件ID
	Caption          string      `json:"caption"`                       // 文件标题
	ReplyToMessageID int32       `json:"reply_to_message_id,omitempty"` // 回复消息ID
	ReplyMarkup      interface{} `json:"reply_markup,omitempty"`        // 回复标记
}

// SendDocument 发送文档
func (bot *BotExt) SendDocument(chatID int64, caption string, fileID string,
	markup *InlineKeyboardMarkup) (*types.Message, error) {

	request := sendDocument{
		ChatID:   chatID,
		Document: fileID,
		Caption:  caption,
	}

	if markup != nil {
		request.ReplyMarkup = markup
	}

	res, err := bot.Call("sendDocument", request)
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

// SendDocumentFile 发送文档文件
func (bot *BotExt) SendDocumentFile(chatID int64, caption string, file []byte,
	filename string, markup *InlineKeyboardMarkup) (*types.Message, error) {

	// 生成请求内容
	formdata := make([]Field, 0)
	formdata = append(formdata, Field{Name: "chat_id", Text: strconv.Itoa(int(chatID))})
	formdata = append(formdata, Field{Name: "document", File: file, FileName: filename})
	if len(caption) > 0 {
		formdata = append(formdata, Field{Name: "caption", Text: caption})
	}
	if markup != nil {
		if jsb, err := markup.ToJSON(); err == nil {
			formdata = append(formdata, Field{Name: "reply_markup", Text: string(jsb)})
		}
	}

	res, err := bot.Upload("sendDocument", formdata)
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
