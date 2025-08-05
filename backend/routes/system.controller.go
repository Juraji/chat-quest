package routes

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/ai"
	"juraji.nl/chat-quest/database"
	"log"
	"net/http"
	"os"
	"time"
)

func SystemController(router *gin.RouterGroup, db *sql.DB) {
	systemRouter := router.Group("/system")

	systemRouter.POST("/tokenizer/count", func(c *gin.Context) {
		body, err := c.GetRawData()
		if err != nil {
			respondBadRequest(c, "Failed to read request body")
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
		version, _ := getIDParam(c, "version")
		log.Printf("Migrating to version: %d", version)
		err := database.GoToVersion(db, uint(version))
		respondEmpty(c, err)
	})

	systemRouter.POST("/shutdown", func(c *gin.Context) {
		c.String(http.StatusOK, "Shutting down...")
		log.Print("Shutdown requested from API, goodbye!")

		go func() {
			// Give Gin some time to process and send the response
			time.Sleep(100 * time.Millisecond)
			os.Exit(0)
		}()
	})
}
