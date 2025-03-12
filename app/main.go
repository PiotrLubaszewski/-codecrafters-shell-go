package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// CommandRegistry manages registering commands
type CommandRegistry struct {
	commands map[string]func([]string)
}

// Creates new Command Registry
func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{commands: make(map[string]func([]string))}
}

// Initializes Command Registry with initial (builtin) values
func InitializeRegistry() *CommandRegistry {
	registry := NewCommandRegistry()
	registry.Register("exit", handleExit)
	registry.Register("echo", handleEcho)
	registry.Register("type", func(args []string) { handleType(args, registry) })
	registry.Register("pwd", handlePwd)
	registry.Register("cd", handleCd)
	return registry
}

// Adds new command to Command Registry
func (cr *CommandRegistry) Register(name string, handler func([]string)) {
	cr.commands[name] = handler
}

// Executes command with given args
func (cr *CommandRegistry) Execute(command string, args []string) {
	if action, exists := cr.commands[command]; exists {
		action(args)
	} else if out, err := exec.Command(command, args...).Output(); err == nil {
		fmt.Print(string(out))
	} else {
		fmt.Printf("%s: command not found\n", command)
	}
}

// Checks if command exists in registry
func (cr *CommandRegistry) Exists(command string) bool {
	_, exists := cr.commands[command]
	return exists
}

// Exits shell
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

// Prints provided text to output
func handleEcho(args []string) {
	fmt.Println(strings.Join(args, " "))
}

// Prints if command is builtin, PATH-vide, or not found
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

// Manages Contatct with user
type Shell struct {
	registry *CommandRegistry
}

// Creates new shell instance
func NewShell(registry *CommandRegistry) *Shell {
	return &Shell{registry: registry}
}

// Reads command and args from user input
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

	input := scanner.Text()
	var args []string
	var current strings.Builder
	inQuote := false

	for i := 0; i < len(input); i++ {
		ch := input[i]

		if ch == '\'' {
			if inQuote {
				inQuote = false
			} else {
				inQuote = true
			}
			continue
		}

		if ch == ' ' && !inQuote {
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
			continue
		}

		current.WriteByte(ch)
	}

	if current.Len() > 0 {
		args = append(args, current.String())
	}

	if len(args) == 0 {
		return "", nil, nil
	}

	return args[0], args[1:], nil
}

// Changes working direcotory
func handleCd(args []string) {
	if len(args) != 1 {
		fmt.Println("String not in pwd: $s", strings.Join(args, " "))
	}

	path := args[0]

	if strings.TrimSpace(path) == "~" {
		path = os.Getenv("HOME")
	}

	if err := os.Chdir(path); err != nil {
		fmt.Fprintf(os.Stdout, "%s: No such file or directory\n", path)
	}
}

// Executes Shell
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
