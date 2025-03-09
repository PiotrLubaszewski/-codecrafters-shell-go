package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Command and their actions
var commands = map[string]func(args []string){
	"exit": handleExit,
	"echo": handleEcho,
}

// Terminates with code/status 0.
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

// Echo command prints the provided text back.
func handleEcho(args []string) {
	msg := strings.Join(args, " ")
	fmt.Fprintln(os.Stdout, msg)
}

// Reads and returns inserted user commands
func readCommandAndArgs() (string, []string, error) {
	fmt.Fprint(os.Stdout, "$ ")
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

// If command exists, executes command, else returns error
func executeCommand(command string, args []string) {
	if action, exists := commands[command]; exists {
		action(args)
	} else {
		fmt.Printf("%s: command not found\n", command)
	}
}

func main() {
	for {
		command, args, err := readCommandAndArgs()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}
		if command == "" {
			continue
		}
		executeCommand(command, args)
	}
}
