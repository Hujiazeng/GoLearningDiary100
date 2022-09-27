package clause

import (
	"fmt"
	"strings"
)

// 生成子句

// 定义方法(返回sql和参数)
type generator func(values ...interface{}) (string, []interface{})

// 定义字典映射 操作枚举:方法
var generators map[Type]generator

func init() {
	generators = make(map[Type]generator)
	generators[INSERT] = _insert
	generators[WHERE] = _where
	generators[SELECT] = _select
	generators[LIMIT] = _limit
	generators[VALUES] = _values
	generators[UPDATE] = _update
	generators[DELETE] = _delete
	generators[COUNT] = _count
}

// 生成SQL语句后面需要跟的问号数量
func genBindVars(num int) string {
	var vars []string
	for i := 0; i < num; i++ {
		vars = append(vars, "?")
	}
	return strings.Join(vars, ", ")
}

func _insert(values ...interface{}) (string, []interface{}) {
	tableName := values[0]
	// .([]string) 断言是字符串列表
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("INSERT INTO %s (%s)", tableName, fields), []interface{}{}
}

// 传递一个列表, 第一个值为表名
func _select(values ...interface{}) (string, []interface{}) {
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("SELECT %s FROM %s", fields, tableName), []interface{}{}
}

func _limit(values ...interface{}) (string, []interface{}) {
	return "LIMIT ?", values
}

func _where(values ...interface{}) (string, []interface{}) {
	desc, vars := values[0], values[1:]
	return fmt.Sprintf("WHERE %s", desc), vars
}

func _values(values ...interface{}) (string, []interface{}) {
	var bindStr string
	var sql strings.Builder
	var vars []interface{}
	sql.WriteString("VALUES ")
	for i, value := range values {
		v := value.([]interface{})
		// 同一张表仅首次需要生成问号数量
		if bindStr == "" {
			bindStr = genBindVars(len(v))
		}
		// "insert into table values(?,?,?,?)";
		sql.WriteString(fmt.Sprintf("(%v)", bindStr))
		// 不是最后一个都需要加逗号
		if i != len(values)-1 {
			sql.WriteString(", ")
		}
		vars = append(vars, v...)
	}

	return sql.String(), vars
}

func _update(values ...interface{}) (string, []interface{}) {
	tableName := values[0]
	// 参数格式: {字段名: 值, 字段名2:值}
	m := values[1].(map[string]interface{})
	var keys []string
	var vars []interface{}
	for k, v := range m {
		keys = append(keys, k+" = ?")
		vars = append(vars, v)
	}
	return fmt.Sprintf("UPDATE %s SET %s", tableName, strings.Join(keys, ", ")), vars
}

func _delete(values ...interface{}) (string, []interface{}) {
	return fmt.Sprintf("DELETE FROM %s", values[0]), []interface{}{}
}

func _count(values ...interface{}) (string, []interface{}) {
	return _select(values[0], []string{"count(*)"})
}
