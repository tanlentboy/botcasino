package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/zhangpanyi/botcasino/storage"
)

// GetAssets 获取资产列表
func GetAssets(w http.ResponseWriter, r *http.Request) {
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

	var request GetAssetsRequest
	if err = json.Unmarshal(data, &request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 获取资产列表
	handler := storage.AssetStorage{}
	assets, err := handler.GetAssets(request.UserID)
	if err != nil && err != storage.ErrNoBucket {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	reply := make(GetAssetsReply, 0, len(assets))
	for _, asset := range assets {
		reply = append(reply, &Asset{
			Asset:  asset.Asset,
			Amount: asset.Amount,
			Freeze: asset.Freeze,
		})
	}

	jsb := []byte("[]")
	if len(reply) > 0 {
		jsb, err = json.Marshal(reply)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	// 返回资产列表
	w.Header().Set("Server", "Casino web server")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsb)
}
