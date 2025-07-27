package routes

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/model"
	"net/http"
)

func SystemPromptsController(router *gin.RouterGroup, db *sql.DB) {
	systemPromptsRouter := router.Group("/system-prompts")

	systemPromptsRouter.GET("/", func(c *gin.Context) {
		prompts, err := model.AllSystemPrompts(db)
		respondList(c, prompts, err)
	})

	systemPromptsRouter.GET("/:systemPromptId", func(c *gin.Context) {
		systemPromptId, err := getID(c, "systemPromptId")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid system prompt ID"})
			return
		}

		prompts, err := model.SystemPromptById(db, systemPromptId)
		respondSingle(c, prompts, err)
	})

	systemPromptsRouter.POST("/", func(c *gin.Context) {
		var newPrompt model.SystemPrompt
		if err := c.ShouldBind(&newPrompt); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid system prompt data"})
			return
		}

		err := model.CreateSystemPrompt(db, &newPrompt)
		respondSingle(c, &newPrompt, err)
	})

	systemPromptsRouter.PUT("/:systemPromptId", func(c *gin.Context) {
		systemPromptId, err := getID(c, "systemPromptId")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid system prompt ID"})
			return
		}
		var prompt model.SystemPrompt
		if err := c.ShouldBind(&prompt); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid system prompt data"})
			return
		}

		err = model.UpdateSystemPrompt(db, systemPromptId, &prompt)
		respondSingle(c, &prompt, err)
	})

	systemPromptsRouter.DELETE("/:systemPromptId", func(c *gin.Context) {
		systemPromptId, err := getID(c, "systemPromptId")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid system prompt ID"})
			return
		}

		err = model.DeleteSystemPrompt(db, systemPromptId)
		respondDeleted(c, err)
	})
}
