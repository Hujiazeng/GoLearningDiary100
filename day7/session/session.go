package session

import (
	"database/sql"
	"day7/clause"
	"day7/log"
	"day7/schema"
	"reflect"
	"strings"
)

// 会话交互对象
type Session struct {
	db *sql.DB

	sql         strings.Builder // re-use
	sqlVars     []interface{}
	clause      *clause.Clause
	cacheSchema *schema.Schema
	tx          *sql.Tx //transaction
}

type CommenDB interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

var _ CommenDB = (*sql.DB)(nil)
var _ CommenDB = (*sql.Tx)(nil)

func (s *Session) DB() CommenDB {
	if s.tx == nil {
		return s.db
	}
	return s.tx
}

// 创建对象
func New(db *sql.DB) *Session {
	return &Session{
		db:      db,
		clause:  &clause.Clause{},
		sql:     strings.Builder{},
		sqlVars: make([]interface{}, 0)}
}

// 注入sql(链式)
func (s *Session) Raw(sql string, vars ...interface{}) *Session {
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
	log.Infof("execute sql: %s, vars: %v", s.sql.String(), s.sqlVars)
	result, err = s.DB().Exec(s.sql.String(), s.sqlVars...)
	return
}

// 查询单条
func (s *Session) QueryRow() *sql.Row {
	defer s.Reset()
	log.Infof("query row sql: %s, vars: %v", s.sql.String(), s.sqlVars)
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...)
}

// 查询多条
func (s *Session) QueryRows() (*sql.Rows, error) {
	defer s.Reset()
	log.Infof("query rows sql: %s, vars: %v", s.sql.String(), s.sqlVars)
	return s.DB().Query(s.sql.String(), s.sqlVars...)
}

// 避免重复解析
func (s *Session) Model(model interface{}) *Session {
	if s.cacheSchema == nil || reflect.ValueOf(s.cacheSchema.Model) != reflect.Indirect(reflect.ValueOf(model)) {
		s.cacheSchema = schema.Parse(model)
	}
	return s
}

func (s *Session) Reset() {
	s.sql.Reset()
	s.sqlVars = nil
}
