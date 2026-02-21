package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Printf("Begin: %s\n", os.Args[0])
	file, err := os.Open("proc.txt")
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	// Process the file (for demonstration, we'll just read it)
	buf := make([]byte, 1024)
	n, err := file.Read(buf)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}
	commands := strings.SplitSeq(string(buf[:n]), "\n")
	for cmd := range commands {
		cmds := strings.Fields(cmd)
		cmdName := cmds[0]
		cmdArgs := cmds[1:]

		fmt.Printf("Command: %s, Args: %v\n", cmdName, cmdArgs)
	}
}
