
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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
	fmt.Println(msg)
}

// Type command checks for builtin commands and unrecognized commands 
func handleType(args []string, commands map[string]func([]string)) {
	if len(args) == 0 {
		fmt.Println("type: missing argument")
		return
	}

	command := args[0]
	var msg string

	if checkCommand(command, commands) {
		msg = fmt.Sprintf("%s is a shell builtin", command)
	} else {
		msg = fmt.Sprintf("%s: not found", command)
	}

	fmt.Println(msg)
}

// Reads and returns inserted user commands
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

// If command exists, executes command, else returns error
func executeCommand(command string, args []string, commands map[string]func([]string)) {
	if action, exists := commands[command]; exists {
		action(args)
	} else {
		fmt.Printf("%s: command not found\n", command)
	}
}

func checkCommand(command string, commands map[string]func([]string)) bool {
	_, exists := commands[command]
	return exists
}

func main() {
	commands := map[string]func([]string){
		"exit": handleExit,
		"echo": handleEcho,
		"type": func(args []string) { handleType(args, commands) },
	}

	for {
		command, args, err := readCommandAndArgs()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}
		if command == "" {
			continue
		}
		executeCommand(command, args, commands)
	}
}
