package log

import (
	"log"
)

type Log interface {
	Println(string)
	Printf(format string, v ...interface{})
}

var logImpl Log = ConsoleLog{}

type ConsoleLog struct {
}

func (consoleLog ConsoleLog) Println(content string) {
	log.Println(content)
}

func (consoleLog ConsoleLog) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func SetImpl(a Log) {
	logImpl = a
}

func GetImpl() Log {
	return logImpl
}
