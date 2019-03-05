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
	for i, e := range list {
		if i == 20 {
			break
		}
		w.WriteString(e + "\n")
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

	// ttyr := bufio.NewScanner(tty)
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
			fmt.Fprintf(ttyw, "\n:::%v\n", b)
			userSearch = userSearch + string(b)
		}

		printInterface(Filter(listIn, userSearch), ttyw)
		ttyw.WriteString(userSearch)
		ttyw.Flush()
	}
}
