package lib

import (
	"fmt"
)

func decimalToHexStr(number int) string {
	return fmt.Sprintf("%04x", number)
}

type MemoryTable struct {
	cursor               int
	variables            map[string]string
	conditionalIdCounter int
}

func NewMeMemoryTable() *MemoryTable {
	return &MemoryTable{variables: make(map[string]string), cursor: 516}
}

func (mt *MemoryTable) Malloc(varName string) string {
	hexMem := decimalToHexStr(mt.cursor)
	mt.cursor += 2
	mt.variables[varName] = hexMem
	return hexMem
}

func (mt *MemoryTable) RetriveMem(varName string) (string, bool) {
	if value, ok := mt.variables[varName]; ok {
		return value, ok
	}
	return "", false
}

func (mt *MemoryTable) RegisterConditional() string {
	mt.conditionalIdCounter++
	return fmt.Sprintf("%02x", mt.conditionalIdCounter)
}

func (node *ASTNode) GenerateAssembly(mt *MemoryTable) string {
	switch node.Type {
	case "program":
		program := "\nPROGRAM:\n"
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
		return fmt.Sprintf("  lda %s\n  sta $%s\n\n", rhs, memAddress)

	case "binary_expr":
		ass := ""

		lhs := node.Children[0].GenerateAssembly(mt)
		ass += lhs

		rhs := node.Children[1].GenerateAssembly(mt)
		if node.Value == "add" {
			ass += fmt.Sprintf("\n  clc\n  adc %s", rhs)
		} else {
			ass += fmt.Sprintf("\n  clc\n  sbc %s", rhs)
		}
		return ass
	case "block":
		program := fmt.Sprintf("%s:\n", node.Value)
		for _, child := range node.Children {
			program += child.GenerateAssembly(mt)
		}
		return program
	case "if_statement":
		ass := ""
		ass += node.Children[0].GenerateAssembly(mt)
		conditionalId := mt.RegisterConditional()
		blockLabel := fmt.Sprintf("IF_%s", conditionalId)

		ass += fmt.Sprintf("%s\n\n", blockLabel)

		for _, child := range node.Children[1].Children {
			ass += child.GenerateAssembly(mt)
		}
		ass += fmt.Sprintf("%s:\n", blockLabel)
		return ass
	case "boolean_expr":
		ass := ""
		lhs := node.Children[0].GenerateAssembly(mt)
		rhs := node.Children[1].GenerateAssembly(mt)
		ass += fmt.Sprintf("  lda %s\n", lhs)
		ass += fmt.Sprintf("  cmp %s\n", rhs)
		switch node.Value {
		case "lrt":
			ass += "  bmi "
		case "grt":
			ass += "  bpl "
		case "eql":
			ass += "  bmi " // TODO : Implement the ==
		}
		return ass
	case "print":
		rhs := node.Children[0].GenerateAssembly(mt)
		ass := fmt.Sprintf("  lda %s\n  sta value\n  lda %s + 1\n  sta value + 1\n  jsr printNumber\n  jsr delay\n\n", rhs, rhs)
		return ass

	case "goto":
		return fmt.Sprintf("  jmp %s\n\n", node.Children[0].Value)
	case "variable":
		memAddress, _ := mt.RetriveMem(node.Value)
		return fmt.Sprintf("$%s", memAddress)
	case "number":
		return fmt.Sprintf("#%s", node.Value)
	}
	return ""
}
