package types

import (
	"encoding/json"
)

// PhotoSize 照片大小
type PhotoSize struct {
	FileID   string `json:"file_id"`   // 文件唯一ID
	Width    int32  `json:"width"`     // 照片宽度
	Height   int32  `json:"height"`    // 照片高度
	FileSize int32  `json:"file_size"` // 文件大小
}

// ToJSON 转换为JSON
func (photoSize *PhotoSize) ToJSON() ([]byte, error) {
	return json.Marshal(photoSize)
}

// FromJSON 从JSON反序列化
func (photoSize *PhotoSize) FromJSON(jsb []byte) error {
	return json.Unmarshal(jsb, photoSize)
}
