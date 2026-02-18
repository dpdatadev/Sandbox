package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"

	"github.com/fatih/color"
	"github.com/google/uuid"
)

var (
	PrintIdentity = color.New(color.Bold, color.FgGreen, color.Italic).PrintfFunc()
	PrintSuccess  = color.New(color.Bold, color.FgGreen, color.Underline).PrintfFunc()
	PrintStdOut   = color.New(color.Bold, color.FgYellow).PrintfFunc()
	PrintStdErr   = color.New(color.Bold, color.FgHiRed).PrintfFunc()
	PrintFailure  = color.New(color.Bold, color.FgRed, color.Underline).PrintfFunc()
	PrintDebug    = color.New(color.Bold, color.FgBlue, color.Italic).PrintfFunc()
)

type IoHelper struct{}

func (io *IoHelper) showHostIpConfig() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal(err)
	}

	for _, addr := range addrs {
		// Check the address type and if it is not a loopback
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				PrintDebug("IPv4 Address: %s\n", ipnet.IP.String())
			} else if ipnet.IP.To16() != nil {
				PrintDebug("IPv6 Address: %s\n", ipnet.IP.String())
			}
		}
	}
}

func (io *IoHelper) printHeap() {
	m := &runtime.MemStats{}
	go runtime.ReadMemStats(m)
	PrintDebug("Allocated Heap: %v MB\n", m.Alloc/1024/1024)
}

// ANSI SQL LEFT style substring
func (io *IoHelper) Left(s string, size int) (string, error) {

	if s == "" {
		return s, errors.New("EMPTY STRING")
	}

	leftSubstr := s[:size]

	return leftSubstr, nil
}

// ANSI SQL RIGHT style substring
func (io *IoHelper) Right(s string, size int) (string, error) {
	if s == "" {
		return s, errors.New("EMPTY STRING")
	}

	appliedSize := max((len(s) - size), 0)

	return s[appliedSize:], nil
}

// Return files for Logging or dumping
func (io *IoHelper) GetFile(fileName string) *os.File {
	if fileName == "" {
		PrintFailure("errors.New(\"\"): %v\n", errors.New("FILE ERROR"))
	}

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		PrintFailure("errors.New(\"\"): %v\n", err)
	}

	return file
}

// Version 4 Google UUID (length 7) (UNSAFE, INTERNAL USE ONLY)
func (io *IoHelper) NewShortUUID() (string, error) {

	uuidString, err := io.Left(uuid.NewString(), 8)

	return uuidString, err
}

// Helper function for displaying/dumping Command info (default Console/Text/Printf())
func (io *IoHelper) ConsoleDump(cmd *Command) {
	if cmd.Stderr != "" || cmd.Status == "FAILED" {
		PrintFailure("Command ID: %v\n", cmd.ID)
		PrintFailure("Command Name: %s\n", cmd.Name)
		PrintFailure("Command Args: %s\n", cmd.Args)
		PrintFailure("Status: %v\n", cmd.Status)
		PrintStdErr("STDERR: %s::<%s>\n", cmd.Stderr, cmd.Error)
		//ConsoleStdErrHandle(cmd.Stderr) //TODO
	} else if cmd.Status == "SUCCESS" {
		PrintIdentity("\nCommand ID: %v\n", cmd.ID)
		PrintIdentity("Command Name: %s\n", cmd.Name)
		PrintIdentity("Command Args: %s\n", cmd.Args)
		PrintSuccess("Status: %v\n", cmd.Status)
		PrintStdOut("STDOUT:\n %s\n", cmd.Stdout)
		fmt.Println()
		//ConsoleStdOutHandle(cmd.Stdout) //TODO
	} else {
		fmt.Println(fmt.Errorf("UNKNOWN ERROR OCCURRED: %s", cmd.ID.String()))
	}
}

func (io *IoHelper) FileDump(cmd *Command, logFile string) {

	log.SetOutput(io.GetFile(logFile))

	if cmd.Stderr != "" || cmd.Status == "FAILED" {
		log.Fatalf("Command ID: %v\n", cmd.ID)
		log.Fatalf("Command Name: %s\n", cmd.Name)
		log.Fatalf("Command Args: %s\n", cmd.Args)
		log.Fatalf("Status: %v\n", cmd.Status)
		log.Fatalf("STDERR: %s::<%s>\n", cmd.Stderr, cmd.Error)
		//ConsoleStdErrHandle(cmd.Stderr) //TODO
	} else if cmd.Status == "SUCCESS" {
		log.Printf("Command ID: %v\n", cmd.ID)
		log.Printf("Command Name: %s\n", cmd.Name)
		log.Printf("Command Args: %s\n", cmd.Args)
		log.Printf("Status: %v\n", cmd.Status)
		log.Printf("STDOUT:\n %s\n", cmd.Stdout)
		//ConsoleStdOutHandle(cmd.Stdout) //TODO
	} else {
		fmt.Println(fmt.Errorf("UNKNOWN ERROR OCCURRED: %s", cmd.ID.String()))
	}
}
