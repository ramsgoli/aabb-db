package repl

import (
	"bufio"
	"fmt"
	"os"
)

func printPrompt(p string) {
	fmt.Print(p)
}

func printInputs(inputs []string) {
	for _, v := range inputs {
		fmt.Printf("received: %s\n", v)
	}
}

func StartRepl(c chan string, con chan bool, prompt string) {
	reader := bufio.NewScanner(os.Stdin)

	printPrompt(prompt)
	for reader.Scan() {
		t := reader.Text()
		if t == ".exit" {
			close(c)
			break
		}

		c <- t

		// wait until we can contnue
		<-con
		printPrompt(prompt)
	}
}
