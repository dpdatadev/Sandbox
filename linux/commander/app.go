package main

import (
	"log"
)

func runTests() {
	//mainTestSuite()
	//singleListTest()
	lineageTest()
	//testGetAllCommands()
	//testGetRecentCommands()
}

func main() {
	var ioHelper IoHelper
	log.SetPrefix("::APP::")
	log.SetFlags(0)
	log.Print("main()::")
	log.Println("DPDIGITAL,LLC::COMMANDER::<INIT>::")
	ipInfo, _ := ioHelper.getHostIpConfig()
	log.Println(ipInfo)
	//runTests()
}
