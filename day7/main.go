package main

import (
	"day7/korm"
	"day7/log"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := korm.NewConnect("mysql", "root:mysql@tcp(10.6.32.62:3306)/daqing2_dev")
	if err != nil {
		log.Error(err)
		return
	}
	// 用户交互对象
	k := korm.New(db)

	// 会话对象
	session := k.NewSession()

	res, err := session.Raw("INSERT INTO `user` (Name, Age) VALUES (?,?)", "hh", 123).Exec()
	fmt.Println(res)
}
