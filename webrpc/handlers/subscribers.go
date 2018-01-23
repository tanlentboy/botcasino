package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/zhangpanyi/botcasino/storage"
)

// GetSubscribers 获取订阅者数量
func GetSubscribers(w http.ResponseWriter, r *http.Request) {
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

	var request GetSubscribersRequest
	if err = json.Unmarshal(data, &request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 查询订阅数量
	handler := storage.SubscriberStorage{}
	count, err := handler.GetSubscriberCount(request.BotID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 返回处理结果
	reply := GetSubscribersReply{Count: int32(count)}
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
