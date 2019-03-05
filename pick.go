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

func readStdin() []string {
	scanner := bufio.NewScanner(os.Stdin)
	var list []string

	for scanner.Scan() {
		list = append(list, scanner.Text())
	}

	return list
}

// let's own 20 rows of the terminal
// 20th row contains: "80/360 > [current search]"
// remaining 11 rows have top matches
// one match is selected as current
func printInterface(list []string, w *bufio.Writer) {
	for i, e := range list {
		if i == 20 {
			break
		}
		w.WriteString(e + "\n")
	}
	w.WriteString("> ")
	w.Flush()
}

func main() {
	listIn := readStdin()

	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)

	if err != nil {
		fmt.Printf("excuse me: %v", err)
	}

	ttyw := bufio.NewWriter(tty)
	printInterface(listIn, ttyw)

	// outScanner := bufio.NewScanner(os.Stdout)
	outScanner := bufio.NewScanner(tty)

	var result string
	for outScanner.Scan() {
		filterResults := Filter(listIn, outScanner.Text())
		ttyw.WriteString(strings.Join(filterResults, " ") + "\n")

		result = filterResults[0]
		break
	}

	// fmt.Println(Filter(listIn, text))
	// os.Stdout = originalStdOut

	// fmt.Println(result)
	fmt.Printf(result)
}
