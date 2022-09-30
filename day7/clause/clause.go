package clause

// 从句
type Clause struct {
	sqlMap  map[Type]string
	varsMap map[Type][]interface{}
}

type Type int

const (
	SELECT Type = iota
	WHERE
	ORDER
	LIMIT
	INSERT
	DELETE
	VALUES
	UPDATE
)

// 写入子句
func (c *Clause) Set(Name Type, values ...interface{}) {
	if c.sqlMap == nil {
		c.sqlMap = make(map[Type]string)
		c.varsMap = make(map[Type][]interface{})
	}
	sql, vars := generators[Name](values...)
	c.sqlMap[Name] = sql
	c.varsMap[Name] = vars
}

// 组装
func (c *Clause) Build(Order ...Type) (s string, vars []interface{}) {
	defer c.Clean()
	for i := 0; i < len(Order); i++ {
		s += c.sqlMap[Order[i]]
		vars = append(vars, c.varsMap[Order[i]]...)
	}
	return
}

// 清空
func (c *Clause) Clean() {
	c.sqlMap = nil
	c.varsMap = nil
}
