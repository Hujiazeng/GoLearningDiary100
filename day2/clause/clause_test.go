package clause

import (
	"testing"
)

func TestClause_Build(t *testing.T) {
	var clause Clause
	clause.Set(LIMIT, 1)
	clause.Set(SELECT, "Hero", []string{"*"})
	clause.Set(WHERE, "Age = ?", 1)

	sql, vars := clause.Build(SELECT, WHERE, LIMIT)
	t.Log(sql)
	if sql != "SELECT * FROM Hero WHERE Age = ? LIMIT ?" {
		t.Fatal("failed to build sql")
	}
	t.Log(vars)
	// if !reflect.DeepEqual(vars, []interface{}{"Age", 1})

}
