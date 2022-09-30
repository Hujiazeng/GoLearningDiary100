package main

import (
	"day7/korm"
	"day7/log"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Animal struct {
	Name string `korm:"primarykey; type:varchar(99)"`
	Id   int    `korm:"notnull; comment:编号; type:int(12); Default:1; unique;"`
	Food string `korm:"notnull; comment:食物; type:varchar(11);"`
}

func main() {
	// Create new connect
	db, err := korm.NewConnect("mysql", "root:mysql@tcp(10.6.32.62:3306)/daqing2_dev")
	if err != nil {
		// custom log
		log.Error(err)
		return
	}
	// Create new korm obj
	k := korm.New(db)

	// Create new session obj
	session := k.NewSession()

	// Parse struct to cache schema
	session.Model(&Animal{})

	if ok := session.HasTable(); ok {
		session.DropTable()
	}
	session.CreateTable()

	cat := &Animal{Name: "cat", Id: 1, Food: "fish"}
	session.Insert(cat)

	session.Where("Name = ?", "cat").Update("Food", "mouse")

	res := []Animal{}
	session.Limit(1).Filter(&res)

	fmt.Println(res)
}
