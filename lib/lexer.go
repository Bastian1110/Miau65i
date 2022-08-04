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
	fmt.Println("Actual Position :", pos)
	fmt.Println("Actual Buffer :", Buffer[pos:])
	if pos >= len(Buffer) {
		return "1"
	} else {
		fmt.Println("Buffer Inside : ", Buffer[pos:])
		space := regexp.MustCompile(`\S`).FindStringIndex(Buffer[pos:])

		if len(space) != 0 {
			fmt.Println("Space Match : ", space)
			pos = space[0]
		} else {
			return "1"
		}

		for token, exp := range regex {
			fmt.Println("Pos", pos, "buffer : ", Buffer[pos:])
			match := exp.FindStringIndex(Buffer[pos:])
			if len(match) != 0 {
				fmt.Println(Buffer[match[0]:match[1]], token)
				pos = match[1] + 1
				return "0"
			}
		}
		return "1"
	}
}

func Test() {
	Buffer = "11 a 4 5"
	x := 0
	for {
		state := getToken()
		if state == "1" {
			break
		}
		x++
		if x == 3 {
			break
		}
	}
}
