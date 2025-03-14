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

	args := ParseArgs(scanner.Text())
	if len(args) == 0 {
		return "", nil, nil
	}
	return args[0], args[1:], nil
}

// Parser is responsible for parsing input arguments
// It adheres to the Single Responsibility Principle by focusing solely on parsing
// and not handling execution logic or external dependencies.
type Parser struct {
	input    string
	args     []string
	current  strings.Builder
	stack    []rune
	escaped  bool
	handlers map[rune]func(*Parser)
}

// NewParser creates a new instance of the parser
// This follows the Open/Closed Principle by allowing extension without modification.
func NewParser(input string) *Parser {
	parser := &Parser{input: input}
	parser.handlers = map[rune]func(*Parser){
		'\\': func(p *Parser) { p.escaped = true },
		'\'': func(p *Parser) { p.toggleQuote('\'') },
		'"':  func(p *Parser) { p.toggleQuote('"') },
		' ':  func(p *Parser) { p.handleSpace() },
	}
	return parser
}

// Parse executes the parsing operation
func (p *Parser) Parse() []string {
	for i := 0; i < len(p.input); i++ {
		p.processRune(rune(p.input[i]))
	}
	p.addCurrentArg()
	return p.args
}

// processRune processes a single input character
// This method ensures the correct sequence of parsing operations.
func (p *Parser) processRune(ch rune) {
	if p.escaped {
		p.current.WriteRune(ch)
		p.escaped = false
		return
	}

	if handler, exists := p.handlers[ch]; exists {
		handler(p)
	} else {
		p.current.WriteRune(ch)
	}
}

// toggleQuote manages opening and closing quotes
func (p *Parser) toggleQuote(ch rune) {
	if len(p.stack) > 0 && p.stack[len(p.stack)-1] == ch {
		p.stack = p.stack[:len(p.stack)-1]
	} else {
		p.stack = append(p.stack, ch)
	}
}

// handleSpace handles argument separation
func (p *Parser) handleSpace() {
	if len(p.stack) == 0 && p.current.Len() > 0 {
		p.addCurrentArg()
	}
}

// addCurrentArg adds the current argument to the list
func (p *Parser) addCurrentArg() {
	if p.current.Len() > 0 {
		p.args = append(p.args, p.current.String())
		p.current.Reset()
	}
}

// ParseArgs is a helper function for ease of use
func ParseArgs(input string) []string {
	parser := NewParser(input)
	return parser.Parse()
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
