package log

import (
	"log"
	"os"
)

// 自定义日志: 颜色区分 报错行数显示

var (
	errorLog = log.New(os.Stdout, "\033[31m[error]\033[0m", log.Lshortfile|log.LstdFlags)
	infoLog  = log.New(os.Stdout, "\033[34m[info]\033[0m", log.Lshortfile|log.LstdFlags)
)

// methods
var (
	Error  = errorLog.Println
	Errorf = errorLog.Printf
	Info   = infoLog.Println
	Infof  = infoLog.Printf
)
