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
	rules["float"] = "([0-9]*[.])+[0-9]+"
	rules["int"] = "[0-9]+"
	rules["var"] = "[a-zA-Z_]+"
	for token, rule := range rules {
		regexParts = append(regexParts, "(?P<"+token+">"+rule+")")
	}
	regex = regexp.MustCompile(strings.Join(regexParts, "|"))
}

func getToken() []string {
	var token [2]string
	match := regex.FindStringSubmatch(buffer)
	if len(match[0]) == 0 {
		token[0] = "err"
		return token[:]
	}
	index := regex.FindStringIndex(buffer)
	buffer = buffer[index[1]:]
	result := make(map[string]string)
	tokenNames := regex.SubexpNames()
	for i, name := range tokenNames {
		if i != 0 && name != "" {
			result[name] = match[i]
			if result[name] != "" {
				token[0] = name
				token[1] = result[name]
				return token[:]
			}
		}
	}
	token[0] = "err"
	return token[:]
}

func Tokenize(input string) [][]string {
	var result [][]string
	buffer = input
	for len(buffer) > 0 {
		token := getToken()
		if token[0] == "err" {
			fmt.Println("Failed tokenize 😐")
		}
		result = append(result, token)
	}
	return result
}
