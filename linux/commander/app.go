package main

import (
	"fmt"
	"log"
)

func runTests() {
	mainTestSuite()
	//singleListTest()
	//lineageTest()
}

func main() {
	log.SetPrefix("::APP::")
	log.SetFlags(0)
	log.Print("main()::")
	fmt.Println("DPDIGITAL,LLC::COMMANDER::<INIT>::")
	runTests()
}
