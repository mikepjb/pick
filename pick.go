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

// Display and Output handling

// take stdin - read the list input to pick
// write info the /dev/tty.. try this now.

func main() {
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		fmt.Printf("excuse me: %v", err)
	}

	ttyw := bufio.NewWriter(tty)
	// ttyr := bufio.NewReader(tty)

	scanner := bufio.NewScanner(os.Stdin)

	var haystack []string

	for scanner.Scan() {
		haystack = append(haystack, scanner.Text())
	}

	for i, e := range haystack {
		if i == 15 {
			break
		}
		ttyw.WriteString(e + "\n")
		// fmt.Println(e) // this goes to stdout
	}

	ttyw.Flush()

	fmt.Print("> ")

	// outScanner := bufio.NewScanner(os.Stdout)
	outScanner := bufio.NewScanner(tty)

	var result string
	for outScanner.Scan() {
		// fmt.Print(outScanner.Text())
		filterResults := Filter(haystack, outScanner.Text())
		fmt.Println(filterResults)

		result = filterResults[0]
		break
	}

	// fmt.Println(Filter(haystack, text))
	// os.Stdout = originalStdOut

	fmt.Println(result)
}
