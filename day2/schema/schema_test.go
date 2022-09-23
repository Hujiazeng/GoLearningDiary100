package schema

import (
	"day2/dialect"
	"testing"
)

type DD struct {
	Name string `gorm:"PRIMARY KEY"`
	Age  int
}

var sqlite3, _ = dialect.GetDialect("sqlite3")

func TestParse(t *testing.T) {
	schema := Parse(&DD{}, sqlite3)
	if schema.Name != "DD" || len(schema.Fields) != 2 {
		t.Fatal("failed to parse")
	}

	if schema.GetField("Name").Tag != "PRIMARY KEY" {
		t.Fatal("failed to parse")
	}
}
