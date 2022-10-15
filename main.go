package main

import (
	"fmt"

	"github.com/ramsgoli/columnar_store/backend"
)

func main() {
	/*
		var u = backend.User{
			Name: [8]byte{'r', 'a', 'm'},
			Age:  24,
		}

		backend.Insert(&u)
	*/

	allUsers := backend.ReadUsers()
	fmt.Printf("Found %d users\n", len(*allUsers))
	firstUser := (*allUsers)[0]
	fmt.Printf("Name: %s, Age: %d\n", firstUser.Name, firstUser.Age)
}
