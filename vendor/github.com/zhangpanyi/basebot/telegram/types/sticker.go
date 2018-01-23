package types

import "encoding/json"

// Sticker 贴纸信息
type Sticker struct {
	FileID   string     `json:"file_id"`             // 文件ID
	Width    uint32     `json:"width"`               // 宽度
	Height   uint32     `json:"height"`              // 高度
	Thumb    *PhotoSize `json:"thumb,omitempty"`     // 缩略图
	Emoji    string     `json:"emoji,omitempty"`     // 表情
	SetName  string     `json:"set_name,omitempty"`  // 集合名字
	FileSize uint32     `json:"file_size,omitempty"` // 文件大小
}

// StickerSet 贴纸集合
type StickerSet struct {
	Name          string     `json:"name"`                     // 集合名称
	Title         string     `json:"title"`                    // 集合标题
	ContainsMasks bool       `json:"contains_masks,omitempty"` // 包含遮罩
	Stickers      []*Sticker `json:"stickers"`                 // 贴纸列表
}

// ToJSON 转换为JSON
func (set *StickerSet) ToJSON() ([]byte, error) {
	return json.Marshal(set)
}

// FromJSON 从JSON反序列化
func (set *StickerSet) FromJSON(jsb []byte) error {
	return json.Unmarshal(jsb, set)
}
