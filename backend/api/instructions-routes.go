package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/model/instructions"
)

func InstructionsRoutes(router *gin.RouterGroup) {
	instructionsRouter := router.Group("/instruction")

	instructionsRouter.GET("", func(c *gin.Context) {
		prompts, err := instructions.AllInstructions()
		respondList(c, prompts, err)
	})

	instructionsRouter.GET("/:templateId", func(c *gin.Context) {
		templateId, ok := getParamAsID(c, "templateId")
		if !ok {
			respondBadRequest(c, "Invalid prompt ID", nil)
			return
		}

		prompts, err := instructions.InstructionById(templateId)
		respondSingle(c, prompts, err)
	})

	instructionsRouter.POST("", func(c *gin.Context) {
		var newPrompt instructions.Instruction
		if err := c.ShouldBind(&newPrompt); err != nil {
			respondBadRequest(c, "Invalid prompt data", nil)
			return
		}

		if !newPrompt.Type.IsValid() {
			respondBadRequest(c, "Invalid template type", nil)
			return
		}

		err := instructions.CreateInstruction(&newPrompt)
		respondSingle(c, &newPrompt, err)
	})

	instructionsRouter.PUT("/:templateId", func(c *gin.Context) {
		templateId, ok := getParamAsID(c, "templateId")
		if !ok {
			respondBadRequest(c, "Invalid prompt ID", nil)
			return
		}
		var prompt instructions.Instruction
		if err := c.ShouldBind(&prompt); err != nil {
			respondBadRequest(c, "Invalid prompt data", nil)
			return
		}

		err := instructions.UpdateInstruction(templateId, &prompt)
		respondSingle(c, &prompt, err)
	})

	instructionsRouter.DELETE("/:templateId", func(c *gin.Context) {
		templateId, ok := getParamAsID(c, "templateId")
		if !ok {
			respondBadRequest(c, "Invalid prompt ID", nil)
			return
		}

		err := instructions.DeleteInstruction(templateId)
		respondEmpty(c, err)
	})

	instructionsRouter.GET("/default-templates", func(c *gin.Context) {
		c.JSON(http.StatusOK, instructions.DefaultTemplates())
	})

	instructionsRouter.GET("/default-templates/:templateKey", func(c *gin.Context) {
		templateKey := c.Param("templateKey")
		_, exists := instructions.DefaultTemplates()[templateKey]
		if !exists {
			respondBadRequest(c, "Invalid template index", nil)
			return
		}

		instruction, err := instructions.ReifyInstructionTemplate(templateKey)
		if err != nil {
			respondInternalError(c, err)
			return
		}

		respondSingle(c, &instruction, err)
	})
}
