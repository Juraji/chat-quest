package routes

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/ai"
	"juraji.nl/chat-quest/migrations"
	"net/http"
	"os"
)

func SystemController(router *gin.RouterGroup, db *sql.DB) {
	systemRouter := router.Group("/system")

	systemRouter.POST("/tokenizer/count", func(c *gin.Context) {
		body, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
			return
		}

		text := string(body)

		tokenCount, err := ai.TokenCount(text)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get token count"})
		} else {
			c.JSON(http.StatusOK, gin.H{"count": tokenCount})
		}
	})

	systemRouter.POST("/migrations/goto/:version", func(c *gin.Context) {
		version, _ := getID(c, "version")
		fmt.Printf("Migrating to version: %d", version)
		err := migrations.GoToVersion(db, uint(version))
		respondEmpty(c, err)
	})

	systemRouter.POST("/shutdown", func(c *gin.Context) {
		done := make(chan bool)

		go func() {
			c.String(200, "Shutting down...")
			done <- true
		}()
		<-done

		os.Exit(0)
	})
}
