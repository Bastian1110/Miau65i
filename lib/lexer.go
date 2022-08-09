package lib

import (
	"fmt"
	"regexp"
	"strings"
)

var Buffer string

var regex *regexp.Regexp

//var space = regexp.Compile(`\s`)
var pos = 0

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

func getToken() string {
	match := regex.FindStringSubmatch(Buffer)
	index := regex.FindStringIndex(Buffer)
	Buffer = Buffer[index[1]:]
	result := make(map[string]string)
	tokenNames := regex.SubexpNames()
	for i, name := range tokenNames {
		if i != 0 && name != "" {
			result[name] = match[i]
			if result[name] != "" {
				return ("Type : " + name + " Token : " + result[name])
			}
		}
	}
	return ""
}

func Test() {
	Buffer = "abc1.2aaa "
	fmt.Println(getToken())
	fmt.Println(getToken())
	fmt.Println(getToken())
}
