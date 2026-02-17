package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"runtime/debug"
	"time"
)

// TESTING
const (
	LOCAL_SQLITE_CMD_DB1   = "testcmd1"
	LOCAL_SQLITE_CMD_DB2   = "testcmd2"
	LOCAL_SQLITE_CMD_DB3   = "testcmd3"
	LOCAL_SQLITE_CMD_DB4   = "testcmd4"
	LOCAL_SQLITE_CMD_DB5   = "testcmd5"
	LOCAL_SQLITE_CMD_TABLE = `CREATE TABLE IF NOT EXISTS commands (
    id	 		  INTEGER PRIMARY KEY AUTOINCREMENT,
	uuid          TEXT,
    name          TEXT NOT NULL,
    status        TEXT NOT NULL,
    created_at    DATETIME NOT NULL,
    started_at    DATETIME,
    finished_at   DATETIME,
    stdout        TEXT,
    stderr        TEXT,
    exit_code     INTEGER
);`
)

//create indexes(?)

func setupTestDatabase(databaseName string, tableSql string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", fmt.Sprintf("./%s.db", databaseName))
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(tableSql) //inject me baby
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("New table on %s created successfully", databaseName)
	PrintDebug("SQL EXECUTED on %s:::\n%s:::", databaseName, tableSql)
	return db, err //defer close!
}

func setupLocalInMemoryCommandService() *CommandService {
	store := NewInMemoryStore()
	//store := NewSqliteStore()
	exec := NewLocalExecutor()
	//store := NewSSHStore()

	//Here we will use the default(Memory Store, Local Execution)
	svc := NewCommandService(store, exec)

	return svc
}

func setupLocalSqliteCommandService(database *sql.DB) *CommandService {

	//defer testDb.Close()

	store := NewSqliteCommandStore(database)

	exec := NewLocalExecutor()

	svc := NewCommandService(store, exec)

	return svc
}

func setupTimeoutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
}

func setupConsoleCommandTestSuite() (CommandRunner, []*Command) {
	// I could just as well call NewCommandService to return a RemoteSSH executor to download logs into a SQLITE Store
	//svc := NewCommandService(sqlStore, sshExec)

	hostInfo := NewCommand("uname", []string{"-a"}, "Local Host Info")

	cmd := NewCommand("ifconfig", []string{""}, "Get Local NIC Config") // TODO, deal with default no args

	cmd1 := NewCommand("ip", []string{"neighbor"}, "Get IP Neighbor Output")

	cmd2 := NewCommand("free", []string{"-g", "-h"}, "Get Active Memory Usage")

	cmd3 := NewCommand("arp", []string{"-a"}, "Get Local ARP Cache")

	cmd4 := NewCommand("sudo", []string{"dd"}, "NAUGHTY COMMAND")

	commands := []*Command{hostInfo, cmd, cmd1, cmd2, cmd3, cmd4}

	consoleCommandRunner := NewCommandRunner(RunnerType_Console)

	return consoleCommandRunner, commands
}

// Orchestrate sample commands with ConsoleRunner for Testing
func ConsoleInMemoryCommandTest() {

	ctx, cancel := setupTimeoutContext()

	defer cancel()

	svc := setupLocalInMemoryCommandService()

	consoleCommandRunner, commands := setupConsoleCommandTestSuite()

	testCommands := consoleCommandRunner.RunCommands(svc, ctx, commands, true)

	ioHelper := &IoHelper{}

	for _, cmd := range testCommands {
		ioHelper.ConsoleDump(cmd)
	}
}

func ConsoleSqliteCommandTest(databaseName string, tableSQL string) {
	ctx, cancel := setupTimeoutContext()

	defer cancel()

	testDb, err := setupTestDatabase(databaseName, tableSQL)

	if err != nil {
		log.Panicf("DB ERROR %v", err)
	}

	defer testDb.Close()

	svc := setupLocalSqliteCommandService(testDb)

	consoleCommandRunner, commands := setupConsoleCommandTestSuite()

	testCommands := consoleCommandRunner.RunCommands(svc, ctx, commands, true)

	ioHelper := &IoHelper{}

	for _, cmd := range testCommands {
		ioHelper.ConsoleDump(cmd)
	}
}

// Testing
func main() {
	fmt.Printf("%s\n", debug.Stack())
	log.SetPrefix("::TEST::")
	log.SetFlags(0)
	log.Print("main()::")
	fmt.Println("TESTRUNNER::<INIT>::")
	//ConsoleInMemoryCommandTest()
	ConsoleSqliteCommandTest(LOCAL_SQLITE_CMD_DB1, LOCAL_SQLITE_CMD_TABLE)
}

//Feb week 3
//SQLITE impl
//Default args
//Client code with baseline v0 API
//Simple http/net or udp exposure

// TODO - phase 2:
// Lineage/history
// Piping multiple commands (still local)
// Get out of dev zone and create actual package structure

// TODO - phase 3:
// Working with HTTP(S)
// Working with SSH
// Test Scrubers/Security
// Metadata/Analysis/Examples
// LAUNCH
