package task

import "fmt"

type Example struct{}

func (j *Example) Run() {
	fmt.Println("testFunc") // 每天打印一遍
}
