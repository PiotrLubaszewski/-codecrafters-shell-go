package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

// Command and their actions
var commands = map[string]func(args []string)  {
	"exit": handleExit, 
	"echo": handleEcho, 
}

// Terminates with code/status 0.
func handleExit(args []string) {
	if len(args) > 1 {
		exit_code, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Invalid exit code")
			return
		}
		os.Exit(exit_code)
	}
	return
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
		return "", scanner.Err()
	}
	trimed := strings.TrimSpace(scanner.Text())
	splited := strings.Split(trimed, " ")
	return splitted[0], args[1:], nil
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
		command, args, err := readCommand()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}
		executeCommand(command, args)
	}
}
