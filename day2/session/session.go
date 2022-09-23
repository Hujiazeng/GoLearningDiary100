package session

import (
	"database/sql"
	"day2/log"
	"strings"
)

type Session struct {
	db      *sql.DB         //数据库连接对象
	sql     strings.Builder // sql语句
	sqlVars []interface{}   // 变量参数
}

// 创建Session对象
func New(db *sql.DB) *Session {
	return &Session{db: db}
}

// 清理sql语句及变量参数
func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
}

// 封装写sql语句
func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, values...)
	return s
}

// 封装执行SQL语句(退出清空session保存的sql, 日志打印执行语句sql)
func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)

	if result, err = s.db.Exec(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}

// 封装查询语句
func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	return s.db.QueryRow(s.sql.String(), s.sqlVars...)
}
