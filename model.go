package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

type model struct {
	Hostname      string
	Username      string
	Os            string
	Kernel        string
	Uptime        time.Time
	Packages      int
	Shell         string
	Theme         string
	Icons         string
	Terminal      string
	CPU           string
	GPU           string
	MemoryCurrent int
	MemoryTotal   int
	Follow        bool
	FetchSpeed    time.Duration
}

func initialModel(follow bool) model {
	log.Debug("Initializing")

	return model{Follow: follow}
}

func (m model) Init() tea.Cmd {
	log.Debug("Init")
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Debug("Update")

	m.fetchData()

	if !m.Follow {
		return m, tea.Quit
	}

	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {
		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	log.Debug("View")
	return fmt.Sprintf(`
    %s@%s
    %s
    OS: %s
    Kernel: %s
    Uptime: %s
    Packages: %d
    Shell: %s
    Theme: %s
    Icons: %s
    Terminal: %s
    CPU: %s
    GPU: %s
    Memory: %d / %d
    Fetch speed: %s
`,
		m.Username,
		m.Hostname,
		strings.Repeat("-", len(m.Hostname)+len(m.Username)+1),
		m.Os,
		m.Kernel,
		time.Since(m.Uptime),
		m.Packages,
		m.Shell,
		m.Theme,
		m.Icons,
		m.Terminal,
		m.CPU,
		m.GPU,
		m.MemoryCurrent,
		m.MemoryTotal,
		m.FetchSpeed,
	)
}

func (m *model) fetchData() {
	start := time.Now()
	m.Hostname = getHostname()
	m.Username = getUsername()
	m.Os = getOS()
	m.Kernel = getKernel()
	m.Packages = getPackages()
	m.Shell = getShell()
	m.Theme = getTheme()
	m.Icons = getIcons()
	m.Terminal = getTerminal()
	m.CPU = getCPU()
	m.GPU = getGPU()
	m.MemoryCurrent = getCurrentMemory()
	m.MemoryTotal = getTotalMemory()
	m.FetchSpeed = time.Since(start)
}

func getTotalMemory() int {
	return 0
}

func getCurrentMemory() int {
	return 0
}

func getGPU() string {
	return "Not implemented"
}

func getCPU() string {
	cpuinfo, err := os.Open("/proc/cpuinfo")
	if err != nil {
		log.Warn("error when reading cpuinfo", "err", err)
	}
	defer cpuinfo.Close()

	key := "model name"
	var modelName string

	s := bufio.NewScanner(cpuinfo)
	for s.Scan() {
		line := s.Text()
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 && strings.TrimSpace(parts[0]) == key {
			modelName = strings.TrimSpace(parts[1])
		}
	}

	// TODO: Add CPU Hz

	return modelName
}

func getTerminal() string {
	return "Not implemented"
}

func getIcons() string {
	return "Not implemented"
}

func getTheme() string {
	return "Not implemented"
}

func getShell() string {
	return "Not implemented"
}

func getPackages() int {
	return 0
}

func getKernel() string {
	return "Not implemented"
}

func getOS() string {
	return "Not implemented"
}

func getHostname() string {
	hostname, _ := os.Hostname()
	return hostname
}

func getUsername() string {
	username := os.Getenv("USER")
	return username
}
