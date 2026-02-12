package main

import (
	"errors"
	"strings"
	"time"

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
*/

type CommandStatus string
type CommandCategory int

const (
	StatusPending CommandStatus = "PENDING"
	StatusRunning CommandStatus = "RUNNING"
	StatusSuccess CommandStatus = "SUCCESS"
	StatusFailed  CommandStatus = "FAILED"

	NIL   CommandCategory = -1
	TEXT  CommandCategory = 0
	WEB   CommandCategory = 1
	DATA  CommandCategory = 2
	OTHER CommandCategory = 3
)

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
	ID    uuid.UUID
	Name  string
	Args  []string
	Notes string

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

func contains(list []string, val string) bool {
	for _, v := range list {
		if v == val {
			return true
		}
	}
	return false
}

func (s *DefaultScrubber) Scrub(
	cmd *Command,
) error {

	name := strings.ToLower(cmd.Name)

	// ---- Allowlist Mode ----
	if s.Policy.AllowMode {
		if !contains(s.Policy.AllowCommands, name) {
			return errors.New("command not in allowlist")
		}
	}

	// ---- Deny Command ----
	if contains(s.Policy.DenyCommands, name) {
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
