package session

import (
	"database/sql"
	"day2/dialect"
	"fmt"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

type Hero struct {
	NickName string `gorm:"PRIMARY KEY"`
	Age      int
}

func (h *Hero) AfterFind(s *Session) error {
	fmt.Println("After Find Hook...")
	return nil
}

func TestCreateTable(t *testing.T) {
	TestDB, _ := sql.Open("sqlite3", "dev.db")
	if TestDB == nil {
		t.Fatal("no driver")
	}
	if err := TestDB.Ping(); err != nil {
		t.Fatal("connect error")
	}
	TestDial, _ := dialect.GetDialect("sqlite3")
	session := New(TestDB, TestDial)
	session.Model(&Hero{})
	session.CreateTable()
	if ok := session.HasTable(); !ok {
		t.Fatal("failed to create table")
	}

}
