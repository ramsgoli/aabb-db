package backend

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/ramsgoli/columnar_store/config"
)

const NAME_SIZE = 8

type User struct {
	Name [NAME_SIZE]byte
	Age  uint8
}

func CreateTable(name string) {
	fmt.Printf("Using %s as the table path\n", config.GetTablePath())
}

func Insert(u *User) {
	var nameFile *os.File
	var ageFile *os.File
	var err error

	nameFile, err = os.OpenFile(config.GetTablePath()+"/user/name", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("Couldn't open name file")
	}
	defer nameFile.Close()

	_, err = nameFile.Write(u.Name[:])
	if err != nil {
		fmt.Println(err)
		panic("Couldn't write name")
	}

	ageFile, err = os.OpenFile(config.GetTablePath()+"/user/age", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		panic("Couldn't open age file")
	}
	defer ageFile.Close()

	err = binary.Write(ageFile, binary.LittleEndian, u.Age)
	if err != nil {
		fmt.Println(err)
		panic("Couldn't write age")
	}
}

func ReadUsers() *[]User {
	var users []User

	// first, read name
	nameFile, err := os.Open(config.GetTablePath() + "/user/name")
	if err != nil {
		panic("Couldn't open name file to read")
	}
	defer nameFile.Close()

	// read the first user, regardless of how many there are
	name := [NAME_SIZE]byte{}
	nameFile.Read(name[:])

	ageFile, err := os.Open(config.GetTablePath() + "/user/age")
	if err != nil {
		panic("Couldn't open age file to read")
	}

	age := [1]uint8{}
	ageFile.Read(age[:])

	users = append(users, User{Name: name, Age: age[0]})

	return &users
}
