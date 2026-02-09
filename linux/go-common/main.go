package main

import (
	"errors"
	"fmt"
	"log"
	"runtime/debug"

	"github.com/google/uuid"
)

// DEV: 2/9
// ////////////////////////////////

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

	appliedSize := (len(s) - size)

	if appliedSize < 0 {
		appliedSize = 0
	}

	return s[appliedSize:], nil
}

// Version 4 Google UUID (length 7) (UNSAFE, INTERNAL USE ONLY)
func NewShortUUID() (string, error) {

	uuidString, err := Left(uuid.NewString(), 8)

	return uuidString, err
}

// ////////////////////////////////

func main() {
	fmt.Printf("%s", debug.Stack())
	log.SetPrefix("::Testing::")
	log.SetFlags(0)
	log.Print("main()::")
	fmt.Println("Hello, World!!!")

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

	fmt.Println(newID)
}
