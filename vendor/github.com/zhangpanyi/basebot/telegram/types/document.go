package types

import (
	"encoding/json"
)

// Documnet 文档信息
type Documnet struct {
	FileID   string     `json:"file_id"`             // 文件ID
	Thumb    *PhotoSize `json:"thumb,omitempty"`     // 缩略图信息
	FileName string     `json:"file_name,omitempty"` // 文件名
	MimeType string     `json:"mime_type,omitempty"` // 文件类型
	FileSize uint32     `json:"file_size,omitempty"` // 文件大小
}

// ToJSON 转换为JSON
func (document *Documnet) ToJSON() ([]byte, error) {
	return json.Marshal(document)
}

// FromJSON 从JSON反序列化
func (document *Documnet) FromJSON(jsb []byte) error {
	return json.Unmarshal(jsb, document)
}
