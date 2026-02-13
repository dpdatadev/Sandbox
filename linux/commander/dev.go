package commander

//https://chatgpt.com/c/698c0190-d8ec-832d-8aee-537b6c64320d
//https://chatgpt.com/c/698c0190-d8ec-832d-8aee-537b6c64320d

import (
	"time"

	"github.com/google/uuid"
)

type CommandStatus string

const (
	StatusPending CommandStatus = "PENDING"
	StatusRunning CommandStatus = "RUNNING"
	StatusSuccess CommandStatus = "SUCCESS"
	StatusFailed  CommandStatus = "FAILED"
)

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
