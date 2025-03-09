package main

import (
	"bufio"
	"fmt"
	"os"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

// Command and their actions
var commands = map[string]func()  {
	"exit 0": func () { os.Exit(0) }, 
}

// Reads and returns inserted user commands
func readCommand() (string, error) {
	fmt.Fprint(os.Stdout, "$ ")
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return "", scanner.Err()
	}
	return strings.TrimSpace(scanner.Text()), nil
}

// If command exists, executes command, else returns error
func executeCommand(command string) {
	if action, exists := commands[command]; exists {
		action()
	} else {
		fmt.Printf("%s: command not found\n", command)
	}
}

func main() {
	// Uncomment this block to pass the first stag
	for {
		command, err := readCommand()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}
		executeCommand(command)
	}
}
