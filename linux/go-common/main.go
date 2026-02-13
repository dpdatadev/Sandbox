package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"runtime/debug"
	"sync"
	"time"

	"github.com/google/uuid"
)

// DEV: 2/9
// ////////////////////////////////

func printAlloc() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Allocated Heap: %v MB\n", m.Alloc/1024/1024)
}

// ANSI SQL LEFT style substring
func Left(s string, size int) (string, error) {

	if s == "" {
		return s, errors.New("EMPTY STRING")
	}

	leftSubstr := s[:size]

	return leftSubstr, nil
}

// ANSI SQL RIGHT style substring
func Right(s string, size int) (string, error) {
	if s == "" {
		return s, errors.New("EMPTY STRING")
	}

	appliedSize := max((len(s) - size), 0)

	return s[appliedSize:], nil
}

// Version 4 Google UUID (length 7) (UNSAFE, INTERNAL USE ONLY)
func NewShortUUID() (string, error) {

	uuidString, err := Left(uuid.NewString(), 8)

	return uuidString, err
}

func CommandTest() (string, error) {

	if runtime.GOOS == "windows" {
		log.Panicf("WINDOWS NOT SUPPORTED")
		os.Exit(-1)
	}

	out, err := exec.Command("ip", "neighbor").Output()
	if err != nil {
		log.Panicf("%s", err)
	}

	log.Print("Test Command (CMD) Successfully Executed")
	output := string(out[:])
	return output, nil
}

func CommandTestRunner(
	svc *CommandService,
	ctx context.Context,
	cmds []*Command,
) []*Command {

	var wg sync.WaitGroup

	finished := make([]*Command, len(cmds))

	for i, cmd := range cmds {
		wg.Add(1)

		go func(i int, cmd *Command) {
			defer wg.Done()

			if err := svc.Run(ctx, cmd); err != nil {
				panic(err) //todo remove all panics from framework code
			}

			finished[i] = cmd
		}(i, cmd)
	}

	wg.Wait()
	return finished
}

func CommandSystemTest() {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	store := NewInMemoryStore()
	exec := NewLocalExecutor()

	svc := NewCommandService(store, exec)

	cmd := NewCommand(
		"echo",
		[]string{"Hello from Phase 1"},
		"test command",
	)

	cmd1 := NewCommand("ip", []string{"neighbor"}, "IP Test Command")

	commands := []*Command{cmd, cmd1}

	commands = append(commands, NewCommand("arp", []string{"-a"}, "Arp Test Command"))

	testCommands := CommandTestRunner(svc, ctx, commands)

	for _, cmd := range testCommands {
		fmt.Println("Command ID:", cmd.ID)
		fmt.Println("Status:", cmd.Status)
		fmt.Println("Stdout:", cmd.Stdout)
	}
}

// ////////////////////////////////
// Create actual unit tests .. TODO
func main() {
	fmt.Printf("%s", debug.Stack())
	log.SetPrefix("::Testing::")
	log.SetFlags(0)
	log.Print("main()::")
	fmt.Println("Hello, World!!!")

	_, hostname := os.Hostname()
	pid := os.Getpid()
	log.Println(hostname)
	log.Println(pid)

	testString := "Hello Mate"
	leftString, err := Left(testString, 7)
	if err == nil {
		fmt.Println(leftString)
	}

	rightString, err1 := Right(testString, 3)
	if err1 == nil {
		fmt.Println(rightString)
	}

	newID, err := NewShortUUID()
	log.Println("::UUID SERVICE Start::")
	fmt.Println(newID)

	defer StartServer()
	//defer printAlloc() //see heap usage after we force GC towards the end
	log.Println("::DB SERVICE START::")
	CreateDB()

	uuid, err := NewShortUUID()
	if err != nil {
		panic("UUID ERROR")
	} else {
		log.Println("::SEEDB SERVICE START::")
		SeedDB(uuid, "TEST NOTES")
	}

	log.Println("::QUERY DB READ STATE START::")
	QueryDBTest()

	//commandTest, _ := CommandTest()

	CommandSystemTest()
	log.Println("::HTTP SERVICE START -- UP -- ::")
	//printAlloc()
	//runtime.GC()
}
