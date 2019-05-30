package main

import (
	"bufio"
	"fmt"
	"github.com/mikepjb/pick/src/filter"
	"io"
	"io/ioutil"
	"os"
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

func readStdin() []string {
	scanner := bufio.NewScanner(os.Stdin)
	var list []string

	for scanner.Scan() {
		list = append(list, scanner.Text())
	}

	return list
}

type errWriter struct {
	w   *bufio.Writer
	err error
}

func (e *errWriter) WriteString(s string) {
	if e.err != nil {
		return
	}
	_, e.err = e.w.WriteString(s)
}

func (e *errWriter) Err() error {
	return e.err
}

func printInterface(cpos int, list []string, w *bufio.Writer) {
	ew := &errWriter{w: w}
	ew.WriteString("\033[19A") // move up 19 lines
	ew.WriteString("\033[1F")  // move to beginning of previous line

	for i, e := range list {
		if i == 20 {
			break
		}
		if cpos == i {
			ew.WriteString("\033[37m")
		}
		ew.WriteString(e)
		if cpos == i {
			ew.WriteString("\033[0m")
		}
		ew.WriteString("\033[K" + "\n")
	}

	if len(list) < 20 {
		fmt.Fprintf(w, strings.Repeat("\033[K"+"\n", (20-len(list))))
	}

	if ew.Err() != nil {
		fmt.Printf("could not write interface: %v\n", ew.Err())
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
	cpos := 0 // choice position
	var b []byte = make([]byte, 1)
	listIn := readStdin()

	ignoreFile, err := ioutil.ReadFile(".gitignore")

	if err == nil {
		ignoreList := strings.Split(strings.TrimSuffix(string(ignoreFile), "\n"), "\n")
		listIn = filter.FilterByList(listIn, ignoreList)
	}

	printInterface(cpos, listIn, ttyw)
	ttyw.WriteString("\033[s") // save cursor position

	for {
		ttyr.Read(b)

		filterResults := filter.Filter(listIn, userSearch)

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

		printInterface(cpos, filter.Filter(listIn, userSearch), ttyw)
		ttyw.WriteString("\033[u")   // return cursor position
		ttyw.WriteString("\033[1D")  // left one
		ttyw.WriteString("\033[0K")  // return cursor position
		ttyw.WriteString(userSearch) // return cursor position
		ttyw.Flush()
	}
}
