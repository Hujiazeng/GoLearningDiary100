package session

import (
	"day2/log"
	"day2/schema"
	"fmt"
	"strings"
)

func (s *Session) Model(value interface{}) *Session {
	// 避免多次解析
	if s.cacheTable == nil || value != s.cacheTable.Model {
		s.cacheTable = schema.Parse(value, s.dialect)
	}
	return s
}

func (s *Session) GetCacheTable() *schema.Schema {
	if s.cacheTable == nil {
		log.Error("no cache model")
	}
	return s.cacheTable
}

// 创建表(存在差异应该写入dialect)
func (s *Session) CreateTable() error {
	table := s.GetCacheTable()
	var columns []string
	// 拼接sql
	for _, field := range table.Fields {
		columns = append(columns, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
	}
	desc := strings.Join(columns, ",")
	_, err := s.Raw(fmt.Sprintf("CREATE TABLE %s (%s);", table.Name, desc)).Exec()
	return err
}

// 删除表
func (s *Session) DropTable() error {
	_, err := s.Raw(fmt.Sprintf("drop table if exists %s", s.GetCacheTable().Name)).Exec()
	return err
}

// 是否存在表
func (s *Session) HasTable() bool {
	// args := []interface{}{s.GetCacheTable().Name}
	row := s.Raw(fmt.Sprintf("select name from sqlite_master where type='table' and name= '%s'", s.GetCacheTable().Name)).QueryRow()
	var tem string
	_ = row.Scan(&tem)
	return tem == s.GetCacheTable().Name
}
