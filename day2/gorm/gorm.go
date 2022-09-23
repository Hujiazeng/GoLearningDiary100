package gorm

import (
	"database/sql"
	"day2/log"
	"day2/session"
)

// Gdb 用于用户层交互

type Gdb struct {
	db *sql.DB
}

func NewGdb(driver, source string) (g *Gdb, err error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Error(err)
		return
	}
	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}
	g = &Gdb{db: db}
	log.Info("Connect database success")
	return
}

// 创建数据库交互对象Session
func (g *Gdb) NewSession() *session.Session {
	return session.New(g.db)
}

// 关闭交互对象gdb
func (g *Gdb) Close() {
	if err := g.db.Close(); err != nil {
		log.Error("failed to close database")
	} else {
		log.Info("Close database success")
	}
}
