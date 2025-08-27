package instructions

import (
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/core/util/controllers"
)

func Routes(router *gin.RouterGroup) {
	instructionPromptsRouter := router.Group("/instruction-templates")

	instructionPromptsRouter.GET("", func(c *gin.Context) {
		prompts, err := AllInstructions()
		controllers.RespondListE(c, prompts, err)
	})

	instructionPromptsRouter.GET("/:templateId", func(c *gin.Context) {
		templateId, ok := controllers.GetParamAsID(c, "templateId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid prompt ID", nil)
			return
		}

		prompts, err := InstructionById(templateId)
		controllers.RespondSingleE(c, prompts, err)
	})

	instructionPromptsRouter.POST("", func(c *gin.Context) {
		var newPrompt InstructionTemplate
		if err := c.ShouldBind(&newPrompt); err != nil {
			controllers.RespondBadRequest(c, "Invalid prompt data", nil)
			return
		}

		if !newPrompt.Type.IsValid() {
			controllers.RespondBadRequest(c, "Invalid template type", nil)
			return
		}

		err := CreateInstruction(&newPrompt)
		controllers.RespondSingleE(c, &newPrompt, err)
	})

	instructionPromptsRouter.PUT("/:templateId", func(c *gin.Context) {
		templateId, ok := controllers.GetParamAsID(c, "templateId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid prompt ID", nil)
			return
		}
		var prompt InstructionTemplate
		if err := c.ShouldBind(&prompt); err != nil {
			controllers.RespondBadRequest(c, "Invalid prompt data", nil)
			return
		}

		err := UpdateInstruction(templateId, &prompt)
		controllers.RespondSingleE(c, &prompt, err)
	})

	instructionPromptsRouter.DELETE("/:templateId", func(c *gin.Context) {
		templateId, ok := controllers.GetParamAsID(c, "templateId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid prompt ID", nil)
			return
		}

		err := DeleteInstruction(templateId)
		controllers.RespondEmptyE(c, err)
	})
}
