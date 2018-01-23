package types

import (
	"encoding/json"
)

// Update 更新信息
type Update struct {
	UpdateID      int32          `json:"update_id"`                // 更新ID
	Message       *Message       `json:"message,omitempty"`        // 消息
	CallbackQuery *CallbackQuery `json:"callback_query,omitempty"` // 查询回调
}

// ToJSON 转换为JSON
func (update *Update) ToJSON() ([]byte, error) {
	return json.Marshal(update)
}

// FromJSON 从JSON反序列化
func (update *Update) FromJSON(jsb []byte) error {
	return json.Unmarshal(jsb, update)
}
