package routes

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/model"
	"net/http"
)

func InstructionPromptsController(router *gin.RouterGroup, db *sql.DB) {
	instructionPromptsRouter := router.Group("/instruction-prompts")

	instructionPromptsRouter.GET("", func(c *gin.Context) {
		prompts, err := model.AllInstructionPrompts(db)
		respondList(c, prompts, err)
	})

	instructionPromptsRouter.GET("/:promptId", func(c *gin.Context) {
		promptId, err := getIDParam(c, "promptId")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid prompt ID"})
			return
		}

		prompts, err := model.InstructionPromptById(db, promptId)
		respondSingle(c, prompts, err)
	})

	instructionPromptsRouter.POST("", func(c *gin.Context) {
		var newPrompt model.InstructionPrompt
		if err := c.ShouldBind(&newPrompt); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid prompt data"})
			return
		}

		err := model.CreateInstructionPrompt(db, &newPrompt)
		respondSingle(c, &newPrompt, err)
	})

	instructionPromptsRouter.PUT("/:promptId", func(c *gin.Context) {
		promptId, err := getIDParam(c, "promptId")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid prompt ID"})
			return
		}
		var prompt model.InstructionPrompt
		if err := c.ShouldBind(&prompt); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid prompt data"})
			return
		}

		err = model.UpdateInstructionPrompt(db, promptId, &prompt)
		respondSingle(c, &prompt, err)
	})

	instructionPromptsRouter.DELETE("/:promptId", func(c *gin.Context) {
		promptId, err := getIDParam(c, "promptId")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid prompt ID"})
			return
		}

		err = model.DeleteInstructionPrompt(db, promptId)
		respondDeleted(c, err)
	})
}
