package session

import (
	"database/sql"
	"day7/log"
	"strings"
)

// 会话交互对象
type Session struct {
	db *sql.DB

	sql     strings.Builder // re-use
	sqlVars []interface{}
}

// 创建对象
func New(db *sql.DB) *Session {
	return &Session{db: db, sql: strings.Builder{}, sqlVars: make([]interface{}, 0)}
}

// 注入sql(链式)
func (s *Session) Raw(sql string, vars []interface{}) *Session {
	_, err := s.sql.WriteString(sql)
	if err != nil {
		log.Error("Raw sql builder write string err")
	}
	s.sqlVars = append(s.sqlVars, vars...)
	return s
}

// 执行
func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Reset()
	log.Info("execute sql: %s, vars: %v", s.sql.String(), s.sqlVars)
	result, err = s.db.Exec(s.sql.String(), s.sqlVars...)
	return
}

// 查询单条
func (s *Session) QueryRow(sql string, vars []interface{}) *sql.Row {
	defer s.Reset()
	log.Info("query row sql: %s, vars: %v", s.sql.String(), s.sqlVars)
	return s.db.QueryRow(sql, vars...)
}

// 查询多条
func (s *Session) QueryRows(sql string, vars []interface{}) (*sql.Rows, error) {
	defer s.Reset()
	log.Info("query rows sql: %s, vars: %v", s.sql.String(), s.sqlVars)
	return s.db.Query(sql, vars...)
}

func (s *Session) Reset() {
	s.sql.Reset()
	s.sqlVars = nil
}
