package main //package commander
//https://go.dev/blog/using-go-modules
import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"maps"
	"os"
	"os/exec"
	"runtime"
	"runtime/debug"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/google/uuid"
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

Study this pattern*
*/

type CommandStatus string // Reporting the status of the command

const (
	StatusPending CommandStatus = "PENDING"
	StatusRunning CommandStatus = "RUNNING"
	StatusSuccess CommandStatus = "SUCCESS"
	StatusFailed  CommandStatus = "FAILED"
)

// TODO
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
)

type IoHelper struct{}

func (io *IoHelper) printAlloc() {
	m := &runtime.MemStats{}
	go runtime.ReadMemStats(m)
	fmt.Printf("Allocated Heap: %v MB\n", m.Alloc/1024/1024)
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
		ConsoleStdErrHandle(cmd.Stderr) //TODO
	} else if cmd.Status == "SUCCESS" {
		PrintIdentity("Command ID: %v\n", cmd.ID)
		PrintIdentity("Command Name: %s\n", cmd.Name)
		PrintIdentity("Command Args: %s\n", cmd.Args)
		PrintSuccess("Status: %v\n", cmd.Status)
		PrintStdOut("STDOUT:\n %s\n", cmd.Stdout)
		ConsoleStdOutHandle(cmd.Stdout) //TODO
	} else {
		fmt.Println(fmt.Errorf("UNKNOWN ERROR OCCURRED: %v", cmd))
	}
}

func (io *IoHelper) DebugDump(cmd *Command, er *ExecutionResult, logFile string) {
	// Open the log file. O_APPEND appends to an existing file, O_CREATE creates the file if it
	// doesn't exist, and O_WRONLY opens the file in write-only mode.

	if logFile == "" {
		log.Printf("errors.New(\"\"): %v\n", errors.New("NO LOG FILE"))
	}

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

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
	log.Println("StartedAt: ", er.StartedAt)
	log.Println("EndedAt: ", er.EndedAt)
	log.Println("Duration: ", er.Duration)
	log.Println("ExitCode: ", er.ExitCode)
	log.Println("::END EXECUTION::")
	log.Println("===========================================================================================")
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

type Command struct {
	ID       uuid.UUID
	Name     string
	Category int
	Args     []string
	Notes    string

	Stdout   string
	Stderr   string
	ExitCode int
	Error    string

	Status CommandStatus

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
	Create(ctx context.Context, cmd *Command) error
	GetByID(ctx context.Context, id uuid.UUID) (*Command, error)
	GetAll(ctx context.Context) (map[uuid.UUID]*Command, error)
	Update(ctx context.Context, cmd *Command) error
	//Delete - Stores don't delete. No compromised Audit trail. Every execution captured.
	//Store size can be managed separately
}

type InMemoryCommandStore struct {
	mu   sync.RWMutex
	data map[uuid.UUID]*Command
}

//TODO - SQLITE memory store/file store
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

// yield iter.Seq[*Command]?
func (s *InMemoryCommandStore) GetAll(ctx context.Context) (map[uuid.UUID]*Command, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(s.data) == 0 {
		return nil, errors.New("NO DATA IN STORE")
	}

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
		cmd *Command,
	) (*ExecutionResult, error)
}

// Local (not Remote) Command Executor (Default)
type LocalExecutor struct{}

func NewLocalExecutor() *LocalExecutor {
	return &LocalExecutor{}
}

func (e *LocalExecutor) Execute(
	ctx context.Context,
	cmd *Command,
) (*ExecutionResult, error) {

	ioHelper := &IoHelper{}

	start := time.Now()

	c := exec.CommandContext(ctx, cmd.Name, cmd.Args...)

	var stdout, stderr bytes.Buffer
	c.Stdout = &stdout
	c.Stderr = &stderr

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

	go ioHelper.DebugDump(cmd, result, "executions.log")

	return result, err
}

// A Command Service must have access to an Executor for managing Commands and a Store for persisting results.
type CommandService struct {
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

func (s *CommandService) Run(
	ctx context.Context,
	cmd *Command,
) error {

	if NewDefaultScrubber().Scrub(cmd) != nil {
		violation := "SECURITY POLICY VIOLATED in COMMAND"
		PrintFailure(violation)
		panic(violation)
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
	result, err := s.Executor.Execute(ctx, cmd)

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
	RunCommands(svc *CommandService, ctx context.Context, cmds []*Command) []*Command
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
	cmds []*Command,
) []*Command {

	var wg sync.WaitGroup

	finished := make([]*Command, len(cmds))

	for i, cmd := range cmds {
		wg.Add(1)

		go func(i int, cmd *Command) {
			defer wg.Done()

			if err := svc.Run(ctx, cmd); err != nil {
				panic(err) //todo remove all panics from framework code
			}

			finished[i] = cmd
		}(i, cmd)
	}

	wg.Wait()
	return finished
}

// TODO create handlers for parsing output and doing stuff (or maybe that code should live somewhere else)
func ConsoleStdOutHandle(stdOut string) {

	if stdOut == "" {
		fmt.Println("STDOUT CANNOT BE HANDLED")
	}

	log.Println("STDOUT HANDLED")
}

func ConsoleStdErrHandle(stdErr string) {

	if stdErr == "" {
		fmt.Println("STDERR CANNOT BE HANDLED")
	}

	log.Println("STDERR HANDLED")
}

//End FRAMEWORK

// Testing
func ConsoleCommandTest() {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	store := NewInMemoryStore()
	//store := NewSqliteStore()
	exec := NewLocalExecutor()
	//store := NewSSHStore()

	svc := NewCommandService(store, exec)

	//TODO pass the list of commands from outside obviously

	hostInfo := NewCommand("uname", []string{"-a"}, "Local Host Info")

	cmd := NewCommand("ifconfig", []string{""}, "Get Local NIC Config")

	cmd1 := NewCommand("ip", []string{"neighbor"}, "Get IP Neighbor Output")

	cmd2 := NewCommand("free", []string{"-g", "-h"}, "Get Active Memory Usage")

	cmd3 := NewCommand("arp", []string{"-a"}, "Get Local ARP Cache")

	commands := []*Command{hostInfo, cmd, cmd1, cmd2, cmd3}

	consoleCommandRunner := NewCommandRunner(RunnerType_Console)

	testCommands := consoleCommandRunner.RunCommands(svc, ctx, commands)

	ioHelper := &IoHelper{}

	for _, cmd := range testCommands {
		ioHelper.ConsoleDump(cmd)
	}
}

func main() {
	fmt.Printf("%s\n", debug.Stack())
	log.SetPrefix("::Testing::")
	log.SetFlags(0)
	log.Print("main()::")
	fmt.Println("Testing Commander")
	ConsoleCommandTest()
}

// TODO - phase 2:
// Get database setup (SQLITE store)
// Piping multiple commands (still local)
// Get out of dev zone and create actual package structure

// TODO - phase 3:
// Working with HTTP(S)
// Working with SSH
// Test Scrubers/Security
// Metadata/Analysis/Examples
// LAUNCH
