package routes

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/model"
	"juraji.nl/chat-quest/util"
	"net/http"
)

func InstructionPromptsController(router *gin.RouterGroup, db *sql.DB) {
	instructionPromptsRouter := router.Group("/instruction-templates")

	instructionPromptsRouter.GET("", func(c *gin.Context) {
		prompts, err := model.AllInstructionPrompts(db)
		util.RespondList(c, prompts, err)
	})

	instructionPromptsRouter.GET("/:templateId", func(c *gin.Context) {
		templateId, err := util.GetIDParam(c, "templateId")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid prompt ID"})
			return
		}

		prompts, err := model.InstructionPromptById(db, templateId)
		util.RespondSingle(c, prompts, err)
	})

	instructionPromptsRouter.POST("", func(c *gin.Context) {
		var newPrompt model.InstructionTemplate
		if err := c.ShouldBind(&newPrompt); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid prompt data"})
			return
		}

		err := model.CreateInstructionPrompt(db, &newPrompt)
		util.RespondSingle(c, &newPrompt, err)
	})

	instructionPromptsRouter.PUT("/:templateId", func(c *gin.Context) {
		templateId, err := util.GetIDParam(c, "templateId")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid prompt ID"})
			return
		}
		var prompt model.InstructionTemplate
		if err := c.ShouldBind(&prompt); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid prompt data"})
			return
		}

		err = model.UpdateInstructionPrompt(db, templateId, &prompt)
		util.RespondSingle(c, &prompt, err)
	})

	instructionPromptsRouter.DELETE("/:templateId", func(c *gin.Context) {
		templateId, err := util.GetIDParam(c, "templateId")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid prompt ID"})
			return
		}

		err = model.DeleteInstructionPrompt(db, templateId)
		util.RespondDeleted(c, err)
	})
}
