package main

import (
	"fmt"

	"github.com/Bastian1110/Miau65i/lib"
)

func main() {
	fmt.Println("Welcome to Miau65i! 🧶")
	lib.Lexer()
	lines := lib.ReadFile("./test.miau")
	for _, line := range lines {
		tokens := lib.Tokenize(line)
		fmt.Println("Tokens : ", tokens)
		parser := lib.NewParser(tokens)
		ast := parser.ParseProgram()
		fmt.Println("AST :", ast)
		ast.Show(0)
	}

}
