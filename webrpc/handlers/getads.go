package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/zhangpanyi/botcasino/storage"
)

// GetAds 获取广告列表
func GetAds(w http.ResponseWriter, r *http.Request) {
	// 验证权限
	if !authentication(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// 解析请求参数
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var request GetAdsRequest
	if err = json.Unmarshal(data, &request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 获取广告列表
	handler := storage.AdStorage{}
	ads, err := handler.GetAds(request.BotID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 回复处理结果
	reply := make(GetAdsReply, 0, len(ads))
	for _, item := range ads {
		reply = append(reply, item)
	}
	jsb, err := json.Marshal(&reply)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Server", "Casino web server")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsb)
}
