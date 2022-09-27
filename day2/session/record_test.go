package session

import (
	"database/sql"
	"day2/dialect"
	"os"
	"testing"
)

var (
	hero1       = &Hero{"Tom", 18}
	hero2       = &Hero{"Hu", 20}
	TestDB      *sql.DB
	TestDial, _ = dialect.GetDialect("sqlite3")
)

func TestMain(m *testing.M) {
	TestDB, _ = sql.Open("sqlite3", "dev.db")
	code := m.Run()
	_ = TestDB.Close()
	os.Exit(code)
}

func NewSession() *Session {
	return New(TestDB, TestDial)
}

func testRecordInit(t *testing.T) *Session {
	// 报错时返回上层调用者行数
	t.Helper()
	session := New(TestDB, TestDial)
	s := session.Model(&Hero{})
	s.DropTable()
	s.CreateTable()
	_, err3 := s.Insert(hero1, hero2)
	if err3 != nil {
		t.Fatal("failed to insert")
	}
	return s

}

func TestSession_Insert(t *testing.T) {
	testRecordInit(t)
}

func TestSession_Find(t *testing.T) {
	s := testRecordInit(t)
	var heros []Hero
	if err := s.Find(&heros); err != nil || len(heros) != 2 {
		t.Fatal("failed find")
	}
}

func TestSession_Limit(t *testing.T) {
	// 链式调用
	s := testRecordInit(t)
	var heros []Hero
	if err := s.Limit(1).Find(&heros); err != nil || len(heros) != 1 {
		t.Fatal("find err")
	}
}

func TestSession_Update(t *testing.T) {
	s := testRecordInit(t)
	var retMap = make(map[string]interface{})
	retMap["Age"] = 99
	affectRows, _ := s.Where("NickName = ?", "Tom").Limit(1).Update(retMap)

	HeroRecord := &Hero{}
	s.First(HeroRecord)
	if affectRows == 0 || HeroRecord.Age != 99 {
		t.Fatal("update error")
	}

}
