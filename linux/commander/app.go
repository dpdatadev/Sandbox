package main

import (
	"fmt"
	"log"
)

func boostrapApp() {
	panic("not impl")
}

func retrieveCommands() *[]Command {
	panic("not impl")
}

func main() {
	log.SetPrefix("::APP::")
	log.SetFlags(0)
	log.Print("main()::")
	fmt.Println("DPDIGITAL,LLC::COMMANDER::<INIT>::")
	boostrapApp()
}
