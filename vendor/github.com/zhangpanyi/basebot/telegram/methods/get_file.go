package methods

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

// 获取文件
type getFile struct {
	FileID string `json:"file_id"` // 文件ID
}

// 文件信息
type fileInfo struct {
	FileID   string `json:"file_id"`   // 文件ID
	FileSize uint32 `json:"file_size"` // 文件大小
	FilePath string `json:"file_path"` // 文件路径
}

// 获取文件响应
type getFileResonpe struct {
	OK     bool      `json:"ok"`     // 是否成功
	Result *fileInfo `json:"result"` // 文件信息
}

// GetFile 获取文件
func (bot *BotExt) GetFile(fileID string) ([]byte, error) {
	// 获取文件路径
	request := getFile{
		FileID: fileID,
	}
	data, err := bot.Call("getFile", &request)
	if err != nil {
		return nil, err
	}

	res := getFileResonpe{}
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}

	// 获取文件数据
	a := [...]string{bot.APIWebsite, "file/bot", bot.Token, "/", res.Result.FilePath}
	respone, err := http.Get(strings.Join(a[:], ""))
	if err != nil {
		return nil, err
	}
	defer respone.Body.Close()

	filedata, err := ioutil.ReadAll(respone.Body)
	if err != nil {
		return nil, err
	}
	return filedata, nil
}
