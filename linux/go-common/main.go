package main

import (
	"fmt"
	"log"
)

func main() {
	log.SetPrefix("::Testing::")
	log.SetFlags(0)
	log.Print("main()::")
	fmt.Println("Hello, World!!!")
}
