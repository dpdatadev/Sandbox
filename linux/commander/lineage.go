package main

//BETA FEATURE

// Lineage/tracker impl
import (
	"context"
	"time"
)

//container List implementation for testing/debugging and InMemory stuff
/*
	https://chatgpt.com/c/698c0190-d8ec-832d-8aee-537b6c64320d
	https://pkg.go.dev/container/list
	for e := chain.Back(); e != nil; e = e.Prev() {
    cmd := e.Value.(*Command)
    fmt.Println(cmd.ID, cmd.Status)
*/

//Batch processing fits into this?
/*
type CommandBatch struct {
	Commands   []*Command
	BatchLabel string
}
*/

//Deadline - version 1 BETA

/*Lineage - History Tracking*/

/*
HTTPCommand → fetch JSON
↓
TransformCommand → jq parse
↓
DBCommand → insert rows
↓
FileCommand → archive CSV
*/

type LineageCommand struct {
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

/////////////////////////////////////////////////////////////

//Direct pointer method - database fields in Command struct or separate struct/table
/*
func (s *CommandService) ChainCommands(
	ctx context.Context,
	cmds []*Command,
) ([]*LineageCommand, error) {
	if len(cmds) == 0 {
		return nil, errors.New("NO COMMANDS GIVEN")
	}

	commandChain := make([]*LineageCommand, 0, len(cmds))

	rootID := cmds[0].ID


}
*/

func (s *CommandService) Chain(
	ctx context.Context,
	cmds []*LineageCommand, //todo add history struct to keep separate table of tracking and we can join on uuid
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
