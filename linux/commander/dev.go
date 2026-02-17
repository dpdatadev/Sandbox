// Single file view of all code, docs, examples, and dependencies; to later be split into production module/package(s)
package main //package commander
//https://go.dev/blog/using-go-modules
import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"maps"
	"os"
	"os/exec"
	"runtime"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

/*

COMMANDER (?) still picking the name
https://chatgpt.com/c/698c0190-d8ec-832d-8aee-537b6c64320d

February 2026 - R&D

Containerized testing:
https://pkg.go.dev/github.com/testcontainers/testcontainers-go
https://medium.com/tiket-com/go-with-cucumber-an-introduction-for-bdd-style-integration-testing-7aca2f2879e4

We can setup containers, configured to our liking, describe the scenarios exactly with Cucumber,
then test each of our pipeline components. SSH, HTTP, SQLITE, YAML, etc., and not have a messy server circus.

https://chatgpt.com/c/698c0190-d8ec-832d-8aee-537b6c64320d
https://tutorialedge.net/golang/executing-system-commands-with-golang/
Paying homage to https://github.com/alexellis/go-execute
Notes, 2/11 - thinking of creating a new command API/service for Linux
servers. To have diagnostic and automated activities embedded in various applications.
For use with my raspberry pi's in the beginning. Do more research on this**.
https://chatgpt.com/c/698c0190-d8ec-832d-8aee-537b6c64320d

Chain of Command
github/dpdatadev/chain-of-command
“Composable command execution framework with persistence, security controls, and pipeline orchestration.”
REDIS BACKED, HTTP2 streaming, ACID compliant, process management and remote execution with policy enforcement. (??)
1) Command identity + audit trail

Each command becomes a first-class entity:

2) Persistence layer

Your CommandStore abstraction enables:

SQLite → local dev

Postgres → production

Redis → caching / queues

This converts ephemeral execution into:

Durable execution history

Companies building:

CI/CD tools

Remote agents

Fleet orchestration

Security audit systems

…all need this.

3) Security scrubber / policy engine

os/exec will happily run:

rm -rf /
sudo shutdown now


Your scrubber introduces:

Blocklists

Regex policies

Allowlists

Role-based execution

Now your framework becomes viable for:

SaaS agents

Remote automation

Multi-tenant systems


Yes, os/exec supports pipes…

…but only at the file descriptor level.

You’re abstracting at the semantic level:

sshCmd.
  Pipe(textCmd).
  Pipe(fileCmd).
  Pipe(httpCmd).
  Execute(ctx)


  Now pipelines can cross protocols:

Source	Destination
SSH	Local shell
Shell	File
HTTP	Parser
File	Database

This is beyond Unix pipes.

It’s execution graphs.


TextCommand   → local shell
SSHCommand    → remote shell
HTTPCommand   → REST call
FileCommand   → write/read
SQLCommand    → database


ExecChain is a composable command execution framework for Go that extends os/exec with persistence, security policies, and multi-protocol pipelines.

Track, audit, and chain shell, SSH, HTTP, and file commands into reproducible execution graphs — with Redis caching and database storage built in.

pipeline := cmdforge.NewPipeline()

pipeline.
    SSH(sshConfig, "journalctl", []string{"-n", "500"}).
    PipeLocal("grep", []string{"ERROR"}).
    PipeHTTPPost("https://ops.internal/logs").
    PipeFile("error_report.txt")

pipeline.Run(ctx)

Different taglines:

“Embedded command orchestration framework with remote execution agents.”

A programmable command orchestration + audit + pipeline system with multi-protocol execution (shell, SSH, HTTP, file, DB) and persistence.

If RunDeck and Ansible had a baby .. but it came out as an embeddable API for Dev teams.


This framework becomes compelling when:

Command execution is part of the product

Not just an operational concern

DevOps handles:

Deploying systems

Embedded orchestration handles:

Operating systems programmatically from within software

Different layers of the stack.

Implements the Chain of Responsibility design pattern:
https://refactoring.guru/design-patterns/chain-of-responsibility/go/example

There are no rogue commands - must be handled in the context of Execution manager,
which validates, scrubs, executes, and handles directed output and logging.
Each component hands off to the next.
Command → Scrubber → Policy Engine → Logger → Executor → Post-Processor → Store

2/15
Post-Processor (Handlers?, this could be a tie into any app or process for Data Extraction/Analysis etc.,)
Processing the data is the responsiblity of another framework or user code and doesn't belong in the framework(maybe)

Study this pattern*
*/

type CommandStatus string // Reporting the status of the command

const (
	StatusPending  CommandStatus = "PENDING"
	StatusRunning  CommandStatus = "RUNNING"
	StatusSuccess  CommandStatus = "SUCCESS"
	StatusFailed   CommandStatus = "FAILED"
	StatusRejected CommandStatus = "REJECTED (SECURITY)"
)

const (
	_ = iota
	CommandType_NIL
	CommandType_TEXT
	CommandType_WEB
	CommandType_DATA
	CommandType_OTHER
)

const (
	_ = iota
	RunnerType_Console
	RunnerType_FlatFile
	RunnerType_HTTP
	RunnerType_UDP
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
		fmt.Println(fmt.Errorf("UNKNOWN ERROR OCCURRED: %v", cmd))
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
		fmt.Println(fmt.Errorf("UNKNOWN ERROR OCCURRED: %v", cmd))
	}
}

// TODO, add a way for the user to add more Deny Commands
var DefaultDenyCommands = []string{
	"sudo",
	"rm",
	"dd",
	"mkfs",
	"shutdown",
	"reboot",
	"halt",
	"poweroff",
	"init",
	"kill",
	"killall",
	"pkill",
	"chmod",
	"chown",
	"mount",
	"umount",
	"iptables",
}

// TODO, add a way for the user to add more Deny Commands
var DefaultDenyPatterns = []string{
	"rm -rf /",
	"rm -rf /*",
	":(){ :|:& };:", // fork bomb
	"dd if=",
	"mkfs.",
	"> /dev/",
}

var DefaultProtectedPaths = []string{
	"/",
	"/boot",
	"/etc",
	"/bin",
	"/usr",
	"/lib",
	"/sys",
	//"/proc",
	"/dev",
}

// TODO 2-17
// Explore lineage/tracking (Linked List design)
/*
lineage tracking
PrevID	Backtracking, failure tracing
NextID	Forward traversal, replay
ParentID	Branch lineage
RootID	Workflow grouping
*/
type Command struct {
	ID       uuid.UUID
	Name     string
	Category int
	Args     []string
	Notes    string

	// Basic lineage
	//PrevID *string
	//NextID *string

	// Optional richer lineage (NOT AVAILABLE yet)
	//ParentID *string // spawned from
	//RootID   *string // workflow root

	Stdout   string
	Stderr   string
	ExitCode int
	Error    string

	Status    CommandStatus
	CreatedAt time.Time
	StartedAt time.Time
	EndedAt   time.Time
}

func NewCommand(name string, args []string, notes string) *Command {
	return &Command{
		ID:        uuid.New(),
		Name:      name,
		Args:      args,
		Notes:     notes,
		Status:    StatusPending,
		CreatedAt: time.Now(),
	}
}

type CommandScrubber interface {
	Scrub(cmd Command) error
}

type ScrubPolicy struct {
	DenyCommands   []string
	DenyPatterns   []string
	ProtectedPaths []string

	AllowCommands []string // optional allowlist mode
	AllowMode     bool
}

type DefaultScrubber struct {
	Policy ScrubPolicy
}

func NewDefaultScrubber() *DefaultScrubber {
	return &DefaultScrubber{
		Policy: ScrubPolicy{
			DenyCommands:   DefaultDenyCommands,
			DenyPatterns:   DefaultDenyPatterns,
			ProtectedPaths: DefaultProtectedPaths,
			AllowMode:      false,
		},
	}
}

func (s *DefaultScrubber) Scrub(
	cmd *Command,
) error {

	name := strings.ToLower(cmd.Name)

	// ---- Allowlist Mode ----
	if s.Policy.AllowMode {
		if !slices.Contains(s.Policy.AllowCommands, name) {
			return errors.New("command not in allowlist")
		}
	}

	// ---- Deny Command ----
	if slices.Contains(s.Policy.DenyCommands, name) {
		return errors.New("command denied by policy: " + name)
	}

	// ---- Argument String ----
	full := name + " " + strings.Join(cmd.Args, " ")
	full = strings.ToLower(full)

	// ---- Pattern Checks ----
	for _, pattern := range s.Policy.DenyPatterns {
		if strings.Contains(full, pattern) {
			return errors.New(
				"command contains dangerous pattern: " + pattern,
			)
		}
	}

	// ---- Protected Paths ----
	for _, arg := range cmd.Args {
		for _, path := range s.Policy.ProtectedPaths {
			if strings.HasPrefix(arg, path) {
				return errors.New(
					"operation on protected path: " + path,
				)
			}
		}
	}

	return nil
}

type CommandStore interface {
	//Store a single Command (InMemory, SQLITE, FlatFile)
	Create(ctx context.Context, cmd *Command) error
	//Retrieve a single Command Record
	GetByID(ctx context.Context, id uuid.UUID) (*Command, error)
	//Retrieve all previous Command Records (return any depends on implementation)
	GetAll(ctx context.Context) ([]*Command, error)
	//Update a single Command record, usually called internally for updating active Commands
	Update(ctx context.Context, cmd *Command) error
	//Delete - Stores don't delete. No compromised Audit trail. Every execution captured.
	//Store size can be managed separately
}

type InMemoryCommandStore struct {
	mu   sync.RWMutex
	data map[uuid.UUID]*Command
}

/* SQLITE impl */ // TODO, testing
// 2/16, flesh out, test, conform to CommandStore interface
type SQLiteCommandStore struct {
	db *sql.DB
}

func NewSqliteCommandStore(db *sql.DB) *SQLiteCommandStore {
	return &SQLiteCommandStore{db: db}
}

func (s *SQLiteCommandStore) GetAll(
	ctx context.Context,
) ([]*Command, error) {

	query := `
        SELECT
            id,
            name,
            status,
            created_at,
            started_at,
            finished_at,
            stdout,
            stderr,
            exit_code,
            metadata_json
        FROM commands
        ORDER BY created_at DESC
    `

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commands []*Command

	for rows.Next() {
		cmd := new(Command)

		err := rows.Scan(
			&cmd.ID,
			&cmd.Name,
			&cmd.Status,
			&cmd.CreatedAt,
			&cmd.StartedAt,
			&cmd.EndedAt,
			&cmd.Stdout,
			&cmd.Stderr,
			&cmd.ExitCode,
		)
		if err != nil {
			return nil, err
		}

		commands = append(commands, cmd)
	}

	return commands, nil
}

func (s *SQLiteCommandStore) GetByID(
	ctx context.Context,
	uuid uuid.UUID,
) (*Command, error) {

	query := `
        SELECT
            id,
            name,
            status,
            created_at,
            started_at,
            finished_at,
            stdout,
            stderr,
            exit_code,
            metadata_json
        FROM commands
        WHERE id = ?
        LIMIT 1
    `
	uuidString := uuid.String()
	row := s.db.QueryRowContext(ctx, query, uuidString)

	cmd := new(Command)

	err := row.Scan(
		&cmd.ID,
		&cmd.Name,
		&cmd.Status,
		&cmd.CreatedAt,
		&cmd.StartedAt,
		&cmd.EndedAt,
		&cmd.Stdout,
		&cmd.Stderr,
		&cmd.ExitCode,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf(
				"command not found: %s",
				uuidString,
			)
		}
		return nil, err
	}

	return cmd, nil
}

func (s *SQLiteCommandStore) Update(
	ctx context.Context,
	cmd *Command,
) error {

	query := `
        UPDATE commands
        SET
            name = ?,
            status = ?,
            started_at = ?,
            finished_at = ?,
            stdout = ?,
            stderr = ?,
            exit_code = ?,
            metadata_json = ?
        WHERE id = ?
    `

	result, err := s.db.ExecContext(
		ctx,
		query,
		cmd.Name,
		cmd.Status,
		cmd.StartedAt,
		cmd.EndedAt,
		cmd.Stdout,
		cmd.Stderr,
		cmd.ExitCode,
		cmd.ID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf(
			"no command updated (id=%s)",
			cmd.ID,
		)
	}

	return nil
}

// TODO
func (s *SQLiteCommandStore) SaveBatch(
	ctx context.Context,
	cmds []*Command,
) error {

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, `
        INSERT INTO commands (
            id, name, status, created_at
        )
        VALUES (?, ?, ?, ?)
    `)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, cmd := range cmds {
		_, err := stmt.ExecContext(
			ctx,
			cmd.ID,
			cmd.Name,
			cmd.Status,
			cmd.CreatedAt,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (s *SQLiteCommandStore) MarkStarted(
	ctx context.Context,
	id string,
	startedAt time.Time,
) error {

	query := `
        UPDATE commands
        SET
            status = 'RUNNING',
            started_at = ?
        WHERE id = ?
    `

	_, err := s.db.ExecContext(
		ctx,
		query,
		startedAt,
		id,
	)

	return err
}

func (s *SQLiteCommandStore) MarkFinished(
	ctx context.Context,
	id string,
	finishedAt time.Time,
	exitCode int,
	stdout string,
	stderr string,
) error {

	query := `
        UPDATE commands
        SET
            status = 'COMPLETED',
            finished_at = ?,
            exit_code = ?,
            stdout = ?,
            stderr = ?
        WHERE id = ?
    `

	_, err := s.db.ExecContext(
		ctx,
		query,
		finishedAt,
		exitCode,
		stdout,
		stderr,
		id,
	)

	return err
}

// TODO
func (s *SQLiteCommandStore) GetRecent(
	ctx context.Context,
	limit uint,
) ([]*Command, error) {

	query := `
        SELECT
            id,
            name,
            status,
            created_at,
            started_at,
            finished_at,
            stdout,
            stderr,
            exit_code,
            metadata_json
        FROM commands
        ORDER BY created_at DESC
        LIMIT ?
    `

	rows, err := s.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commands []*Command

	for rows.Next() {
		cmd := new(Command)

		err := rows.Scan(
			&cmd.ID,
			&cmd.Name,
			&cmd.Status,
			&cmd.CreatedAt,
			&cmd.StartedAt,
			&cmd.EndedAt,
			&cmd.Stdout,
			&cmd.Stderr,
			&cmd.ExitCode,
		)
		if err != nil {
			return nil, err
		}

		commands = append(commands, cmd)
	}

	return commands, nil
}

func (s *SQLiteCommandStore) Create(
	ctx context.Context,
	cmd *Command,
) error {

	query := `
        INSERT INTO commands (
            id,
            name,
            status,
            created_at,
            started_at,
            finished_at,
            stdout,
            stderr,
            exit_code,
            metadata_json
        )
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `

	_, err := s.db.ExecContext(
		ctx,
		query,
		cmd.ID,
		cmd.Name,
		cmd.Status,
		cmd.CreatedAt,
		cmd.StartedAt,
		cmd.EndedAt,
		cmd.Stdout,
		cmd.Stderr,
		cmd.ExitCode,
	)

	return err
}

//TODO - Postgres store

func NewInMemoryStore() *InMemoryCommandStore {
	return &InMemoryCommandStore{
		data: make(map[uuid.UUID]*Command),
	}
}

// Depending on how large the map gets, this will help recycle memory.
// See here --> https://medium.com/@caring_smitten_gerbil_914/go-maps-and-hidden-memory-leaks-what-every-developer-should-know-17b322b177eb
// This may not be neccessary since we are storing single value pointers as opposed to 128 byte buckets
// Though, we know now that some of the command data will grow fairly large, UTF-8 text.
func (s *InMemoryCommandStore) shrinkMap(old map[uuid.UUID]*Command) map[uuid.UUID]*Command {
	newMap := make(map[uuid.UUID]*Command, len(old))
	maps.Copy(newMap, old)
	return newMap
}

func (s *InMemoryCommandStore) Create(
	ctx context.Context,
	cmd *Command,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[cmd.ID] = cmd
	return nil
}

func (s *InMemoryCommandStore) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*Command, error) {

	s.mu.RLock()
	defer s.mu.RUnlock()

	cmd, ok := s.data[id]
	if !ok {
		return nil, errors.New("command not found")
	}

	return cmd, nil
}

// InMemoryStore puts data in Map (for various reasons), but API always works with an Array/Slice
func (s *InMemoryCommandStore) GetAll(ctx context.Context) ([]*Command, error) {

	// Honor context cancellation
	//Overkill here, add in SQL and Network Stores
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.data) == 0 {
		return []*Command{}, nil
	}

	mapToList := make([]*Command, 0, len(s.data))

	for _, cmd := range s.data {
		mapToList = append(mapToList, cmd)
	}

	return mapToList, nil
}

// internal function for InMemoryCommandStore to return the in memory map as is (not a List)
func (s *InMemoryCommandStore) memoryMap() (map[uuid.UUID]*Command, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.data, nil
}

func (s *InMemoryCommandStore) Update(
	ctx context.Context,
	cmd *Command,
) error {

	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[cmd.ID] = cmd
	return nil
}

//SQL Store Implementation (SQLITE default)
/*
CREATE TABLE IF NOT EXISTS commands (
    id	 		  PRIMARY KEY AUTOINCREMENT
	uuid          TEXT,
    name          TEXT NOT NULL,
    status        TEXT NOT NULL,
    created_at    DATETIME NOT NULL,
    started_at    DATETIME,
    finished_at   DATETIME,
    stdout        TEXT,
    stderr        TEXT,
    exit_code     INTEGER
);

CREATE INDEX idx_commands_created_at
ON commands(created_at DESC);
*/

//2/15 TODO
// Lineage tracking? Command History Linked structure?
// Graph/chained output for debugging
// Traverse and Rewind capability(?)

// Stores the results/metadata of a Command Operation
type ExecutionResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
	Error    string

	StartedAt time.Time
	EndedAt   time.Time
	Duration  time.Duration // TODO, we should handle this somewhere, important info
}

// Define the Executor contract for the Service layer. All Commands are Executed by Executors.
type CommandExecutor interface {
	Execute(
		ctx context.Context,
		cmd *Command, debug bool,
	) (*ExecutionResult, error)
}

// Local (not Remote) Command Executor (Default)
type LocalExecutor struct{}

func NewLocalExecutor() *LocalExecutor {
	return &LocalExecutor{}
}

func (le *LocalExecutor) debugDump(cmd *Command, er *ExecutionResult, logFile string) {
	// Open the log file. O_APPEND appends to an existing file, O_CREATE creates the file if it
	// doesn't exist, and O_WRONLY opens the file in write-only mode.

	ioHelper := &IoHelper{}

	file := ioHelper.GetFile(logFile)

	// Ensure the file is closed when the main function exits.
	defer file.Close()

	// Set the standard logger's output to the file.
	log.SetOutput(file)

	// Log messages will now be written to "application.log" instead of stderr.
	log.Println("===========================================================================================")
	log.Println("::BEGIN EXECUTION::")
	log.Println("Time: ", time.Now())
	log.Println("Name: ", cmd.Name)
	log.Println("Args: ", cmd.Args)
	log.Println("Notes: ", cmd.Notes)
	log.Println("Status: ", cmd.Status)
	log.Println("Command Started: ", cmd.StartedAt)
	log.Println("Execution Ended: ", er.EndedAt)
	log.Println("Duration: ", er.Duration)
	log.Println("ExitCode: ", er.ExitCode)
	log.Println("::END EXECUTION::")
	log.Println("===========================================================================================")
}

func (e *LocalExecutor) Execute(
	ctx context.Context,
	cmd *Command, debug bool,
) (*ExecutionResult, error) {

	start := time.Now()

	c := exec.CommandContext(ctx, cmd.Name, cmd.Args...)

	var stdout, stderr bytes.Buffer
	c.Stdout = &stdout
	c.Stderr = &stderr

	//Security Audit/Scub check happens in Service, only valid commands make it to the Executor
	err := c.Run()

	end := time.Now()

	result := &ExecutionResult{
		Stdout:    stdout.String(),
		Stderr:    stderr.String(),
		StartedAt: start,
		EndedAt:   end,
		Duration:  end.Sub(start),
	}

	if c.ProcessState != nil {
		result.ExitCode = c.ProcessState.ExitCode()
	}

	if err != nil {
		result.Error = err.Error()
	}

	if debug {
		go e.debugDump(cmd, result, "executions.log")
	}

	return result, err
}

// This contract defines the API/entrypoint for end-user capabilities
// The implementation can specify anything else involved with Running the service
// ie, the CommandService Executes, Stores, and Retrieves data but is all encompassed by a call to "Run"
// RunHistory is a way to view and interact with prior Service Runs
// ie, the CommandService impl will return (Store.GetAll()) all prior or specifc (Store.GetById()) Service Runs
/*
type Service interface {
	Run(ctx context.Context, a any) error
	RunHistory(ctx context.Context, limit uint) ([]*any, error)
}
*/

// A Command Service must have access to an Executor for managing Commands and a Store for persisting results.
// The service handles implementation specific details and communication.
type CommandService struct {
	//Service
	Store    CommandStore
	Executor CommandExecutor
}

func NewCommandService(
	store CommandStore,
	exec CommandExecutor,
) *CommandService {
	return &CommandService{
		Store:    store,
		Executor: exec,
	}
}

// TODO
func (s *CommandService) CommandHistory(ctx context.Context, limit uint) ([]*Command, error) {
	var ok error
	var allCommands []*Command
	var commandSubset []*Command

	if limit == 0 {
		allCommands, ok = s.Store.GetAll(ctx)
		if ok != nil {
			return nil, errors.New("ERR - RETRIEVAL")
		}
		return allCommands, nil
	}

	if limit > 0 {
		var subsetLimit uint
		if limit > 1 {
			subsetLimit = limit - 1
		}
		commandSubset = make([]*Command, limit)
		allCommands, ok = s.Store.GetAll(ctx)
		if ok != nil {
			return nil, errors.New("ERR - RETRIEVAL")
		}
		for i := range subsetLimit {
			commandSubset = append(commandSubset, allCommands[i])
		}
		return commandSubset, nil
	}

	return nil, errors.New("unexpected error has occurred")
}

func (s *CommandService) RunCommand(
	ctx context.Context,
	cmd *Command, debug bool,
) error {

	if NewDefaultScrubber().Scrub(cmd) != nil {
		violation := "SECURITY POLICY TRIGGERED"

		// Mark rejected
		cmd.Status = StatusRejected
		cmd.StartedAt = time.Now()
		cmd.Stdout = violation
		cmd.Stderr = violation
		cmd.EndedAt = time.Now()
		// Keep track of our rejections (Audit everything. Track everything.)
		// We may also keep security violations in a separate text log
		s.Store.Update(ctx, cmd)
		//PrintFailure(cmd.Name)
		//PrintFailure(string(StatusRejected))
		PrintFailure(violation)
		return errors.New(string(StatusRejected))
	}

	// Persist initial record
	if err := s.Store.Create(ctx, cmd); err != nil {
		return err
	}

	// Mark running
	cmd.Status = StatusRunning
	cmd.StartedAt = time.Now()
	s.Store.Update(ctx, cmd)

	// Execute
	result, err := s.Executor.Execute(ctx, cmd, debug)

	// Apply results
	cmd.Stdout = result.Stdout
	cmd.Stderr = result.Stderr
	cmd.ExitCode = result.ExitCode
	cmd.Error = result.Error
	cmd.EndedAt = result.EndedAt

	if err != nil {
		cmd.Status = StatusFailed
	} else {
		cmd.Status = StatusSuccess
	}

	return s.Store.Update(ctx, cmd)
}

// Types built for Automation and Testing purposes
type CommandRunner interface {
	RunCommands(svc *CommandService, ctx context.Context, cmds []*Command, debug bool) []*Command
}

type ConsoleCommandRunner struct{}
type HTTPCommandRunner struct{}
type UDPCommandRunner struct{}
type FlatFileCommandRunner struct{}

// TODO!
func NewCommandRunner(runnerType uint) CommandRunner {
	// Not yet implemented
	/*
		switch runnerType {
		case RunnerType_Console:
			return &ConsoleCommandRunner{}
		case RunnerType_FlatFile:
			return &FlatFileCommandRunner{}
		case RunnerType_HTTP:
			return &HTTPCommandRunner{}
		case RunnerType_UDP:
			return &UDPCommandRunner{}
		}
	*/
	log.Printf("<Remove in Test>::Default Runner Selected: %v", runnerType)
	return &ConsoleCommandRunner{}
}

func (ccr *ConsoleCommandRunner) RunCommands(
	svc *CommandService,
	ctx context.Context,
	cmds []*Command, debug bool,
) []*Command {

	var wg sync.WaitGroup

	finished := make([]*Command, len(cmds))

	for i, cmd := range cmds {
		wg.Add(1)

		go func(i int, cmd *Command) {
			defer wg.Done()

			if err := svc.RunCommand(ctx, cmd, debug); err != nil {
				PrintFailure("\nERR --> See Logs::<<%v>>::\n", err)
			}

			finished[i] = cmd
		}(i, cmd)
	}

	wg.Wait()
	return finished
}

//End FRAMEWORK
