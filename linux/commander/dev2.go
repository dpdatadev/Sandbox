package main

type CommandBatch struct {
	Commands   []*Command
	BatchLabel string
}

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

/*
	IN memory version for testing/debugging with container/list (Doubly Linked List)
	https://chatgpt.com/c/698c0190-d8ec-832d-8aee-537b6c64320d
	https://pkg.go.dev/container/list
	for e := chain.Back(); e != nil; e = e.Prev() {
    cmd := e.Value.(*Command)
    fmt.Println(cmd.ID, cmd.Status)
}

*/

/*
func _hydrate(commands []*Command) *CommandBatch {
	//cb := &CommandBatch{_commands: commands, _batchLabel: fmt.Sprintf("_batch_out_%s", uuid.NewString)}
	/*for _, cmd := range commands {

	}

}


func Lineage(cHistory CommandBatch) *list.List {
	//Batch save execution/history tree
	//l := list.List{}

}
*/

/*

Reverse Traversal
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


Walk Forward
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


Chain Function - setup the Lineage
Using the builtin- pointer database route
  	prev_id TEXT, (not optional)
    next_id TEXT, (not optional)
    parent_id TEXT,
    root_id TEXT,
func (s *CommandService) Chain(
    ctx context.Context,
    cmds []*Command,
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

        err := s.Store.Save(ctx, cmds[i])
        if err != nil {
            return err
        }
    }

    return nil
}



*/
