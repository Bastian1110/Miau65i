package lib

import (
	"fmt"
	"regexp"
	"strings"
)

var buffer string
var regex *regexp.Regexp

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

func getToken() []string {
	var token [2]string
	buffer = strings.TrimSpace(buffer)
	if buffer == "" {
		return []string{"eof", ""}
	}
	match := regex.FindStringSubmatch(buffer)
	if len(match) == 0 {
		fmt.Println("No match found") // Debug print
		token[0] = "err"
		return token[:]
	}
	index := regex.FindStringIndex(buffer)
	if len(index) < 2 {
		token[0] = "err"
		return token[:]
	}
	buffer = buffer[index[1]:]
	result := make(map[string]string)
	tokenNames := regex.SubexpNames()
	for i, name := range tokenNames {
		if i != 0 && name != "" && match[i] != "" {
			result[name] = match[i]
			token[0] = name
			token[1] = result[name]
			return token[:]
		}
	}
	token[0] = "err"
	return token[:]
}

func Tokenize(input string) [][]string {
	var result [][]string
	buffer = input
	max := 0
	for len(buffer) > 0 && max < 10 {
		token := getToken()
		if token[0] == "err" {
			fmt.Println("Failed tokenize 😐")
		}
		result = append(result, token)
		max++
	}
	return result
}
