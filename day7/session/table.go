package session

import (
	"database/sql"
	"day7/schema"
	"fmt"
	"strings"
)

// 表操作

// 创建表
func (s *Session) CreateTable(schema *schema.Schema) (sql.Result, error) {
	var colomns []string
	var primarykeys []string
	var uniquekeys []string
	for i := 0; i < len(schema.Fields); i++ {
		field := schema.Fields[i]
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
	sql := fmt.Sprintf("CREATE TABLE `%s` (%s)", schema.Name, desc)

	return s.Raw(sql).Exec()
}

// 删除表
func (s *Session) DropTable(schema *schema.Schema) (sql.Result, error) {
	sql := fmt.Sprintf("DROP TABLE %s", schema.Name)
	return s.Raw(sql).Exec()
}

// 是否存在表
func (s *Session) HasTable(schema *schema.Schema) bool {
	sql := fmt.Sprintf("SHOW TABLEs like '%s'", schema.Name)
	row := s.Raw(sql).QueryRow(s.sql.String(), s.sqlVars)
	var name string
	row.Scan(&name)
	return name == schema.Name
}
