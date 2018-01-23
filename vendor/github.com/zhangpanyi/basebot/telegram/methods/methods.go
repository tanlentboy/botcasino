package methods

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/zhangpanyi/basebot/telegram/types"
)

// Field 字段信息
type Field struct {
	Name     string // 字段名
	Text     string // 文本
	File     []byte // 文件
	FileName string // 文件名
}

// IsFile 是否文件
func (f *Field) IsFile() bool {
	return len(f.FileName) > 0 && f.File != nil
}

// Bot 机器人信息
type Bot struct {
	ID        int64  `json:"id"`         // 机器人唯一ID
	FirstName string `json:"first_name"` // 机器人First name
	LastName  string `json:"last_name"`  // 机器人Last name
	UserName  string `json:"username"`   // 机器人用户名
}

// BotExt 机器人扩展信息
type BotExt struct {
	Bot
	Token      string // 机器人Token
	APIWebsite string // 机器人API网站
}

// ErrorResonpe 错误响应
type ErrorResonpe struct {
	OK          bool   `json:"ok"`          // 是否成功
	ErrorCode   int32  `json:"error_code"`  // 错误代码
	Description string `json:"description"` // 描述信息
}

// GetMeResonpe 获取信息响应
type GetMeResonpe struct {
	OK     bool `json:"ok"`     // 是否成功
	Result *Bot `json:"result"` // 机器人信息
}

// SendMessageResonpe 发送消息响应
type SendMessageResonpe struct {
	OK     bool           `json:"ok"`     // 是否成功
	Result *types.Message `json:"result"` // 消息内容
}

// 格式化响应
func (bot *BotExt) formatRespone(res *http.Response) ([]byte, error) {
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		respone := ErrorResonpe{}
		err = json.Unmarshal(body, &respone)
		if err != nil {
			reason := fmt.Sprintf("http status code: %v, %v", res.StatusCode, string(body))
			return body, errors.New(reason)
		}
		reason := fmt.Sprintf("error code: %v, description: %v", respone.ErrorCode, respone.Description)
		return body, errors.New(reason)
	}

	return body, nil
}

// Call 调用方法
func (bot *BotExt) Call(methodName string, request interface{}) ([]byte, error) {
	jsb, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(jsb)
	a := [...]string{bot.APIWebsite, "bot", bot.Token, "/", methodName}
	res, err := http.Post(strings.Join(a[:], ""), "application/json", reader)
	if err != nil {
		return nil, err
	}
	return bot.formatRespone(res)
}

// Upload 上传文件
func (bot *BotExt) Upload(methodName string, formdata []Field) ([]byte, error) {
	// 构造数据
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	for _, field := range formdata {
		if !field.IsFile() {
			// 写入文本
			fw, err := w.CreateFormField(field.Name)
			if err != nil {
				return nil, err
			}
			if _, err = fw.Write([]byte(field.Text)); err != nil {
				return nil, err
			}
		} else {
			// 写入文件
			fw, err := w.CreateFormFile(field.Name, field.FileName)
			if err != nil {
				return nil, err
			}
			if _, err := io.Copy(fw, bytes.NewReader(field.File)); err != nil {
				return nil, err
			}
		}
	}
	w.Close()

	// 执行请求
	a := [...]string{bot.APIWebsite, "bot", bot.Token, "/", methodName}
	req, err := http.NewRequest("POST", strings.Join(a[:], ""), &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return bot.formatRespone(res)
}
