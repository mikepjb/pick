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
func printInterface() {
}

func main() {
	originalStdOut := os.Stdout

	scanner := bufio.NewScanner(os.Stdin)

	var haystack []string

	for scanner.Scan() {
		haystack = append(haystack, scanner.Text())
	}

	for _, e := range haystack {
		fmt.Println(e)
	}

	fmt.Print("> ")

	outScanner := bufio.NewScanner(os.Stdout)

	var result string
	for outScanner.Scan() {
		// fmt.Print(outScanner.Text())
		filterResults := Filter(haystack, outScanner.Text())
		fmt.Println(filterResults)

		result = filterResults[0]
		break
	}

	// fmt.Println(Filter(haystack, text))
	os.Stdout = originalStdOut
	fmt.Println(result)
}
