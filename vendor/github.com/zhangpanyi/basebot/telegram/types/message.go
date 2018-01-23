package types

import (
	"encoding/json"
)

// Message 消息信息
type Message struct {
	MessageID            int32            `json:"message_id"`                        // 消息唯一ID
	From                 *User            `json:"from,omitempty"`                    // 消息来源
	Date                 int32            `json:"date"`                              // 消息日期
	Chat                 *Chat            `json:"chat"`                              // 聊天信息
	ForwardFrom          *User            `json:"forward_from,omitempty"`            // 转发消息来源
	ForwardFromChat      *Chat            `json:"forward_from_chat,omitempty"`       // 转发消息来源聊天信息
	ForwardFromMessageID int32            `json:"forward_from_message_id,omitempty"` // 转发消息唯一ID
	ForwardSignature     string           `json:"forward_signature,omitempty"`       // 转发消息签名
	ForwardDate          int32            `json:"forward_date,omitempty"`            // 转发日期
	ReplyToMessage       *Message         `json:"reply_to_message,omitempty"`        // 回复的消息
	EditDate             int32            `json:"edit_date,omitempty"`               // 编辑日期
	AuthorSignature      string           `json:"author_signature,omitempty"`        // 作者签名
	Text                 string           `json:"text,omitempty"`                    // 消息文本
	Entities             []*MessageEntity `json:"entities,omitempty"`                // 文本Entity信息
	Caption              string           `json:"caption,omitempty"`                 // 媒体标题
	CaptionEntities      []*MessageEntity `json:"caption_entities,omitempty"`        // 标题Entity信息
	Photo                []*PhotoSize     `json:"photo,omitempty"`                   // 照片尺寸
	Sticker              *StickerSet      `json:"sticker,omitempty"`                 // 贴纸集合
	Document             *Documnet        `json:"document,omitempty"`                // 文档信息
	Contact              *Contact         `json:"contact,omitempty"`                 // 联系人
}

// ToJSON 转换为JSON
func (message *Message) ToJSON() ([]byte, error) {
	return json.Marshal(message)
}

// FromJSON 从JSON反序列化
func (message *Message) FromJSON(jsb []byte) error {
	return json.Unmarshal(jsb, message)
}
