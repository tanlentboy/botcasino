package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/zhangpanyi/botcasino/envelopes/privatechat"
	"github.com/zhangpanyi/botcasino/withdraw"

	"github.com/zhangpanyi/basebot/logger"
)

// RestoreOrder 恢复订单
func RestoreOrder(w http.ResponseWriter, r *http.Request) {
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

	var request RestoreOrderRequest
	if err = json.Unmarshal(data, &request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 执行恢复订单
	future, err := withdraw.RestoreFuture(request.OrderID)
	logger.Errorf("Restore order from webrpc, OrderID: %d, %v",
		request.OrderID, err)

	// 写入操作记录
	handler := privatechat.Withdraw{}
	handler.HandleWithdrawFuture(future)

	// 返回处理结果
	reply := RestoreOrderReply{OK: true}
	if err != nil {
		reason := err.Error()
		reply.Reason = &reason
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
