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

// Nowy rejestr komend
func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{commands: make(map[string]func([]string))}
}

// Inicjalizacja rejestru komend
func InitializeRegistry() *CommandRegistry {
	registry := NewCommandRegistry()
	registry.Register("exit", handleExit)
	registry.Register("echo", handleEcho)
	registry.Register("type", func(args []string) { handleType(args, registry) })
	registry.Register("pwd", handlePwd)
	return registry
}

// Rejestracja nowej komendy
func (cr *CommandRegistry) Register(name string, handler func([]string)) {
	cr.commands[name] = handler
}

// Wykonanie komendy
func (cr *CommandRegistry) Execute(command string, args []string) {
	if action, exists := cr.commands[command]; exists {
		action(args)
	}  else if out, err := exec.Command(command, args...).Output(); err == nil {
		fmt.Print(string(out))
	} else {
		fmt.Printf("%s: command not found\n", command)
	}
}

// Sprawdzenie, czy komenda istnieje
func (cr *CommandRegistry) Exists(command string) bool {
	_, exists := cr.commands[command]
	return exists
}

// Obsługuje wyjście z programu
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

// Obsługuje polecenie echo
func handleEcho(args []string) {
	fmt.Println(strings.Join(args, " "))
}

// Obsługuje polecenie type
func handleType(args []string, registry *CommandRegistry) {
	if len(args) == 0 {
		fmt.Println("type: missing argument")
		return
	}

	command := args[0]
	if registry.Exists(command) {
		fmt.Printf("%s is a shell builtin\n", command)
	} else if path, err := exec.LookPath(command); err == nil {
		fmt.Printf("%s is %s\n", command, path)
	} else {
		fmt.Printf("%s: not found\n", command)
	}
}

// Zarządza interakcją z użytkownikiem
type Shell struct {
	registry *CommandRegistry
}

// Nowa instancja shella
func NewShell(registry *CommandRegistry) *Shell {
	return &Shell{registry: registry}
}

// Wczytuje i zwraca komendę użytkownika
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

// Wypisuje ścieżkę 
func handlePwd(args []string) {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(path) 
}

// Uruchamia pętlę shella
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

func main() {
	registry := InitializeRegistry()
	shell := NewShell(registry)
	shell.Run()
}
