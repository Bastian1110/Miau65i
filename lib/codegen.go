package lib

import (
	"fmt"
)

type MemoryTable struct {
	symbols map[string]bool
}

func (node *ASTNode) GenerateAssembly() string {
	switch node.Type {
	case "assignment":
		// Assuming the right-hand side is a simple expression or a number for simplicity
		rhs := node.Children[0].GenerateAssembly()          // Recursively generate code for the right-hand side
		return fmt.Sprintf("MOV %s, %s\n", node.Value, rhs) // 'MOV' is just an example
	case "binary_expr":
		if node.Value == "add" {
			left := node.Children[0].GenerateAssembly()
			right := node.Children[1].GenerateAssembly()
			// Simplified: assuming the left part is in a register and we add the right part to it
			return fmt.Sprintf("%sADD %s, %s\n", left, left, right)
		}
	case "number":
		return node.Value // Directly return the number, assuming it will be used in a context that makes sense
	}
	return ""
}
