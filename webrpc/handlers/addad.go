package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/zhangpanyi/botcasino/storage"
)

// AddAd 增加广告
func AddAd(w http.ResponseWriter, r *http.Request) {
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

	var request AddAdRequest
	if err = json.Unmarshal(data, &request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 添加广告
	handler := storage.AdStorage{}
	id, err := handler.AddAd(request.BotID, request.AD)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 回复处理结果
	reply := AddAdReply{OK: true, ID: id}
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
