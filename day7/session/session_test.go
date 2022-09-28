package session

import (
	"database/sql"
	"day7/schema"
	"reflect"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func testNewSession() (*Session, error) {
	// 获取session对象
	// db, err := sql.Open("mysql", "dq2_developer:dc2019wsy@tcp(baef37530ed44018ba901fb7dbff5770in01.internal.cn-south-1.mysql.rds.myhuaweicloud.com)/daqing2_character_dev")
	db, err := sql.Open("mysql", "root:mysql@tcp(10.6.32.62:3306)/daqing2_dev")
	if err != nil {
		return nil, err
	}
	session := New(db)
	return session, nil
}

func TestRaw(t *testing.T) {
	s, err := testNewSession()
	if err != nil {
		t.Fatal("session error")
	}

	var testList = []interface{}{}
	testList = append(testList, 1)
	s.Raw("SELECT * FROM g", testList)
	if s.sql.String() != "SELECT * FROM g" {
		t.Fatal("raw sql error")
	}
	t.Log(s.sqlVars)
	if ok := reflect.DeepEqual(testList, s.sqlVars); !ok {
		t.Fatal("raw sqlVars error")
	}
}

type Hero struct {
	NickName string `korm:"primarykey;  type:varchar(99)"`
	Age      int    `korm:"notnull; comment:年龄; type:int(12); Default:520; unique;"`
}

func TestCreateTable(t *testing.T) {
	s, err := testNewSession()
	if err != nil {
		t.Fatal("session error")
	}
	schema := schema.Parse(&Hero{})
	s.CreateTable(schema)
}

func TestDropTable(t *testing.T) {
	s, err := testNewSession()
	if err != nil {
		t.Fatal("session error")
	}
	schema := schema.Parse(&Hero{})
	ok := s.HasTable(schema)
	if ok {
		_, err = s.DropTable(schema)
		if err != nil {
			t.Fatal("drop table err")
		}
	}
}
