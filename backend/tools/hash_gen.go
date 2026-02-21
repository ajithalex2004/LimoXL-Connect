package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	bytes, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	fmt.Println(string(bytes))
}
