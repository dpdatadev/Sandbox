package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

//TODO, BETA - use actual unit testing, bdd, mocks, etc.,
//Replacing all these with real unit tests is gonna suck (2/21)

// TESTING
const (
	LOCAL_PARSE_TXT_FILE   = "proc.txt"
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

func setupTestDatabase(databaseName string, tableSql string /*,overwrite bool*/) (*sql.DB, error) {
	//check if file already exists on file system - if so delete if flag is set
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

func getTestDataBase(databaseName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", fmt.Sprintf("./%s.db", databaseName))
	if err != nil {
		log.Fatal(err)
	}
	return db, nil
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

// TODO
type TestCommand struct {
	Name        string
	Args        []string
	Description string
}

func setupConsoleCommandTestFromFileSuite() (CommandRunner, []*Command) {
	var CmdIOHelper CmdIOHelper

	commands := CmdIOHelper.ParseCommands(LOCAL_PARSE_TXT_FILE)

	if len(commands) == 0 {
		PrintFailure("No commands parsed from file, check for errors and ensure proc.txt is in the correct location with valid commands")
	}

	consoleCommandRunner := NewCommandRunner(RunnerType_Console)

	return consoleCommandRunner, commands
}

func setupConsoleCommandTestSuite() (CommandRunner, []*Command) {
	// I could just as well call NewCommandService to return a RemoteSSH executor to download logs into a SQLITE Store
	//svc := NewCommandService(sqlStore, sshExec)

	hostInfo := NewCommand("uname", []string{"-a"}, "Local Host Info")

	cmd := NewCommand("ip", []string{"addr"}, "Get Local IP Config") // TODO, deal with default no args

	cmd1 := NewCommand("ip", []string{"neighbor"}, "Get IP Neighbor Output")

	cmd2 := NewCommand("ifconfig", []string{}, "Get another local IP config output")

	//TEST BAD ARGS
	cmd3 := NewCommand("ifconfig", []string{""}, "Intentionally fail")

	cmd4 := NewCommand("free", []string{"-g", "-h"}, "Get Active Memory Usage")

	cmd5 := NewCommand("arp", []string{"-a"}, "Get Local ARP Cache")

	//TEST SECURITY POLICY
	cmd6 := NewCommand("sudo", []string{"dd"}, "NAUGHTY COMMAND")

	commands := []*Command{hostInfo, cmd, cmd1, cmd2, cmd3, cmd4, cmd5, cmd6}

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

	var CmdIOHelper CmdIOHelper

	for _, cmd := range testCommands {
		CmdIOHelper.ConsoleDump(cmd)
	}
}

func ConsoleSqliteCommandTest(databaseName string, tableSQL string) {
	ctx, cancel := setupTimeoutContext()

	defer cancel()

	testDb, _ := setupTestDatabase(databaseName, tableSQL)

	defer testDb.Close()

	svc := setupLocalSqliteCommandService(testDb)

	consoleCommandRunner, commands := setupConsoleCommandTestSuite()

	testCommands := consoleCommandRunner.RunCommands(svc, ctx, commands, true)

	var CmdIOHelper CmdIOHelper

	for _, cmd := range testCommands {
		CmdIOHelper.ConsoleDump(cmd)
	}
}

func ConsoleSqliteCommandFileTest(databaseName string, tableSQL string) {
	ctx, cancel := setupTimeoutContext()

	defer cancel()

	testDb, _ := setupTestDatabase(databaseName, tableSQL)

	defer testDb.Close()

	svc := setupLocalSqliteCommandService(testDb)

	consoleCommandRunner, commands := setupConsoleCommandTestFromFileSuite()

	testCommands := consoleCommandRunner.RunCommands(svc, ctx, commands, true)

	var CmdIOHelper CmdIOHelper

	for _, cmd := range testCommands {
		CmdIOHelper.ConsoleDump(cmd)
	}
}

func testGetAllCommands() {
	db, _ := getTestDataBase(LOCAL_SQLITE_CMD_DB1)
	defer db.Close()

	s := NewSqliteCommandStore(db)
	ctx, cancel := setupTimeoutContext()
	defer cancel()

	cmds, err := s.GetAll(ctx)

	if err != nil {
		log.Fatal(err)
	}

	for _, cmd := range cmds {
		PrintDebug(fmt.Sprintf("Command: %s, Status: %s\n", cmd.Name, cmd.Status))
	}
}

func testGetRecentCommands() {
	db, _ := getTestDataBase(LOCAL_SQLITE_CMD_DB1)
	defer db.Close()

	s := NewSqliteCommandStore(db)
	ctx, cancel := setupTimeoutContext()
	defer cancel()

	cmds, err := s.GetRecent(ctx, 5)

	if err != nil {
		log.Fatal(err)
	}

	for _, cmd := range cmds {
		PrintDebug(fmt.Sprintf("Command: %s, Status: %s\n", cmd.Name, cmd.Status))
	}
}

func singleListTest() {
	s := new(SList[Command])
	s.Value = Command{ID: uuid.New(), Args: []string{}, Notes: "test!"}

	cmd := NewCommand("ip", []string{"neighbor"}, "test")
	cmd1 := NewCommand("ip", []string{"addr"}, "test")

	s.Append(cmd)
	s.Append(cmd1)

	PrintDebug("Current Command: %s\n", s.Values()[0].Name)
	PrintDebug("Next Command: %s\n", s.ForwardNode().Value.Name)

	PrintDebug("SList length: %d\n", s.Len())
	PrintDebug("SList values: %v\n", s.Values())

	s.Print()

}

func lineageTest() {

	ctx, cancel := setupTimeoutContext()
	defer cancel()

	testDb, _ := getTestDataBase(LOCAL_SQLITE_CMD_DB1)
	defer testDb.Close()

	store := NewSqliteCommandStore(testDb)

	recentCommands, err := store.GetRecent(ctx, 5)
	if err != nil {
		log.Fatal(err)
	}

	hs := &HistoryService{
		AuditCommands: recentCommands,
	}

	lineageObjects := hs.BeginChain()
	err = hs.LinkChain(lineageObjects)

	if err != nil {
		PrintStdErr("Error linking lineage: %v", err)
		return
	}

	for _, obj := range lineageObjects {
		fmt.Printf("ID: %s, BatchID: %s, PrevID: %v, Status: %s, NextID: %v, RootID: %v\n",
			obj.ID, obj.BatchID, obj.PrevID, obj.Status, obj.NextID, obj.RootID)
	}

	WriteLineageToFile(lineageObjects, "chain_3.txt")
}

// Testing
func mainTestSuite() {
	log.SetPrefix("::Test Runs::")
	log.SetFlags(0)
	log.Print("mainTestSuite()::")
	ConsoleSqliteCommandFileTest(LOCAL_SQLITE_CMD_DB5, LOCAL_SQLITE_CMD_TABLE)
	//ConsoleInMemoryCommandTest()
	//ConsoleSqliteCommandTest(LOCAL_SQLITE_CMD_DB1, LOCAL_SQLITE_CMD_TABLE)
}

func runTests() {
	//log.SetPrefix("::APP::")
	//log.SetFlags(0)
	//log.Print("main()::")
	//log.Println("DPDIGITAL,LLC::COMMANDER::<INIT>::")
	mainTestSuite()
}

//Feb week 3
//SQLITE impl (complete)
//Default args (complete)
//Client code with baseline v0 API
//Simple http/net or udp exposure

// TODO - phase 2:
// Lineage/history (alpha complete)
//Add Priority Queue capability? (cmd 1, priority 1, cmd 2, priority -1, etc.,)
// Piping multiple commands (still local)
// Get out of dev zone and create actual package structure (started in alpha..)

// TODO - phase 3:
// Working with HTTP(S)
// Working with SSH
// Test Scrubers/Security (complete for alpha, but needs more work)
// Metadata/Analysis/Examples
// LAUNCH
