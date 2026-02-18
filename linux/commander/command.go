package main

import (
	"time"

	"github.com/google/uuid"
)

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

	Status    string
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
