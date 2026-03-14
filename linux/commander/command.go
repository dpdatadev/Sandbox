package main

import (
	"os/user"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Reporting the status of the command
const (
	StatusTracked  = "LINEAGE"
	StatusPending  = "PENDING"
	StatusRunning  = "RUNNING"
	StatusSuccess  = "SUCCESS"
	StatusFailed   = "FAILED"
	StatusRejected = "REJECTED (SECURITY)"
)

const (
	CommandType_LIN = iota //lineage object = 0
	CommandType_NIL
	CommandType_TEXT
	CommandType_WEB
	CommandType_DATA
	CommandType_OTHER
)

type Command struct {
	ID       uuid.UUID
	Name     string
	Category int
	Args     []string
	Notes    string

	// Basic lineage (handled in DAO (CommandLineage))
	//PrevID *string
	//NextID *string

	// Optional richer lineage (NOT AVAILABLE yet)
	//ParentID *string // spawned from
	//RootID   *string // workflow root

	Stdout   string
	Stderr   string
	ExitCode int
	Error    string

	Status    string
	CreatedAt time.Time
	StartedAt time.Time
	EndedAt   time.Time
}

func (c *Command) ExecString() string {
	return c.Name + " " + strings.Join(c.Args, " ")
}

func (c *Command) GetUserName() string {
	current_user, err := user.Current()
	if err != nil {
		PrintStdErr("USER LOOKUP err OCCURRED: ", err)
	}

	username := current_user.Username

	return username
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
