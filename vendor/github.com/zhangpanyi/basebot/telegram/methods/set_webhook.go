package methods

import (
	"encoding/json"
	"strconv"
)

// 设置Webhook
type setWebhook struct {
	URL            string   `json:"url"`             // 回调地址
	Certificate    []byte   `json:"certificate"`     // 证书
	MaxConnections int32    `json:"max_connections"` // 最大连接数(1-100)
	AllowedUpdates []string `json:"allowed_updates"` // 允许更新类型
}

// SetWebhook 设置webhook
func (bot *BotExt) SetWebhook(url string, certificate []byte, maxConnections int32,
	allowedUpdates []string) error {
	// 生成请求内容
	formdata := make([]Field, 0)
	formdata = append(formdata, Field{Name: "url", Text: url})
	formdata = append(formdata, Field{Name: "certificate", File: certificate, FileName: "public.pem"})
	formdata = append(formdata, Field{Name: "max_connections", Text: strconv.FormatInt(int64(maxConnections), 10)})
	if len(allowedUpdates) > 0 {
		updates, err := json.Marshal(allowedUpdates)
		if err != nil {
			return err
		}
		formdata = append(formdata, Field{Name: "allowed_updates", Text: string(updates)})
	}

	_, err := bot.Upload("setWebhook", formdata)
	return err
}
