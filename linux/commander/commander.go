package main //package commander
//https://go.dev/blog/using-go-modules
import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"maps"
	"os/exec"
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


Study this pattern*
*/

type CommandStatus string // Reporting the status of the command

const (
	StatusPending CommandStatus = "PENDING"
	StatusRunning CommandStatus = "RUNNING"
	StatusSuccess CommandStatus = "SUCCESS"
	StatusFailed  CommandStatus = "FAILED"
)

const (
	NIL = iota
	TEXT
	WEB
	DATA
	OTHER
)

var (
	PrintIdentity = color.New(color.Bold, color.FgGreen, color.Italic).PrintfFunc()
	PrintSuccess  = color.New(color.Bold, color.FgGreen, color.Underline).PrintfFunc()
	PrintStdOut   = color.New(color.Bold, color.FgYellow).PrintfFunc()
	PrintStdErr   = color.New(color.Bold, color.FgHiRed).PrintfFunc()
)

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
	"/proc",
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
	Update(ctx context.Context, cmd *Command) error
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
	Duration  time.Duration
}

// Define the Executor contract for the Service layer. All Commands pass through the Executor.
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

// Testing

func CommandTestRunner(
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

func CommandSystemTest() {
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

	cmd := NewCommand("ifconfig", []string{""}, "get local ip info")

	cmd1 := NewCommand("ip", []string{"neighbor"}, "IP Test Command")

	commands := []*Command{cmd, cmd1}

	commands = append(commands, NewCommand("arp", []string{"-a"}, "Arp Test Command"))

	testCommands := CommandTestRunner(svc, ctx, commands)

	for _, cmd := range testCommands {
		PrintIdentity("Command ID: %v\n", cmd.ID)
		PrintSuccess("Status: %v\n", cmd.Status)
		PrintStdOut("STDOUT: %s\n", cmd.Stdout)

		if cmd.Stderr != "" {
			PrintStdErr("STDERR: %s::<%s>\n", cmd.Stderr, cmd.Error)
		}
	}
}

func main() {
	fmt.Printf("%s", debug.Stack())
	log.SetPrefix("::Testing::")
	log.SetFlags(0)
	log.Print("main()::")
	fmt.Println("Testing Commander")
	CommandSystemTest()
}
