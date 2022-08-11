package lib

import (
	"bufio"
	"fmt"
	"os"
)

func ReadFile(path string) []string {
	var res []string
	if path[len(path)-5:] != ".miau" {
		fmt.Println("Expecting a miau file ... 🥺")
	}
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Cant open the file 😬")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		res = append(res, scanner.Text())
	}
	return res
}
