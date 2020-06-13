package gui

import "fmt"

type Log struct {
}

func (consoleLog Log) Println(content string) {
	Print(content + "\n")
}

func (consoleLog Log) Printf(format string, v ...interface{}) {
	Print(fmt.Sprintf(format, v...) + "\n")
}
