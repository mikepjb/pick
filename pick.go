package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
)

const (
	KeyReturn = 0xA
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

	fmt.Fprintf(w, "%d > ", len(list))

	w.Flush()
}

func main() {
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)

	if err != nil {
		fmt.Printf("excuse me: %v", err)
	}

	ttyw := bufio.NewWriter(tty)

	listIn := readStdin()
	printInterface(listIn, ttyw)

	ttyr := io.Reader(tty)

	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	// restore the echoing state when exiting
	defer exec.Command("stty", "-F", "/dev/tty", "echo").Run()

	// ensure we use echo even on sigterm/ctrl+c
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		exec.Command("stty", "-F", "/dev/tty", "echo").Run()
		os.Exit(1)
	}()

	var userSearch string

	var b []byte = make([]byte, 1)

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
