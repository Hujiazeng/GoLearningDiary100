package clause

import "testing"

func TestSet(t *testing.T) {
	c := &Clause{}
	c.Set(SELECT, "Hero", []string{"*"})
	if c.sqlMap[SELECT] != "SELECT (*) FROM Hero " {
		t.Fatal("set select error")
	}
}
