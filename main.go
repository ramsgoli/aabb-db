package main

import (
	"fmt"

	"github.com/ramsgoli/columnar_store/executor"
	"github.com/ramsgoli/columnar_store/repl"
)

func main() {
	/*
		colMetadata := []tables.ColMetadata{
			{ColName: [8]byte{'n', 'a', 'm', 'e'}, Type: 0},
			{ColName: [8]byte{'a', 'g', 'e'}, Type: 1},
		}
		newTable := tables.TableMetadata{TableName: [8]byte{'u', 's', 'e', 'r'}, NumCols: 2, ColMetadata: &colMetadata}
		err := tables.CreateTable(&newTable)
		if err != nil {
			panic(err)
		}
	*/
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
