package log

import (
	"log"
	"os"
)

// define new log
var (
	errorLog = log.New(os.Stdout, "\033[31m[error]\033m", log.Lshortfile|log.LstdFlags)
	infoLog  = log.New(os.Stdout, "\033[34m[info]\033m", log.Lshortfile|log.LstdFlags)
)

// export functions
var (
	Error  = errorLog.Println
	Errorf = errorLog.Printf
	Info   = infoLog.Println
	Infof  = infoLog.Printf
)
