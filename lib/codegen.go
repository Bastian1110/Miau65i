package lib

import (
	"fmt"
)

func decimalToHexStr(number int) string {
	return fmt.Sprintf("%04x", number)
}

type MemoryTable struct {
	cursor    int
	variables map[string]string
}

func NewMeMemoryTable() *MemoryTable {
	return &MemoryTable{variables: make(map[string]string), cursor: 516}
}

func (mt *MemoryTable) Malloc(varName string) string {
	hexMem := decimalToHexStr(mt.cursor)
	mt.cursor++
	mt.variables[varName] = hexMem
	return hexMem
}

func (mt *MemoryTable) RetriveMem(varName string) (string, bool) {
	if value, ok := mt.variables[varName]; ok {
		return value, ok
	}
	return "", false
}

func (node *ASTNode) GenerateAssembly(mt *MemoryTable) string {
	switch node.Type {
	case "program":
		program := "PROGRAM:\n"
		for _, child := range node.Children {
			program += child.GenerateAssembly(mt)
		}
		return program
	case "assignment":
		memAddress, inMalloc := mt.RetriveMem(node.Value)
		if !inMalloc {
			memAddress = mt.Malloc(node.Value)
		}
		rhs := node.Children[0].GenerateAssembly(mt)
		return fmt.Sprintf("LDA %sSTA $%s\n\n", rhs, memAddress)

	case "binary_expr":
		ass := ""

		lhs := node.Children[0].GenerateAssembly(mt)
		ass += lhs

		rhs := node.Children[1].GenerateAssembly(mt)
		if node.Value == "add" {
			ass += fmt.Sprintf("CLC\nADC %s", rhs)
		} else {
			ass += fmt.Sprintf("CLC\nSBC %s", rhs)
		}
		return ass
	case "block":
		program := fmt.Sprintf("%s:\n", node.Value)
		for _, child := range node.Children {
			program += child.GenerateAssembly(mt)
		}
		return program
	case "goto":
		return fmt.Sprintf("JMP %s\n", node.Children[0].Value)
	case "variable":
		memAddress, _ := mt.RetriveMem(node.Value)
		return fmt.Sprintf("$%s\n", memAddress)
	case "number":
		return fmt.Sprintf("#%s\n", node.Value)
	}
	return ""
}
