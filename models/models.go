package models

import (

	// examples
	// https://upper.io/db.v3/examples
	"github.com/zhangpanyi/basebot/logger"

	db "upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
	"upper.io/db.v3/mysql"
)

// 数据库连接池
var pools *sqlbuilder.Database

// Connect 连接到MySQL
func Connect(settings db.ConnectionURL, conns int) error {
	if pools != nil {
		return nil
	}

	db, err := mysql.Open(settings)
	if err != nil {
		return err
	}

	db.SetLogging(false)
	db.SetMaxOpenConns(conns)
	db.SetMaxIdleConns(conns)

	if err = db.Ping(); err != nil {
		db.Close()
		logger.Errorf("Failed to connect to mysql, %v", err)
	}
	pools = &db
	return nil
}
