package clause

import (
	"fmt"
	"strings"
)

type generator func(values ...interface{}) (sql string, vars []interface{})

// 子句生成器
var generators map[Type]generator

func init() {
	generators = make(map[Type]generator, 0)
	generators[SELECT] = _select
	generators[INSERT] = _insert
	generators[VALUES] = _values
	generators[ORDER] = _order
	generators[LIMIT] = _limit
	generators[WHERE] = _where
	generators[UPDATE] = _update
}

func genUnknowStr(num int) string {
	temp := []string{}
	for i := 0; i < num; i++ {
		temp = append(temp, "?")
	}
	return strings.Join(temp, ",")
}

func _select(values ...interface{}) (string, []interface{}) {
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	sql := fmt.Sprintf("SELECT %s FROM `%s` ", fields, tableName)
	return sql, []interface{}{}
}

func _insert(values ...interface{}) (string, []interface{}) {
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("INSERT INTO `%s` (%s)", tableName, fields), []interface{}{}
}

func _values(values ...interface{}) (string, []interface{}) {
	unknowStr := genUnknowStr(len(values))
	return fmt.Sprintf(" VALUES (%s)", unknowStr), values
}

func _limit(values ...interface{}) (string, []interface{}) {
	num := values[0].(int)
	return fmt.Sprintf(" LIMIT %d", num), []interface{}{}
}

func _where(values ...interface{}) (string, []interface{}) {
	desc, args := values[0], values[1:]
	return fmt.Sprintf(" WHERE %s", desc), args
}

func _order(values ...interface{}) (string, []interface{}) {
	var temp []string
	for i := 0; i < len(values); i++ {
		temp = append(temp, values[0].(string))
	}
	return fmt.Sprintf(" ORDER BY %s", strings.Join(temp, ",")), []interface{}{}
}

func _update(values ...interface{}) (string, []interface{}) {
	tableName := values[0]
	mapping := values[1].(map[string]interface{})
	var temp []string
	var vars []interface{}
	for k, v := range mapping {
		temp = append(temp, k+" = ?")
		vars = append(vars, v)
	}
	return fmt.Sprintf("UPDATE %s SET %s", tableName, strings.Join(temp, ",")), vars
}
