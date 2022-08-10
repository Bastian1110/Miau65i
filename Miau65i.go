package main

import (
	"fmt"

	"github.com/Bastian1110/Miau65i/lib"
)

func main() {
	fmt.Println("Welcome to Miau65i! 🧶")
	lib.Lexer()
	lib.ReadFile("./test.miau")

}
