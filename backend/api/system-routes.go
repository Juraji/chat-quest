package api

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/database"
	"juraji.nl/chat-quest/core/log"
	"juraji.nl/chat-quest/core/providers"
	"juraji.nl/chat-quest/core/system"
)

func SystemRoutes(router *gin.RouterGroup) {
	systemRouter := router.Group("/system")

	systemRouter.POST("/tokenizer/count", func(c *gin.Context) {
		body, err := c.GetRawData()
		if err != nil {
			respondBadRequest(c, "Failed to read request body", nil)
			return
		}

		text := string(body)

		tokenCount, err := providers.TokenCount(text)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get token count"})
		} else {
			c.JSON(http.StatusOK, gin.H{"count": tokenCount})
		}
	})

	systemRouter.POST("/stop-current-generation", func(c *gin.Context) {
		system.StopCurrentGeneration.EmitBG(nil)
		respondEmpty(c, nil)
	})

	systemRouter.POST("/migrations/goto/:version", func(c *gin.Context) {
		version, _ := getParamAsID(c, "version")
		log.Get().Info("Migrating to version", zap.Int("version", version))

		database.GoToVersion(database.GetDB(), uint(version))
		respondEmpty(c, nil)
	})

	systemRouter.POST("/shutdown", func(c *gin.Context) {
		c.String(http.StatusOK, "Shutting down...")
		log.Get().Info("Shutting down by API...")

		go func() {
			// Give Gin some time to process and send the response
			time.Sleep(100 * time.Millisecond)
			os.Exit(0)
		}()
	})
}
