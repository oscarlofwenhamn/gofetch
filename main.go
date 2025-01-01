package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

func main() {
	log.SetLevel(log.DebugLevel)
	logfile, err := os.OpenFile("./debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("error when opening logfile: ", err)
		os.Exit(1)
	}
	defer logfile.Close()
	log.SetOutput(logfile)
	var follow bool

	flag.BoolVar(&follow, "follow", false, "keep polling the system for information updates")
	flag.Parse()

	p := tea.NewProgram(initialModel(follow))
	log.Info("Starting program")
	if _, err := p.Run(); err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}
}
