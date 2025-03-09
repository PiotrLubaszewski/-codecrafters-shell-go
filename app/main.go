package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// CommandRegistry zarządza rejestracją i wykonywaniem komend
type CommandRegistry struct {
	commands map[string]func([]string)
}

// Nowy rejestr komend
func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{commands: make(map[string]func([]string))}
}

// Rejestracja nowej komendy
func (cr *CommandRegistry) Register(name string, handler func([]string)) {
	cr.commands[name] = handler
}

// Wykonanie komendy
func (cr *CommandRegistry) Execute(command string, args []string) {
	if action, exists := cr.commands[command]; exists {
		action(args)
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
	if len(args) > 0 {
		exitCode, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Invalid exit code")
			return
		}
		os.Exit(exitCode)
	} else {
		os.Exit(0)
	}
}

// Obsługuje polecenie echo
func handleEcho(args []string) {
	msg := strings.Join(args, " ")
	fmt.Println(msg)
}

// Obsługuje polecenie type
func handleType(args []string, registry *CommandRegistry) {
	if len(args) == 0 {
		fmt.Println("type: missing argument")
		return
	}

	command := args[0]
	var msg string

	if registry.Exists(command) {
		msg = fmt.Sprintf("%s is a shell builtin", command)
	} else {
		msg = fmt.Sprintf("%s: not found", command)
	}

	fmt.Println(msg)
}

// Wczytuje i zwraca komendę użytkownika
func readCommandAndArgs() (string, []string, error) {
	fmt.Print("$ ")
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return "", nil, scanner.Err()
	}

	trimmed := strings.TrimSpace(scanner.Text())
	splitted := strings.Fields(trimmed)
	if len(splitted) == 0 {
		return "", nil, nil
	}

	return splitted[0], splitted[1:], nil
}

func main() {
	registry := NewCommandRegistry()

	registry.Register("exit", handleExit)
	registry.Register("echo", handleEcho)
	registry.Register("type", func(args []string) { handleType(args, registry) })

	for {
		command, args, err := readCommandAndArgs()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}
		if command == "" {
			continue
		}
		registry.Execute(command, args)
	}
}
