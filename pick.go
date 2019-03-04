package main

import (
	"fmt"
	"strings"
)

func Filter(col []string, filter string) []string {
	result := []string{}

	for _, e := range col {
		if strings.Contains(e, filter) {
			result = append(result, e)
		}
	}
	return result
}

func main() {
	fmt.Println("hello")
}
