package telegram

import (
	"strings"
)

// Pre 预格式化文本
// 参考：https://core.telegram.org/bots/api#markdown-style
func Pre(s string) string {
	s = strings.Replace(s, "*", "\\*", -1)
	s = strings.Replace(s, "_", "\\_", -1)
	s = strings.Replace(s, "`", "\\`", -1)
	return s
}
