package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/mikepjb/pick/src/ui"
)

const (
	keyReturn = 0x0A
	backspace = 0x7F
	keyCtrlW  = 0x17
	keyCtrlN  = 0x0E
	keyCtrlP  = 0x10
)

func fuzzy(filter string) (*regexp.Regexp, error) {
	newFilter := "(?i)" // ?i makes is case insensitive

	regexpString := newFilter + strings.Join(strings.Split(filter, ""), ".*?")
	return regexp.Compile(regexpString)
}

type match struct {
	filePath string
	rank     int
}

type matches []match

// TODO: correct make errors
func (m matches) Len() int           { return len(m) }
func (m matches) Less(i, j int) bool { return m[i].rank < m[j].rank }
func (m matches) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }

func filter(col []string, filter string) []string {
	matches := matches{}

	freg, err := fuzzy(strings.Replace(filter, ".", "\\.", -1))

	if err != nil {
		fmt.Errorf("bad regexp compilation %v", err)
	}

	for _, e := range col {
		fuzzyMatches := freg.FindAllString(e, -1)

		if len(fuzzyMatches) != 0 {
			fmatch := fuzzyMatches[len(fuzzyMatches)-1]

			if len(fmatch) >= len(filter) {
				matches = append(matches, match{e, len(fmatch)})
			}
		}
	}

	sort.Sort(matches)

	filePaths := []string{}

	for _, match := range matches {
		filePaths = append(filePaths, match.filePath)
	}

	return filePaths
}

func readStdin() []string {
	scanner := bufio.NewScanner(os.Stdin)
	var list []string

	for scanner.Scan() {
		list = append(list, scanner.Text())
	}

	return list
}

func printInterface(cpos int, list []string, w *bufio.Writer) {
	w.WriteString("\033[19A") // move up 19 lines
	w.WriteString("\033[1F")  // move to beginning of previous line

	for i, e := range list {
		if i == 20 {
			break
		}
		if cpos == i {
			w.WriteString("\033[37m")
			w.WriteString("\033[0m")
		}
		w.WriteString(e)
		if cpos == i {
			w.WriteString("\033[37m")
		}
		w.WriteString("\033[K" + "\n")
	}

	if len(list) < 20 {
		fmt.Fprintf(w, strings.Repeat("\033[K"+"\n", (20-len(list))))
	}

	fmt.Fprintf(w, "%03d > ", len(list))
	w.Flush()
}

func contains(filePath string, list []string) bool {
	for _, ignoreMatch := range list {
		if strings.HasPrefix(filePath, strings.TrimPrefix(ignoreMatch, "/")) {
			return true
		}
	}
	return false
}

func filterByList(source []string, list []string) []string {
	filteredList := []string{}

	for _, filePath := range source {
		if !contains(filePath, list) {
			filteredList = append(filteredList, filePath)
		}
	}

	return filteredList
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
	cpos := 0 // choice position
	var b []byte = make([]byte, 1)
	listIn := readStdin()

	ignoreFile, err := ioutil.ReadFile(".gitignore")

	if err == nil {
		ignoreList := strings.Split(strings.TrimSuffix(string(ignoreFile), "\n"), "\n")
		listIn = filterByList(listIn, ignoreList)
	}

	printInterface(cpos, listIn, ttyw)
	ttyw.WriteString("\033[s") // save cursor position

	for {
		ttyr.Read(b)

		filterResults := filter(listIn, userSearch)

		switch b[0] {
		case keyReturn:
			fmt.Printf(filterResults[cpos])
			return
		case backspace:
			userSearch = userSearch[:len(userSearch)-1]
		case keyCtrlW:
			userSearch = ""
		case keyCtrlN:
			cpos++
		case keyCtrlP:
			cpos--
		default:
			userSearch = userSearch + string(b)
		}

		printInterface(cpos, filter(listIn, userSearch), ttyw)
		ttyw.WriteString("\033[u")   // return cursor position
		ttyw.WriteString("\033[1D")  // left one
		ttyw.WriteString("\033[0K")  // return cursor position
		ttyw.WriteString(userSearch) // return cursor position
		ttyw.Flush()
	}
}
