package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"maragu.dev/gomponents/html"
)

// TODO, web dev with Gin and Gomponents (move to separate module/app)
func main() {
	// Create a Gin router with default middleware (logger and recovery)
	r := gin.Default()

	// Define a simple GET endpoint
	r.GET("/ping", func(c *gin.Context) {
		// Return JSON response
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/api/commands/recent", func(c *gin.Context) {
		ctx := c.Request.Context() //Is this the right context to use here for the store operations?
		db, err := GetSQLITEDB(LOCAL_SQLITE_CMD_DB4)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		store := NewSqliteCommandStore(db)

		recentCommands, err := store.GetRecent(ctx, 5)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, recentCommands)
	})

	r.GET("/", func(c *gin.Context) {
		html := html.Body()
		html.Render(c.Writer)
	})

	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	r.Run()
}
