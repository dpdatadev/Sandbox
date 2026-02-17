package main

// Testing for lineage/tracker impl
import (
	"context"
	"time"
)

type LCommand struct {
	ID string

	// Execution lineage
	PrevID *string
	NextID *string

	// Optional richer lineage
	ParentID *string // spawned from
	RootID   *string // workflow root

	Status    string
	Stdout    string
	CreatedAt time.Time
}

func (s *CommandService) Chain(
	ctx context.Context,
	cmds []*LCommand,
) error {

	if len(cmds) == 0 {
		return nil
	}

	rootID := cmds[0].ID

	for i := range cmds {

		// Root assignment
		cmds[i].RootID = &rootID

		if i > 0 {
			prev := cmds[i-1].ID
			cmds[i].PrevID = &prev
		}

		if i < len(cmds)-1 {
			next := cmds[i+1].ID
			cmds[i].NextID = &next
		}

		/*
			err := s.Store.Create(ctx, cmds[i])
			if err != nil {
				return err
			}
		*/
	}

	return nil
}

/*
func (s *CommandService) WalkForward(
	ctx context.Context,
	startID string,
) ([]*Command, error) {

	var lineage []*Command
	currentID := startID

	for {
		cmd, err := s.Store.GetByID(ctx, currentID)
		if err != nil {
			return nil, err
		}

		lineage = append(lineage, cmd)

		if cmd.NextID == nil {
			break
		}

		currentID = *cmd.NextID
	}

	return lineage, nil
}
*/

/*
func (s *CommandService) WalkBackward(
    ctx context.Context,
    startID string,
) ([]*Command, error) {

    var lineage []*Command
    currentID := startID

    for {
        cmd, err := s.Store.GetByID(ctx, currentID)
        if err != nil {
            return nil, err
        }

        lineage = append(lineage, cmd)

        if cmd.PrevID == nil {
            break
        }

        currentID = *cmd.PrevID
    }

    return lineage, nil
}
*/
