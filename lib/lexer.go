package lib

import (
	"fmt"
	"regexp"
	"strings"
)

var buffer string
var regex *regexp.Regexp

type Token struct {
	Type  string
	Value string
}

func Lexer() {
	var rules = make(map[string]string)
	var regexParts []string

	// Lexer Rules
	rules["int"] = `\b[0-9]+\b`              // Integers
	rules["var"] = `\b[a-z_][a-zA-Z_0-9]*\b` // Variables (identifiers, lowercase to differentiate from labels)
	rules["label"] = `\b[A-Z][A-Z_0-9]*\b`   // Labels (uppercase identifiers)
	rules["ass"] = `=`                       // Assignment
	rules["add"] = `\+`                      // Addition
	rules["sub"] = `\-`                      // Subtraction (if needed)
	rules["grt"] = `>`                       // Greater than
	rules["lrt"] = `<`                       // Less than
	rules["opar"] = `\(`                     // Open parenthesis
	rules["cpar"] = `\)`                     // Close parenthesis
	rules["okey"] = `\{`                     // Open curly brace
	rules["ckey"] = `\}`                     // Close curly brace
	rules["if"] = `\bif\b`                   // if condition
	rules["goto"] = `\bgoto\b`               // goto keyword
	rules["print"] = `\bprint\b`             // print function

	for token, rule := range rules {
		if token == "if" || token == "goto" || token == "print" || token == "label" {
			regexParts = append([]string{"(?P<" + token + ">" + rule + ")"}, regexParts...)
		} else {
			regexParts = append(regexParts, "(?P<"+token+">"+rule+")")
		}
	}

	regex = regexp.MustCompile(strings.Join(regexParts, "|"))
}

func getToken() Token {
	buffer = strings.TrimSpace(buffer)
	if buffer == "" {
		return Token{Type: "eof", Value: ""}
	}
	match := regex.FindStringSubmatch(buffer)
	if len(match) == 0 {
		fmt.Println("No match found") // Debug print
		return Token{Type: "err", Value: ""}
	}
	index := regex.FindStringIndex(buffer)
	if len(index) < 2 {
		return Token{Type: "err", Value: ""}
	}
	buffer = buffer[index[1]:]
	tokenNames := regex.SubexpNames()
	for i, name := range tokenNames {
		if i != 0 && name != "" && match[i] != "" {
			return Token{Type: name, Value: match[i]}
		}
	}
	return Token{Type: "err", Value: ""}
}

func Tokenize(input string) []Token {
	var tokens []Token
	buffer = input
	for len(buffer) > 0 {
		token := getToken()
		if token.Type == "err" {
			fmt.Println("Failed tokenize 😐")
			break // Exit if there's an error to avoid infinite loop
		} else if token.Type != "eof" {
			tokens = append(tokens, token)
		} else {
			break // Exit if we reach the end of file
		}
	}
	return tokens
}
