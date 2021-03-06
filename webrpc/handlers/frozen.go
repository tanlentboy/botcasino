package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/zhangpanyi/basebot/logger"
	"github.com/zhangpanyi/botcasino/storage"
)

// Frozen 冻结资产
func Frozen(w http.ResponseWriter, r *http.Request) {
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

	var request FrozenRequest
	if err = json.Unmarshal(data, &request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 冻结资产
	handler := storage.AssetStorage{}
	err = handler.FrozenAsset(request.UserID, request.Asset, request.Amount)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	logger.Errorf("Frozen user asset from webrpc, UserID: %d, Asset: %s, Amount: %d",
		request.UserID, request.Asset, request.Amount)

	reply := FrozenAssetReply{OK: true}
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
