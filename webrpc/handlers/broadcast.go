package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/zhangpanyi/botcasino/envelopes/expiretimer"
	"github.com/zhangpanyi/botcasino/pusher"
	"github.com/zhangpanyi/botcasino/storage"
)

// Broadcast 广播消息
func Broadcast(w http.ResponseWriter, r *http.Request) {
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

	var request BroadcastRequest
	if err = json.Unmarshal(data, &request); err != nil || len(request.Message) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 广播消息
	handler := storage.SubscriberStorage{}
	subscribers, err := handler.GetSubscribers(request.BotID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for _, userID := range subscribers {
		pusher.To(expiretimer.GetBot(), userID, request.Message, true, nil)
	}

	reply := BroadcastReply{OK: true}
	jsb, err := json.Marshal(&reply)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 返回资产列表
	w.Header().Set("Server", "Casino web server")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsb)
}
