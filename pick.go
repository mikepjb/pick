package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func Filter(col []string, filter string) []string {
	result := []string{}

	regexpFilter := strings.Join(strings.Split(filter, ""), "+.*") + "+"

	for _, e := range col {
		match, _ := regexp.MatchString(regexpFilter, e)
		if match {
			result = append(result, e)
		}
	}
	return result
}

// let's own 12 rows of the terminal
// 12th row contains: "80/360 > [current search]"
// remaining 11 rows have top matches
// one match is selected as current
// func printInterface() {
// }

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
