package korm

import (
	"database/sql"
	"day7/log"
	"day7/session"
)

// 用户交互对象
type Korm struct {
	db *sql.DB
}

// 创建对象
func New(db *sql.DB) *Korm {
	return &Korm{db: db}
}

// 创建连接
func NewConnect(driver, dataSourceName string) (db *sql.DB, err error) {
	db, err = sql.Open(driver, dataSourceName)
	if err != nil {
		log.Errorf("failed to open driver")
		return
	}
	// make sure connection to the database is still alive
	if err = db.Ping(); err != nil {
		return
	}

	log.Info("connect database success...")
	return
}

// 关闭连接
func (k *Korm) Close() error {
	err := k.db.Close()
	if err != nil {
		log.Info("Close database connection error")
		return err
	}
	log.Info("Close database connection success...")
	return nil
}

// 创建会话对象
func (k *Korm) NewSession() *session.Session {
	return session.New(k.db)
}
