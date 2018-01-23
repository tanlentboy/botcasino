package types

import (
	"encoding/json"
)

const (
	// EntityMention 提及用户
	EntityMention = "mention "
	// EntityHashTag 哈希标记
	EntityHashTag = "hashtag"
	// EntityBotCommand 机器人命令
	EntityBotCommand = "bot_command"
	// EntityURL 超链接
	EntityURL = "url"
	// EntityEmail 电子邮箱
	EntityEmail = "email"
	// EntityBold 加粗
	EntityBold = "bold"
	// EntityItalic 斜体
	EntityItalic = "italic"
	// EntityCode 代码
	EntityCode = "code"
	// EntityPre 预格式化文本
	EntityPre = "pre"
	// EntityTextLink 文本链接
	EntityTextLink = "text_link"
	// EntityTextMention 文本提及
	EntityTextMention = "text_mention "
)

// MessageEntity Entity信息
type MessageEntity struct {
	Type   string `json:"type"`           // Entity类型
	Offset int32  `json:"offset"`         // 偏移
	Length int32  `json:"length"`         // 长度
	URL    string `json:"url,omitempty"`  // 地址
	User   *User  `json:"user,omitempty"` // 用户信息
}

// ToJSON 转换为JSON
func (entity *MessageEntity) ToJSON() ([]byte, error) {
	return json.Marshal(entity)
}

// FromJSON 从JSON反序列化
func (entity *MessageEntity) FromJSON(jsb []byte) error {
	return json.Unmarshal(jsb, entity)
}
