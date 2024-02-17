package lib

import (
	"fmt"
	"os"
)

type SymbolTable struct {
	symbols map[string]bool
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{symbols: make(map[string]bool)}
}

func (s *SymbolTable) DeclareVariable(name string) {
	s.symbols[name] = true
}

func (s *SymbolTable) IsDeclared(name string) bool {
	_, exists := s.symbols[name]
	return exists
}

func (node *ASTNode) AnalyzeSemantics(symbolTable *SymbolTable) {
	switch node.Type {
	case "assignment":
		varName := node.Value
		symbolTable.DeclareVariable(varName)
	case "variable":
		varName := node.Value
		if !symbolTable.IsDeclared(varName) {
			fmt.Printf("Semantic Error: Variable '%s' used before declaration\n", varName)
			fmt.Println("Exited before compilation.")
			os.Exit(1)
		}
	}

	for _, child := range node.Children {
		child.AnalyzeSemantics(symbolTable)
	}
}
