package session

import (
	"database/sql"
	"fmt"
	"strings"
)

// 表操作

// 创建表
func (s *Session) CreateTable() (sql.Result, error) {
	var colomns []string
	var primarykeys []string
	var uniquekeys []string
	for i := 0; i < len(s.cacheSchema.Fields); i++ {
		field := s.cacheSchema.Fields[i]
		tempColomn := fmt.Sprintf("`%s` %s %s ", field.Name, field.Type, field.UnSign)
		if field.Default != "" {
			tempColomn = fmt.Sprintf("%s Default %s", tempColomn, field.Default)
		}
		if field.NotNull != "" {
			tempColomn = fmt.Sprintf("%s NOT NULL", tempColomn)
		}
		if field.Comment != "" {
			tempColomn = fmt.Sprintf("%s COMMENT '%s'", tempColomn, field.Comment)
		}
		if field.Primary {
			primarykeys = append(primarykeys, fmt.Sprintf("PRIMARY KEY (`%s`) USING BTREE", field.Name))
		}
		if field.Unique {
			uniquekeys = append(uniquekeys, fmt.Sprintf("UNIQUE KEY `%s` (`%s`) USING BTREE", "uk_"+strings.ToLower(field.Name), field.Name))
		}
		colomns = append(colomns, tempColomn)
	}
	colomns = append(colomns, primarykeys...)
	colomns = append(colomns, uniquekeys...)
	desc := strings.Join(colomns, ",")
	sql := fmt.Sprintf("CREATE TABLE `%s` (%s)", s.cacheSchema.Name, desc)

	return s.Raw(sql).Exec()
}

// 删除表
func (s *Session) DropTable() (sql.Result, error) {
	sql := fmt.Sprintf("DROP TABLE %s", s.cacheSchema.Name)
	return s.Raw(sql).Exec()
}

// 是否存在表
func (s *Session) HasTable() bool {
	sql := fmt.Sprintf("SHOW TABLEs like '%s'", s.cacheSchema.Name)
	row := s.Raw(sql).QueryRow()
	var name string
	row.Scan(&name)
	return name == s.cacheSchema.Name
}
