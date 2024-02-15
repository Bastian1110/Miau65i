package lib

import (
	"fmt"
	"strings"
)

const (
	TOKEN_INT   = "int"
	TOKEN_VAR   = "var"
	TOKEN_ASS   = "ass"
	TOKEN_ADD   = "add"
	TOKEN_SUB   = "sub"
	TOKEN_GRT   = "grt"
	TOKEN_LRT   = "lrt"
	TOKEN_OPAR  = "opar"
	TOKEN_CPAR  = "cpar"
	TOKEN_OKEY  = "okey"
	TOKEN_CKEY  = "ckey"
	TOKEN_IF    = "if"
	TOKEN_GOTO  = "goto"
	TOKEN_PRINT = "print"
	TOKEN_LABEL = "label"
)

type Parser struct {
	tokens  []Token
	current int
}

type ASTNode struct {
	Type     string
	Value    string
	Children []*ASTNode
}

func (node *ASTNode) Show(level int) {
	indent := strings.Repeat("  ", level)
	fmt.Printf("%s%s: %s\n", indent, node.Type, node.Value)
	for _, child := range node.Children {
		child.Show(level + 1)
	}
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) ParseProgram() *ASTNode {
	node := &ASTNode{Type: "program", Children: []*ASTNode{}}
	max := 0
	for !p.atEnd() && max < 5 {
		node.Children = append(node.Children, p.parseStatement())
		max++
	}
	return node
}

func (p *Parser) parseStatement() *ASTNode {
	if p.match(TOKEN_VAR) {
		p.consume("Expecting variable name.", TOKEN_VAR)
		return p.parseAssignmentStatement()
	}
	return nil
}

func (p *Parser) parseAssignmentStatement() *ASTNode {
	varName := p.previous().Value
	p.consume("Expect '=' after variable name.", TOKEN_ASS)
	expr := p.parseExpression()
	return &ASTNode{Type: "assignment", Value: varName, Children: []*ASTNode{expr}}
}

func (p *Parser) parseExpression() *ASTNode {
	node := p.parseTerm()

	for !p.atEnd() && (p.match(TOKEN_ADD, TOKEN_SUB)) {
		p.consume("Expecting arithmetic name.", TOKEN_ADD, TOKEN_SUB)
		operator := p.previous().Type
		right := p.parseTerm()
		node = &ASTNode{
			Type:     "binary_expr",
			Value:    operator,
			Children: []*ASTNode{node, right},
		}
	}

	return node
}

func (p *Parser) parseTerm() *ASTNode {
	if p.match(TOKEN_INT) {
		p.consume("Expecting int", TOKEN_INT)
		return &ASTNode{Type: "number", Value: p.previous().Value}
	} else if p.match(TOKEN_VAR) {
		p.consume("Expecting var", TOKEN_VAR)
		return &ASTNode{Type: "variable", Value: p.previous().Value}
	}

	return nil
}

func (p *Parser) match(types ...string) bool {
	fmt.Println("Current : ", p.current, p.tokens[p.current], " Matching : ", types)
	if p.atEnd() {
		return false
	}
	for _, typ := range types {
		if !p.atEnd() && p.tokens[p.current].Type == typ {
			return true
		}
	}
	return false
}

func (p *Parser) consume(errorMessage string, expectedTypes ...string) Token {
	fmt.Println("Parsing Token :", p.tokens[p.current], " Expected Type : ", expectedTypes)
	for _, typ := range expectedTypes {
		if !p.atEnd() && p.tokens[p.current].Type == typ {
			token := p.tokens[p.current]
			p.current++
			return token
		}
	}
	panic(fmt.Sprintf("Error: %s", errorMessage))
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) atEnd() bool {
	return p.current >= len(p.tokens)
}
