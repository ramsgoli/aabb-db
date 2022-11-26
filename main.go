package main

import (
	"fmt"

	"github.com/ramsgoli/columnar_store/executor"
	"github.com/ramsgoli/columnar_store/repl"
)

func main() {
	command := make(chan string)
	cont := make(chan bool)
	go repl.StartRepl(command, cont, "> ")

	for t := range command {
		err := executor.Execute(t)
		if err != nil {
			panic(err)
		}
		cont <- true
	}
	fmt.Println("done!")
}
