package main

import (
	"day7/korm"
	"day7/log"

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
	_ = k.NewSession()

}
