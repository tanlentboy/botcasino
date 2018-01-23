package models

import (
	"time"

	db "upper.io/db.v3"
)

// HistoryTableName 数据库表名
const HistoryTableName = "history"

// History 历史记录
type History struct {
	UserID     int64     `db:"user_id"`     // 用户ID
	Describe   string    `db:"describe"`    // 描述信息
	InsertedAt time.Time `db:"inserted_at"` // 插入日期
}

// InsertHistory 插入历史
func InsertHistory(userID int64, describe string) error {
	if pools == nil {
		return db.ErrNotConnected
	}

	col := (*pools).InsertInto(HistoryTableName).Columns("user_id", "describe")
	q := col.Values(userID, describe)
	_, err := q.Exec()
	return err
}

// GetUserHistory 获取用户历史记录
func GetUserHistory(userID int64, page int, number int) ([]History, uint, error) {
	if pools == nil {
		return nil, 0, db.ErrNotConnected
	}

	cond := db.Cond{"user_id": userID}

	// 查询记录数量
	count, err := (*pools).Collection(HistoryTableName).Find(cond).Count()
	if err != nil {
		return nil, 0, err
	}

	// 计算总页数
	pagesum := int(count) / number
	if int(count)%number != 0 {
		pagesum++
	}
	if page > pagesum {
		page = pagesum
	}

	// 查询指定历史
	var history []History
	q := (*pools).SelectFrom(HistoryTableName).Where(cond)
	err = q.OrderBy("id DESC").Offset((page - 1) * number).Limit(number).All(&history)
	if err != nil {
		return nil, 0, err
	}

	return history, uint(pagesum), nil
}
