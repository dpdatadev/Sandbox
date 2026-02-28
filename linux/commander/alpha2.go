package main

import (
	"context"
	"fmt"
	"log"
	"os/user"
	"time"

	hub "github.com/goforj/execx"
)

func GetUser() string {
	current_user, err := user.Current()
	if err != nil {
		PrintStdErr("USER LOOKUP OCCURRED: ", err)
	}

	username := current_user.Username

	return username
}

func execxTest() {
	// Run executes the command and returns the result and any error.

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := hub.
		Command("printf", "hello\nworld\n").
		Pipe("tr", "a-z", "A-Z").
		Env("MODE=demo").
		WithContext(ctx).
		OnStdout(func(line string) {
			fmt.Println("OUT:", line)
		}).
		OnStderr(func(line string) {
			fmt.Println("ERR:", line)
		}).
		Run()

	if !res.OK() {
		log.Fatalf("command failed: %v", err)
	}

	fmt.Printf("Stdout: %q\n", res.Stdout)
	fmt.Printf("Stderr: %q\n", res.Stderr)
	fmt.Printf("ExitCode: %d\n", res.ExitCode)
	fmt.Printf("Error: %v\n", res.Err)
	fmt.Printf("Duration: %v\n", res.Duration)
	// OUT: HELLO
	// OUT: WORLD
	// Stdout: "HELLO\nWORLD\n"
	// Stderr: ""
	// ExitCode: 0
	// Error: <nil>
	// Duration: 10.123456ms
}

func main() {
	execxTest()
}
