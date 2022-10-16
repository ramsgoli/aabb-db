package repl

import (
	"bufio"
	"fmt"
	"os"
)

func printPrompt() {
	fmt.Print("> ")
}

func printInputs(inputs []string) {
	for _, v := range inputs {
		fmt.Printf("received: %s\n", v)
	}
}

func StartRepl(c chan string, con chan bool) {
	reader := bufio.NewScanner(os.Stdin)

	printPrompt()
	for reader.Scan() {
		t := reader.Text()
		if t == ".exit" {
			close(c)
			break
		}

		c <- t

		// wait until we can contnue
		<-con
		printPrompt()
	}
}
