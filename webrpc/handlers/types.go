package handlers

import (
	"github.com/zhangpanyi/botcasino/storage"
)

// Asset 资产信息
type Asset struct {
	Asset  string `json:"asset"`  // 资产名称
	Amount uint32 `json:"amount"` // 资产总额
	Freeze uint32 `json:"freeze"` // 冻结资产
}

// AddAdRequest 添加广告请求
type AddAdRequest struct {
	BotID int64  `json:"bot_id"` // 机器人ID
	AD    string `json:"ad"`     // 广告内容
}

// AddAdReply 添加广告回复
type AddAdReply struct {
	OK bool   `json:"ok"`              // 是否成功
	ID uint32 `json:"ad_id,omitempty"` // 广告ID
}

// DelAdRequest 删除广告请求
type DelAdRequest struct {
	BotID int64  `json:"bot_id"` // 机器人ID
	ID    uint32 `json:"ad_id"`  // 广告ID
}

// DelAdReply 删除广告回复
type DelAdReply struct {
	OK bool `json:"ok"` // 是否成功
}

// GetAdsRequest 获取广告列表请求
type GetAdsRequest struct {
	BotID int64 `json:"bot_id"` // 机器人ID
}

// GetAdsReply 获取广告列表回复
type GetAdsReply []*storage.Ad

// AddAssetRequest 增加资产请求
type AddAssetRequest struct {
	UserID int64  `json:"user_id"` // 用户ID
	Asset  string `json:"asset"`   // 资产类型
	Amount uint32 `json:"amount"`  // 资产数量
}

// AddAssetReply 增加资产回复
type AddAssetReply struct {
	OK bool `json:"ok"` // 是否成功
}

// BroadcastRequest 广播消息请求
type BroadcastRequest struct {
	BotID   int64  `json:"bot_id"`  // 机器人ID
	Message string `json:"message"` // 消息内容
}

// BroadcastReply 广播消息回复
type BroadcastReply struct {
	OK bool `json:"ok"` // 是否成功
}

// DeductAssetRequest 扣除资产请求
type DeductAssetRequest struct {
	UserID int64  `json:"user_id"` // 用户ID
	Asset  string `json:"asset"`   // 资产类型
	Amount uint32 `json:"amount"`  // 资产数量
}

// DeductAssetReply 扣除资产回复
type DeductAssetReply struct {
	OK bool `json:"ok"` // 是否成功
}

// FrozenRequest 冻结资产请求
type FrozenRequest struct {
	UserID int64  `json:"user_id"` // 用户ID
	Asset  string `json:"asset"`   // 资产类型
	Amount uint32 `json:"amount"`  // 资产数量
}

// FrozenAssetReply 冻结资产回复
type FrozenAssetReply struct {
	OK bool `json:"ok"` // 是否成功
}

// UnfrozenRequest 解冻资产请求
type UnfrozenRequest struct {
	UserID int64  `json:"user_id"` // 用户ID
	Asset  string `json:"asset"`   // 资产类型
	Amount uint32 `json:"amount"`  // 资产数量
}

// UnfrozenAssetReply 解冻资产回复
type UnfrozenAssetReply struct {
	OK bool `json:"ok"` // 是否成功
}

// GetAssetsRequest 获取资产列表请求
type GetAssetsRequest struct {
	UserID int64 `json:"user_id"` // 用户ID
}

// GetAssetsReply 获取资产列表回复
type GetAssetsReply []*Asset

// RestoreOrderRequest 恢复订单请求
type RestoreOrderRequest struct {
	OrderID int64 `json:"order_id"` // 订单ID
}

// RestoreOrderReply 恢复订单回复
type RestoreOrderReply struct {
	OK     bool    `json:"ok"`               // 是否成功
	Reason *string `json:"reason,omitempty"` // 失败原因
}

// GetSubscribersRequest 获取订阅数量请求
type GetSubscribersRequest struct {
	BotID int64 `json:"bot_id"` // 机器人ID
}

// GetSubscribersReply 获取订阅数量回复
type GetSubscribersReply struct {
	Count int32 `json:"count"` // 订阅者数量
}
