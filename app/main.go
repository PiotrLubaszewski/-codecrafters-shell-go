package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// CommandRegistry manages registering and executing commands
type CommandRegistry struct {
	commands map[string]func([]string)
}

// NewCommandRegistry creates a new command registry
func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{commands: make(map[string]func([]string))}
}

// Register adds a new command to the registry
func (cr *CommandRegistry) Register(name string, handler func([]string)) {
	cr.commands[name] = handler
}

// Execute runs a command with the given arguments
func (cr *CommandRegistry) Execute(command string, args []string) {
	if handler, exists := cr.commands[command]; exists {
		handler(args)
	} else {
		cr.runExternalCommand(command, args)
	}
}

// runExternalCommand attempts to execute a system command
func (cr *CommandRegistry) runExternalCommand(command string, args []string) {
	cmd := exec.Command(command, args...)
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("%s: command not found\n", command)
		return
	}
	fmt.Print(string(output))
}

// Exists checks if a command is registered
func (cr *CommandRegistry) Exists(command string) bool {
	_, exists := cr.commands[command]
	return exists
}

// InitializeRegistry sets up built-in commands
func InitializeRegistry() *CommandRegistry {
	registry := NewCommandRegistry()
	registry.Register("exit", handleExit)
	registry.Register("echo", handleEcho)
	registry.Register("type", func(args []string) { handleType(args, registry) })
	registry.Register("pwd", handlePwd)
	registry.Register("cd", handleCd)
	return registry
}

// Shell handles user interaction
type Shell struct {
	registry *CommandRegistry
}

// NewShell creates a new shell instance
func NewShell(registry *CommandRegistry) *Shell {
	return &Shell{registry: registry}
}

// Run starts the shell loop
func (s *Shell) Run() {
	for {
		command, args, err := s.readCommandAndArgs()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			continue
		}
		if command == "" {
			continue
		}
		s.registry.Execute(command, args)
	}
}

// readCommandAndArgs reads and parses user input
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

	args := parseArgs(scanner.Text())
	if len(args) == 0 {
		return "", nil, nil
	}
	return args[0], args[1:], nil
}

// parseArgs handles quoted arguments properly
func parseArgs(input string) []string {
	var args []string
	var current strings.Builder
	inQuote := false

	for i := 0; i < len(input); i++ {
		ch := input[i]
		switch {
		case ch == '\'':
			inQuote = !inQuote
		case ch == ' ' && !inQuote:
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		default:
			current.WriteByte(ch)
		}
	}
	if current.Len() > 0 {
		args = append(args, current.String())
	}
	return args
}

// handleExit terminates the shell
func handleExit(args []string) {
	exitCode := 0
	if len(args) > 0 {
		if code, err := strconv.Atoi(args[0]); err == nil {
			exitCode = code
		} else {
			fmt.Println("Invalid exit code")
			return
		}
	}
	os.Exit(exitCode)
}

// handleEcho prints the provided text
func handleEcho(args []string) {
	fmt.Println(strings.Join(args, " "))
}

// handleType checks if a command is built-in or in PATH
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

// handlePwd prints the current working directory
func handlePwd(args []string) {
	if path, err := os.Getwd(); err == nil {
		fmt.Println(path)
	} else {
		fmt.Println(err)
	}
}

// handleCd changes the working directory
func handleCd(args []string) {
	if len(args) != 1 {
		fmt.Println("Usage: cd <directory>")
		return
	}

	path := args[0]
	if path == "~" {
		path = os.Getenv("HOME")
	}

	if err := os.Chdir(path); err != nil {
		fmt.Fprintf(os.Stderr, "%s: No such file or directory\n", path)
	}
}

func main() {
	registry := InitializeRegistry()
	shell := NewShell(registry)
	shell.Run()
}
