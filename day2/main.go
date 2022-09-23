package main

import (
	"day2/gorm"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	gdb, _ := gorm.NewGdb("sqlite3", "dev.db")
	defer gdb.Close()
	s := gdb.NewSession()
	s.Raw("drop table if exists User;").Exec()
	s.Raw("create table User(Name text);").Exec()
	s.Raw("create table User(Name text);").Exec()
	s.Raw("insert into User(`Name`) values (?), (?)", "Tom", "Hu").Exec()
	ret := s.Raw("select Name from User;").QueryRow()
	var name string
	ret.Scan(&name)
	fmt.Println(name)

}
