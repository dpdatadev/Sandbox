package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CommandDTO struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Args      string `json:"args"`
	Notes     string `json:"notes"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

func (cdto *CommandDTO) New(cmd *Command) (*CommandDTO, error) {

	if cmd == nil {
		return nil, errors.New("invalid command: nil value provided")
	}

	return &CommandDTO{
		ID:        cmd.ID.String(),
		Name:      cmd.ExecString(),
		Notes:     cmd.Notes,
		Status:    cmd.Status,
		CreatedAt: cmd.CreatedAt.String(),
	}, nil
}

// TODO, web dev with Gin and Gomponents (move to separate module/app)
func main() {
	// Create a Gin router with default middleware (logger and recovery)
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	// Define a simple GET endpoint
	r.GET("/pingtest", func(c *gin.Context) {
		// Return JSON response
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/apitest/commands/recent", func(c *gin.Context) {
		ctx := c.Request.Context() //Is this the right context to use here for the store operations?
		db, err := GetSQLITEDB(LOCAL_SQLITE_CMD_DB4)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		store := NewSqliteCommandStore(db)

		recentCommands, err := store.GetRecent(ctx, 5)
		sendCommands := make([]*CommandDTO, 0, len(recentCommands))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		for _, cmd := range recentCommands {
			PrintDebug(fmt.Sprintf("Command: %s, Status: %s\n", cmd.Name, cmd.Status))
			cmdDTO, err := (&CommandDTO{}).New(cmd)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			sendCommands = append(sendCommands, cmdDTO)
		}

		c.JSON(http.StatusOK, sendCommands)
	})

	r.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		ctx := c.Request.Context() //Is this the right context to use here for the store operations?
		db, err := GetSQLITEDB(LOCAL_SQLITE_CMD_DB4)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
			return
		}

		store := NewSqliteCommandStore(db)

		htmlCommands, err := store.GetRecent(ctx, 5)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
			return
		}

		for _, cmd := range htmlCommands {
			PrintDebug(fmt.Sprintf("Command: %s, Status: %s\n", cmd.Name, cmd.Status))
			//cmdDTO, err := (&CommandDTO{}).New(cmd)
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"Title":    "Commander Home",
			"Message":  "Welcome to the Commander App!",
			"Commands": htmlCommands,
		})
	})

	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	r.Run()
}
