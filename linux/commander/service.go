package main

import (
	"context"
	"errors"
	"time"
)

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

// TODO (remove?)
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
		//ioHelper := &IoHelper{}
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
		PrintFailure(violation)
		//ioHelper.FileDump(cmd, "security.log")
		//get weird results when trying to write to console and log file at same time
		//probably some stdout concurrency thing I'm not aware of (2/17, TODO)
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
