package types

import (
	"encoding/json"
)

// User 用户信息
type User struct {
	ID           int64  `json:"id"`                      // 用户唯一ID
	IsBot        bool   `json:"is_bot"`                  // 是否机器人
	FirstName    string `json:"first_name"`              // First name
	LastName     string `json:"last_name,omitempty"`     // Last name
	UserName     string `json:"username,omitempty"`      // 用户名
	LanguageCode string `json:"language_code,omitempty"` // 语言代码
}

// ToJSON 转换为JSON
func (user *User) ToJSON() ([]byte, error) {
	return json.Marshal(user)
}

// FromJSON 从JSON反序列化
func (user *User) FromJSON(jsb []byte) error {
	return json.Unmarshal(jsb, user)
}
