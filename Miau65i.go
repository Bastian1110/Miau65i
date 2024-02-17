package main

import (
	"fmt"
	"strings"

	"github.com/Bastian1110/Miau65i/lib"
)

func main() {
	fmt.Println("Welcome to Miau65i! 🧶")
	lib.Lexer()
	lines := lib.ReadFile("./test.miau")
	fullSourceCode := strings.Join(lines, "\n")
	tokens := lib.Tokenize(fullSourceCode)
	fmt.Println("Tokens : ", tokens)
	parser := lib.NewParser(tokens)
	ast := parser.ParseProgram()
	fmt.Println("AST :", ast)
	ast.Show(0)
	sa := lib.NewSymbolTable()
	ast.AnalyzeSemantics(sa)

}
