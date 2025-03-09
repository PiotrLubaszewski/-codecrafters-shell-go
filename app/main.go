package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
    "path/filepath"
)

// CommandRegistry zarządza rejestracją i wykonywaniem komend
type CommandRegistry struct {
	commands map[string]func([]string)
}

// NewCommandRegistry tworzy nowy rejestr komend
func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{commands: make(map[string]func([]string))}
}

// Register dodaje komendę do rejestru
func (cr *CommandRegistry) Register(name string, handler func([]string)) {
	cr.commands[name] = handler
}

// Execute wykonuje komendę z rejestru lub jako polecenie zewnętrzne
func (cr *CommandRegistry) Execute(command string, args []string) {
	if action, exists := cr.commands[command]; exists {
		action(args)
	} else {
		cr.executeExternal(command, args)
	}
}

// executeExternal wykonuje komendę zewnętrzną, jeśli nie jest w rejestrze
func (cr *CommandRegistry) executeExternal(command string, args []string) {
	out, err := exec.Command(command, args...).Output()
	if err != nil {
		fmt.Printf("%s: command not found\n", command)
		return
	}
	fmt.Print(string(out))
}

// Shell reprezentuje powłokę z rejestrem komend
type Shell struct {
	registry *CommandRegistry
}

// NewShell tworzy nową instancję powłoki
func NewShell(registry *CommandRegistry) *Shell {
	return &Shell{registry: registry}
}

// Run uruchamia pętlę shella
func (s *Shell) Run() {
	for {
		command, args, err := s.readCommandAndArgs()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}
		if command == "" {
			continue
		}
		s.registry.Execute(command, args)
	}
}

// readCommandAndArgs wczytuje i zwraca komendę użytkownika
func (s *Shell) readCommandAndArgs() (string, []string, error) {
	fmt.Print("$ ")
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", nil, err
		}
		fmt.Println("\nExit")
		os.Exit(0)
	}

	trimmed := strings.TrimSpace(scanner.Text())
	splitted := strings.Fields(trimmed)
	if len(splitted) == 0 {
		return "", nil, nil
	}

	return splitted[0], splitted[1:], nil
}

// handleExit obsługuje wyjście z programu
func handleExit(args []string) {
	exitCode := 0
	if len(args) > 0 {
		var err error
		exitCode, err = strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Invalid exit code")
			return
		}
	}
	os.Exit(exitCode)
}

// handleEcho obsługuje polecenie echo
func handleEcho(args []string) {
	fmt.Println(strings.Join(args, " "))
}

// handleType obsługuje polecenie type
func handleType(args []string) {
	if len(args) == 0 {
		fmt.Println("type: missing argument")
		return
	}

	command := args[0]
	if path, err := exec.LookPath(command); err == nil {
		fmt.Printf("%s is %s\n", command, path)
	} else {
		fmt.Printf("%s: not found\n", command)
	}
}

func handlePwd(args []string) {
	ex, err := os.Executable()
    if err != nil {
        panic(err)
    }
    exPath := filepath.Dir(ex)
    fmt.Println(exPath)
}

// InitializeRegistry tworzy i inicjalizuje rejestr komend
func InitializeRegistry() *CommandRegistry {
	registry := NewCommandRegistry()
	registry.Register("exit", handleExit)
	registry.Register("echo", handleEcho)
	registry.Register("type", handleType)
	registry.Register("pwd", handlePwd)
	return registry
}

func main() {
	registry := InitializeRegistry()
	shell := NewShell(registry)
	shell.Run()
}
