package groupchat

import (
	"math/rand"
	"sync"
	"time"
)

// 随机器
var randx *rand.Rand

func init() {
	// 初始化随机数种子
	randx = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// Cover 封面信息
type Cover struct {
	FileName string
	FileData []byte
}

// 封面列表
var coverArray []Cover
var coverOnce sync.Once

// GetRedEnvelopesCover 获取红包封面
func GetRedEnvelopesCover() *Cover {
	if len(coverArray) == 0 {
		return nil
	}
	return &coverArray[randx.Intn(len(coverArray))]
}

// SetRedEnvelopesCoverForOnce 设置红包封面
func SetRedEnvelopesCoverForOnce(fileids []Cover) {
	coverOnce.Do(func() {
		coverArray = fileids
	})
}
