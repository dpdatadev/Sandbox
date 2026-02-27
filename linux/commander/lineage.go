package main

//BETA FEATURE

// Lineage/tracker impl
import (
	"context"
	"errors"
	"fmt"
	"strconv"
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

// Several ways to implement Chain tracking/lineage - DB persistence likely desired
func (s *SQLiteCommandStore) SaveBatchHistory(
	ctx context.Context,
	cmds []*Command,
) error {

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, `
        INSERT INTO commands (
            id, name, status, created_at
        )
        VALUES (?, ?, ?, ?)
    `)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, cmd := range cmds {
		_, err := stmt.ExecContext(
			ctx,
			cmd.ID,
			cmd.Name,
			cmd.Status,
			cmd.CreatedAt,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// //////////////////////////////////
type SingleList[T any] interface {
	Values() []T
	Len() int
	Head() (*SList[T], error)
	Tail() (*SList[T], error)
	Append(value *T) (int, error)
	ForwardNode() *SList[T]
	ForwardValue() *T
	Print()
}

// Simple Single List
// When Double list is needed - use Golang container/list
type SList[T any] struct {
	Value T
	Next  *SList[T]
}

func (sl *SList[T]) Append(value *T) (int, error) {
	if sl == nil {
		return 0, fmt.Errorf("list is nil")
	}

	current := sl
	index := 0

	for current.Next != nil {
		current = current.Next
		index++
	}

	current.Next = &SList[T]{Value: *value}
	return index + 1, nil
}

func (sl *SList[T]) Head() (*SList[T], error) {

	if sl != nil {
		return sl, nil
	}

	return nil, errors.New("list is empty")
}

func (sl *SList[T]) Tail() (*SList[T], error) {

	if sl == nil {
		return nil, errors.New("list is empty")
	}

	current := sl
	for current.Next != nil {
		current = current.Next
	}

	return current, nil
}

// Probably redundant
func (sl *SList[T]) ForwardValue() *T {
	if sl == nil || sl.Next == nil {
		return nil
	}
	return &sl.Next.Value
}

func (sl *SList[T]) ForwardNode() *SList[T] {
	if sl == nil {
		return nil
	}
	return sl.Next
}

func (sl *SList[T]) Values() []T {
	var values []T
	for current := sl; current != nil; current = current.Next {
		values = append(values, current.Value)
	}
	return values
}

// O(n) complexity - TODO store length for O(1)
func (sl *SList[T]) Len() int {
	count := 0
	for current := sl; current != nil; current = current.Next {
		count++
	}
	return count
}

func (sl *SList[T]) Print() {

	if sl == nil {
		panic("init error")
	}

	counter := 1
	for sl != nil {
		counter++
		PrintSuccess("Index: %d :: Value: %v\n", counter, sl.Value)
		sl = sl.Next
	}
}

type Lineage interface {
	BeginChain() []*CommandLineage          //Step 1 - (Hydrate) - create CommandLineage objects from Command objects (copying relevant fields and adding lineage metadata)
	LinkChain(cmds []*CommandLineage) error //Step 2 - (Chain together) - assign PrevID, NextID, RootID to CommandLineage objects to create a linked tracking chain
	LogLineage(lineage []*CommandLineage, lineageFileName string) error
	//Persist(ctx context.Context, cmds []*CommandLineage) error //Step 3 - save lineage tracking objects to database (or other store) for queryable lineage history
	//WalkForward(ctx context.Context, startID string) ([]*Command, error)
	//WalkBackward(ctx context.Context, startID string) ([]*Command, error)
}

type DBHistoryService struct {
	AuditCommands []*Command
	Store         CommandStore
}

// TODO, beta thoughts (think on this) -- I think the second struct needs to be removed and keep PrevID and NextID on Command
// Then any Command can easily be checked for lineage then move forward or backward instead of
// checking a different table/output.
type CommandLineage struct {
	ID      string
	BatchID string

	// Execution lineage
	PrevID string //* want to see the actual string value stored, not the address/reference of the previous object in memory (which is what a pointer would give us)
	NextID string //*

	// Optional richer lineage
	//ParentID string  // spawned from (copied from Command object in Lineage creation via HydrateLineage())
	RootID string //* workflow root (copied from first CommandLineage in ChainLineage())

	Status    string
	Stdout    string
	CreatedAt time.Time
}

// ///////////////////////////////////////////////////////////
// TODO - improve https://chatgpt.com/c/698c0190-d8ec-832d-8aee-537b6c64320d
func (hs *DBHistoryService) BeginChain() []*CommandLineage {

	if len(hs.AuditCommands) == 0 {
		return []*CommandLineage{}
	}

	lineageObjects := make([]*CommandLineage, 0, len(hs.AuditCommands))

	shortUUID, err := (&CmdIOHelper{}).NewShortUUID()

	var batchSuffix string

	if err != nil {
		PrintStdErr("UUID function fail: %v", err)
		batchSuffix = strconv.FormatInt(time.Now().UnixNano(), 10)
	} else {
		batchSuffix = shortUUID
	}

	batchID := fmt.Sprintf("batch__%s", batchSuffix)
	now := time.Now()

	for _, cmd := range hs.AuditCommands {

		lineageObject := &CommandLineage{
			ID:        cmd.ID.String(),
			BatchID:   batchID,
			Status:    cmd.Status,
			Stdout:    cmd.Stdout,
			CreatedAt: now, // or cmd.CreatedAt
		}

		lineageObjects = append(lineageObjects, lineageObject)
	}

	return lineageObjects
}

func (hs *DBHistoryService) LinkChain(
	cmds []*CommandLineage, //todo add history struct to keep separate table of tracking and we can join on uuid
) error {

	if len(cmds) == 0 {
		return errors.New("Chain Empty! No Commands to Link!")
	}

	rootID := cmds[0].ID

	for i := range cmds {
		// Root assignment
		cmds[i].RootID = rootID //&

		if i > 0 {
			//copy of UUID value (string)
			prev := cmds[i-1].ID
			cmds[i].PrevID = prev //&
		}

		if i < len(cmds)-1 {
			next := cmds[i+1].ID
			cmds[i].NextID = next //&
		}
	}

	return nil
}

// Write lineage graph to file
func (hs *DBHistoryService) LogLineage(lineage []*CommandLineage, lineageFileName string) error {

	f := (&CmdIOHelper{}).GetFileWrite(lineageFileName)
	if f == nil {
		err := errors.New("LINEAGE FILE ERROR")
		PrintFailure("errors.New(\"\"): %v\n", err)
		return err
	}
	defer f.Close()

	for _, cmd := range lineage {
		line := fmt.Sprintf("ID: %s, BatchID: %s, PrevID: %v, NextID: %v, Status: %s, RootID: %v\n",
			cmd.ID, cmd.BatchID, cmd.PrevID, cmd.NextID, cmd.Status, cmd.RootID)
		_, err := f.WriteString(line)
		if err != nil {
			return err
		}
	}
	return nil
}

//TODO - implement DB persistence for lineage tracking objects (could be a separate table with foreign key to Commands or a JSON blob in Commands table)
//BETA
//Database lineage will come
/*
func (dbs *DBHistoryService) WalkForward(
	ctx context.Context,
	startID string,
) ([]*Command, error) {

	var lineage []*Command
	currentID := startID

	for {
		cmd, err := dbs.Store.GetByID(ctx, currentID)
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


func (dbs *DBHistoryService) WalkBackward(
    ctx context.Context,
    startID string,
) ([]*Command, error) {

    var lineage []*Command
    currentID := startID

    for {
        cmd, err := dbs.Store.GetByID(ctx, currentID)
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
