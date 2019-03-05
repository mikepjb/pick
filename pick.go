package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/mikepjb/pick/src/ui"
)

const (
	KeyReturn = 0xA
)

// Inserts '+.*' between all chars to make the regex fuzzy
func Fuzzy(filter string) string {
	var newFilter string

	for _, r := range filter {
		newFilter += string(r)
		if string(r) != "\\" {
			newFilter += "+.*"
		}
	}

	return newFilter
}

func Filter(col []string, filter string) []string {
	result := []string{}

	regexpFilter := Fuzzy(strings.Replace(filter, ".", "\\.", -1))

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

func printInterface(list []string, w *bufio.Writer) {
	w.WriteString("\033[20A") // move up 20 lines

	for i, e := range list {
		if i == 20 {
			break
		}
		w.WriteString(e + "\033[K" + "\n")
		// ttyw.WriteString("\033[20A")
	}

	if len(list) < 20 {
		fmt.Fprintf(w, strings.Repeat("\033[K"+"\n", (20-len(list))))
	}

	fmt.Fprintf(w, "%03d > ", len(list))

	w.Flush()
}

func main() {
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)

	if err != nil {
		fmt.Printf("excuse me: %v", err)
	}

	ttyw := bufio.NewWriter(tty)
	ttyr := io.Reader(tty)

	ui.SetAndProtectTerm()

	var userSearch string
	var b []byte = make([]byte, 1)
	listIn := readStdin()
	printInterface(listIn, ttyw)

	for {
		ttyr.Read(b)

		filterResults := Filter(listIn, userSearch)

		switch b[0] {
		case KeyReturn:
			fmt.Printf(filterResults[0])
			return
		default:
			// debug statement, find the keycode
			// fmt.Fprintf(ttyw, "\n:::%v\n", b)
			userSearch = userSearch + string(b)
			// \033[2K
		}

		ttyw.WriteString("\033[s") // save cursor position
		printInterface(Filter(listIn, userSearch), ttyw)
		ttyw.WriteString("\033[u") // return cursor position
		// ttyw.WriteString(userSearch)
		ttyw.WriteString(string(b))
		ttyw.Flush()
	}
}
