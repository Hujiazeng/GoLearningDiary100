package korm

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func testNewConnect(t *testing.T) *sql.DB {
	conn, err := NewConnect("mysql", "dq2_developer:dc2019wsy@tcp(baef37530ed44018ba901fb7dbff5770in01.internal.cn-south-1.mysql.rds.myhuaweicloud.com)/daqing2_character_dev")
	if err != nil {
		t.Fatal("connect error")
	}
	return conn
}

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}
func testInitKrom() *Korm {
	db := testNewConnect(&testing.T{})
	return &Korm{db: db}
}

func TestNewConnect(t *testing.T) {
	_, err := NewConnect("mysql", "dq2_developer:dc2019wsy@tcp(baef37530ed44018ba901fb7dbff5770in01.internal.cn-south-1.mysql.rds.myhuaweicloud.com)/daqing2_character_dev")
	if err != nil {
		t.Fatal(err)
	}

}

func BenchmarkConnect(b *testing.B) {
	_, err := NewConnect("mysql", "dq2_developer:dc2019wsy@tcp(baef37530ed44018ba901fb7dbff5770in01.internal.cn-south-1.mysql.rds.myhuaweicloud.com)/daqing2_character_dev")
	if err != nil {
		b.Fatal(err)
	}
}

func TestClose(t *testing.T) {
	korm := testInitKrom()
	korm.Close()
	if err := korm.db.Ping(); err == nil {
		t.Fatal("close connection error")
	}
}
