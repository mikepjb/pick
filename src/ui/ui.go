package ui

import (
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

// remove buffering & echo but recover these on exit
func SetAndProtectTerm() {
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
}
