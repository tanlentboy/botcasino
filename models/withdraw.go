package models

import (
	"time"

	"upper.io/db.v3"
)

// WithdrawTableName 数据库表名
const WithdrawTableName = "withdraw"

// WithdrawStatus 提现状态
type WithdrawStatus int

const (
	_ WithdrawStatus = iota
	// WithdrawStatusWaiting 等待提现
	WithdrawStatusWaiting
	// WithdrawStatusProcessing 正在提现
	WithdrawStatusProcessing
	// WithdrawStatusSuccessful 提现成功
	WithdrawStatusSuccessful
	// WithdrawStatusFailure 提现失败
	WithdrawStatusFailure
)

// Withdraw 提现记录
type Withdraw struct {
	OrderID    int64          `db:"id"`               // 订单ID
	UserID     int64          `db:"user_id"`          // 用户ID
	To         string         `db:"to"`               // 帐户名
	AssetID    string         `db:"asset_id"`         // 资产ID
	Amount     uint32         `db:"amount"`           // 资产金额(分)
	Fee        uint32         `db:"fee"`              // 手续费用(分)
	Real       uint32         `db:"real"`             // 真实手续费(分)
	Status     WithdrawStatus `db:"status"`           // 提现状态
	Reason     *string        `db:"reason,omitempty"` // 错误原因
	InsertedAt time.Time      `db:"inserted_at"`      // 插入日期
}

// GetWithdrawOrder 获取提现订单
func GetWithdrawOrder(orderID int64) (*Withdraw, error) {
	if pools == nil {
		return nil, db.ErrNotConnected
	}

	var withdraw Withdraw
	query := (*pools).SelectFrom(WithdrawTableName).Where("id = ?", orderID)
	if err := query.One(&withdraw); err != nil {
		return nil, err
	}
	return &withdraw, nil
}

// InsertWithdraw 插入提现记录
func InsertWithdraw(userID int64, to, assetID string, amount, fee uint32) (int64, error) {
	if pools == nil {
		return 0, db.ErrNotConnected
	}

	col := (*pools).InsertInto(WithdrawTableName).Columns("user_id", "to",
		"asset_id", "amount", "fee", "status")
	q := col.Values(userID, to, assetID, amount, fee, WithdrawStatusWaiting)
	res, err := q.Exec()
	if err != nil {
		return 0, err
	}

	insertID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return insertID, nil
}

// UpdateWithdrawFee 更新手续费
func UpdateWithdrawFee(orderID int64, fee uint32) error {
	if pools == nil {
		return db.ErrNotConnected
	}

	query := (*pools).Update(WithdrawTableName).Set("real", fee)
	_, err := query.Where("id = ?", orderID).Exec()
	return err
}

// UpdateWithdraw 更新提现记录
func UpdateWithdraw(orderID int64, status WithdrawStatus, reason *string) error {
	if pools == nil {
		return db.ErrNotConnected
	}

	query := (*pools).Update(WithdrawTableName).Set("status", status).Set("reason", reason)
	_, err := query.Where("id = ?", orderID).Exec()
	return err
}
