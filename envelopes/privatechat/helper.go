package privatechat

import (
	"github.com/zhangpanyi/botcasino/config"
)

// 语言翻译
func tr(userID int64, key string) string {
	return config.GetLanguge().Value("zh_CN", key)
}
