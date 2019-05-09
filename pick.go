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
	KeyReturn = 0x0A
	Backspace = 0x7F
	KeyCtrlW  = 0x17
	KeyCtrlN  = 0x0E
	KeyCtrlP  = 0x10
)

// Inserts '+.*' between all chars to make the regex fuzzy
func Fuzzy(filter string) (*regexp.Regexp, error) {
	newFilter := "(?i)" // ?i makes is case insensitive

	return regexp.Compile(
		newFilter + strings.Join(strings.Split(filter, ""), ".*?"))
}

type Match struct {
	filePath string
	rank     int
}

type Matches []Match

func (m Matches) Len() int           { return len(m) }
func (m Matches) Less(i, j int) bool { return m[i].rank < m[j].rank }
func (m Matches) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }

func Filter(col []string, filter string) []string {
	matches := Matches{}

	fuzzy, err := Fuzzy(strings.Replace(filter, ".", "\\.", -1))

	if err != nil {
		fmt.Errorf("bad regexp compilation %v\n", err)
	}

	for _, e := range col {
		fuzzyMatches := fuzzy.FindAllString(e, -1)

		if len(fuzzyMatches) != 0 {
			match := fuzzyMatches[len(fuzzyMatches)-1]

			if len(match) >= len(filter) {
				matches = append(matches, Match{e, len(match)})
			}
		}

		// match, _ := regexp.MatchString(regexpFilter, e)
		// if strings.Contains(e, filter) {
		// 	exactMatches = append(exactMatches, e)
		// } else if match {
		// 	fuzzyMatches = append(fuzzyMatches, e)
		// }
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
			w.WriteString("\033[40m")
		}
		w.WriteString(e)
		if cpos == i {
			w.WriteString("\033[0m")
		}
		w.WriteString("\033[K" + "\n")
	}

	if len(list) < 20 {
		fmt.Fprintf(w, strings.Repeat("\033[K"+"\n", (20-len(list))))
	}

	fmt.Fprintf(w, "%03d > ", len(list))
	w.Flush()
}

// for each filePath in a list
// match against ALL gitignores
// if any match -> return false
// otherwise -> return true
// based on this, include in new list
func match(filePath string, list []string) bool {
	for _, ignoreMatch := range list {
		// ignoring / prefix allows entries in .gitignore to work but is a poor fix.
		// regex matching should be employed (or whatever .gitignore uses to match)
		if strings.HasPrefix(filePath, strings.TrimPrefix(ignoreMatch, "/")) {
			return true
		}
	}
	return false
}

func filterByList(source []string, list []string) []string {
	filteredList := []string{}

	for _, filePath := range source {
		if !match(filePath, list) {
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

		filterResults := Filter(listIn, userSearch)

		switch b[0] {
		case KeyReturn:
			fmt.Printf(filterResults[cpos])
			return
		case Backspace:
			userSearch = userSearch[:len(userSearch)-1]
		case KeyCtrlW:
			userSearch = ""
		case KeyCtrlN:
			cpos++
		case KeyCtrlP:
			cpos--
		default:
			userSearch = userSearch + string(b)
		}

		printInterface(cpos, Filter(listIn, userSearch), ttyw)
		ttyw.WriteString("\033[u")   // return cursor position
		ttyw.WriteString("\033[1D")  // left one
		ttyw.WriteString("\033[0K")  // return cursor position
		ttyw.WriteString(userSearch) // return cursor position
		ttyw.Flush()
	}
}
