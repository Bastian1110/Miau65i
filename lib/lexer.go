package lib

import (
	"fmt"
	"regexp"
)

var Buffer string

var regex = make(map[string]*regexp.Regexp)

//var space = regexp.Compile(`\s`)
var pos = 0

func Lexer() {
	var rules = make(map[string]string)
	rules["float"] = "([0-9]*[.])+[0-9]+"
	rules["int"] = "[0-9]+"
	rules["var"] = "[a-zA-Z0-9_]+"
	for token, rule := range rules {
		regex[token] = regexp.MustCompile(rule)
	}
}

func getToken() string {
	fmt.Println("---------")
	fmt.Println("Actual Buffer :", Buffer)
	if len(Buffer) < 0 {
		return "1"
	} else {
		space := regexp.MustCompile(`\S`).FindStringIndex(Buffer[pos:])

		if len(space) == 0 {
			return "1"
		} else if space[0] != 0 {
			Buffer = Buffer[space[0]:]
		} else {
			Buffer = Buffer[0:]
		}

		for token, exp := range regex {
			fmt.Println("Atemmpt : ", token)
			match := exp.FindStringIndex(Buffer)
			if len(match) != 0 {
				fmt.Println("Token Faund !")
				fmt.Println(Buffer[match[0]:match[1]], token)
				Buffer = Buffer[match[1]:]
				return "0"
			}
			fmt.Println("Token Not Faund")
		}
		return "1"
	}
}

func Test() {
	re := regexp.MustCompile("(?P<string>[a-zA-Z]+)|(?P<int>[0-9]+)")
	input := "aa 1"
	groupNames := re.SubexpNames()
	fmt.Println(groupNames)
	firstMatch := re.FindStringIndex(input)
	fmt.Println(input[firstMatch[0]:firstMatch[1]])
	input = input[firstMatch[1]:]
	secondMatch := re.FindStringIndex(input)
	fmt.Println(input[secondMatch[0]:secondMatch[1]])
}
