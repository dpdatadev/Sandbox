package main

//TODO

import (
	"errors"
	"fmt"
	"log"
	"testing"
)

// DEV: 2/9
// ////////////////////////////////
func leftTest(s string, size int) (string, error) {

	if s == "" {
		return s, errors.New("EMPTY STRING")
	}

	leftSubstr := s[:size]

	return leftSubstr, nil
}

// ////////////////////////////////
// Tests

func MatchString(s1 string, s2 string) bool {
	return s1 == s2
}

func testSubString(t *testing.T) {
	testString := "TestString"
	expectedString := "Tes"
	want, err := Left(testString, 3)

	if !MatchString(expectedString, want) || err != nil {
		t.Errorf("testSubString() FAILED = %q, %v", want, err)
	}
}

func testMain() {
	log.SetPrefix("::Testing::")
	log.SetFlags(0)
	log.Print("main()::")
	fmt.Println("Hello, World!!!")
}
