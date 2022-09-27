package clause

import "strings"

type Type int

const (
	INSERT Type = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBY
	UPDATE
	DELETE
	COUNT
)

type Clause struct {
	sql     map[Type]string
	sqlVars map[Type][]interface{}
}

// 根据传入的type调用方法生成sql和vars, 并写入clause中
func (c *Clause) Set(name Type, vars ...interface{}) {
	if c.sql == nil {
		c.sql = make(map[Type]string)
		c.sqlVars = make(map[Type][]interface{})
	}
	sql, vars := generators[name](vars...)
	c.sql[name] = sql
	c.sqlVars[name] = vars
}

// 生成总sql和变量列表
func (c *Clause) Build(orders ...Type) (string, []interface{}) {
	var sqls []string
	var vars []interface{}
	for _, name := range orders {
		if sql, ok := c.sql[name]; ok {
			sqls = append(sqls, sql)
			vars = append(vars, c.sqlVars[name]...)
		}
	}
	return strings.Join(sqls, " "), vars
}
