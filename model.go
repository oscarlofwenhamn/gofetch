package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

type model struct {
	Hostname      string
	Username      string
	Os            string
	KernelName    string
	KernelVersion string
	KernelMachine string
	Uptime        time.Duration
	Packages      string
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
    OS: %s %s
    Kernel: %s
    Uptime: %s
    Packages: %s
    Shell: %s
    Theme: %s
    Icons: %s
    Terminal: %s
    CPU: %s
    GPU: %s
    Memory: %.0fMiB / %.0fMiB
    Fetch speed: %s
`,
		m.Username,
		m.Hostname,
		strings.Repeat("-", len(m.Hostname)+len(m.Username)+1),
		m.Os,
		m.KernelMachine,
		m.KernelVersion,
		m.Uptime,
		m.Packages,
		m.Shell,
		m.Theme,
		m.Icons,
		m.Terminal,
		m.CPU,
		m.GPU,
		float64(m.MemoryCurrent)/(1024),
		float64(m.MemoryTotal)/(1024),
		m.FetchSpeed,
	)
}

func (m *model) fetchData() {
	start := time.Now()
	m.Hostname = getHostname()
	m.Username = getUsername()
	m.KernelName, m.KernelVersion, m.KernelMachine = getKernel()
	m.Os = getOS()
	m.Uptime = getUptime()
	m.Packages = getPackages()
	m.Shell = getShell()
	m.Theme = getTheme()
	m.Icons = getIcons()
	m.Terminal = getTerminal()
	m.CPU = getCPU()
	m.GPU = getGPU()
	m.MemoryCurrent, m.MemoryTotal = getMemory()
	m.FetchSpeed = time.Since(start)
}

func getMemory() (int, int) {
	freeCmd := exec.Command("free")
	out, err := freeCmd.Output()
	if err != nil {
		log.Warn("error when running free", "err", err)
		return 0, 0
	}
	i := bytes.IndexByte(out, ':')
	if i == -1 {
		log.Warn("invalid output from free")
		return 0, 0
	}
	out = bytes.TrimLeft(out[i+1:], " ")
	total, err := readInt(out)
	if err != nil {
		log.Warn("error when reading total memory", "err", err)
	}

	i = bytes.IndexByte(out, ' ')
	out = bytes.TrimLeft(out[i:], " ")
	current, err := readInt(out)
	if err != nil {
		log.Warn("error when reading total memory", "err", err)
	}

	return current, total
}

func readInt(s []byte) (int, error) {
	i := bytes.IndexByte(s, ' ')
	b := s[:i]
	val, err := strconv.Atoi(string(b))
	return val, err
}

func getUptime() time.Duration {
	b, err := os.ReadFile("/proc/uptime")
	if err != nil {
		log.Warn("error when trying to open /proc/uptime", "err", err)
		return 0
	}

	up := string(bytes.Split(b, []byte{' '})[0])

	uptime, err := strconv.ParseFloat(up, 64)
	if err != nil {
		log.Warn("error when parsing uptime float", "err", err, "uptime", uptime)
		return 0
	}
	uptimeSeconds := int(uptime)

	return time.Duration(uptimeSeconds) * time.Second
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
	modelNameKey := "model name"
	var modelName string

	mHzKey := "cpu MHz"
	var cpuGHz float64

	vals, err := getFromKeyValueFile("/proc/cpuinfo", ":", []string{modelNameKey, mHzKey})
	if err != nil {
		log.Warn("error when reading values from cpuinfo", "err", err)
	}

	modelName, ok := vals[modelNameKey]
	if !ok {
		log.Warn("no model name found")
		modelName = "Invalid"
	}

	cpuMHz, ok := vals[mHzKey]
	if !ok {
		log.Warn("no clock speed found")
	} else {

		cpuGHz, err = strconv.ParseFloat(strings.TrimSpace(cpuMHz), 32)
		if err != nil {
			log.Warn("error when converting cpu hz", "err", err)
		}
		cpuGHz /= 1000
	}

	// TODO: Add CPU siblings(?)
	// QUESTION: Is it accurate and correct to mathematically round clock speed,
	// or should it rather be truncated?
	return fmt.Sprintf("%s @ %.3fGHz", modelName, cpuGHz)
}

func getTerminal() string {
	var term string
	if os.Getenv("WT_SESSION") != "" {
		term = "Windows Terminal"
	}
	if program := os.Getenv("TERM_PROGRAM"); program != "" {
		if term == "" {
			term = program
		} else {
			term += " (" + program + ")"
		}
	}
	return term
}

func getIcons() string {
	return "Not implemented"
}

func getTheme() string {
	return "Not implemented"
}

func getShell() string {
	shellPath := os.Getenv("SHELL")
	if shellPath == "" {
		log.Warn("no shell path found")
		return "Invalid"
	}
	shellVersionCommand := exec.Command(shellPath, "--version")
	output, err := shellVersionCommand.Output()
	if err != nil {
		log.Warn("error when fetching shell version", "err", err)
		return "Invalid"
	}
	for i, b := range output {
		if b == '\n' {
			output = output[:i]
		}
	}

	// TODO: Remove "unwanted information", e.g. within ()
	return string(output)
}

func getPackages() string {
	var packageCounters []func() (string, error)
	if _, err := exec.LookPath("dpkg"); err == nil {
		packageCounters = append(packageCounters, GetDpkgPackageCount)
	}
	if _, err := exec.LookPath("snap"); err == nil {
		packageCounters = append(packageCounters, GetSnapPackageCount)
	}

	var out string
	for _, provider := range packageCounters {
		count, err := provider()
		if err != nil {
			log.Warn("error when fetching count", "err", err)
		}
		if out != "" {
			out += ", "
		}
		out += count
	}
	return out
}

func getKernel() (string, string, string) {
	unameCmd := exec.Command("uname", "-smr")
	unameOut, err := unameCmd.Output()
	if err != nil {
		log.Warn("error when running uname command", "err", err)
		return "Invalid", "Invalid", "Invalid"
	}
	info := strings.Split(strings.TrimSpace(string(unameOut)), " ")
	log.Debug(info)
	return info[0], info[1], info[2]
}

func getOS() string {
	osNameKey := "PRETTY_NAME"
	vals, err := getFromKeyValueFile("/etc/os-release", "=", []string{osNameKey})
	if err != nil {
		log.Warn("error when reading from os-release")
		return "Invalid"
	}

	return strings.Trim(vals[osNameKey], "\"")
}

func getHostname() string {
	hostname, _ := os.Hostname()
	return hostname
}

func getUsername() string {
	username := os.Getenv("USER")
	return username
}

func getFromKeyValueFile(path, separator string, keys []string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	s := bufio.NewScanner(file)

	vals := make(map[string]string)
	for s.Scan() {
		line := s.Text()
		parts := strings.SplitN(line, separator, 2)

		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		for i, k := range keys {
			if key == k {
				keys = append(keys[:i], keys[i+1:]...)
				vals[k] = strings.TrimSpace(parts[1])
			}
		}
	}
	return vals, nil
}
