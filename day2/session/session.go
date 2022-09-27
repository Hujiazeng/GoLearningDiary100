package session

import (
	"database/sql"
	"day2/clause"
	"day2/dialect"
	"day2/log"
	"day2/schema"
	"strings"
)

type Session struct {
	db      *sql.DB         //数据库连接对象
	sql     strings.Builder // sql语句
	sqlVars []interface{}   // 变量参数

	dialect    dialect.Dialect // 差异支持
	cacheTable *schema.Schema  // 缓存表(解析耗时避免重复解析)
	clause     clause.Clause   // 提供生成sql语句
	tx         *sql.Tx         // 事务
}

// 创建Session对象
func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{db: db, dialect: dialect}
}

// 清理sql语句及变量参数
func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
	s.clause = clause.Clause{}
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
func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if rows, err = s.db.Query(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}
