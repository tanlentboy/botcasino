package handlers

import (
	"net/http"
	"strings"

	"github.com/zhangpanyi/botcasino/config"

	"github.com/zhangpanyi/basebot/logger"
)

// 身份验证
func authentication(r *http.Request) bool {
	result := strings.Split(r.RemoteAddr, ":")
	if len(result) != 2 {
		return false
	}

	dynamicCfg := config.GetDynamic()
	logger.Infof("Remote call rpc, %s", result[0])
	for _, addr := range dynamicCfg.WhiteList {
		if result[0] == addr {
			return true
		}
	}
	return false
}
