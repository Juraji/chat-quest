package instructions

import (
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/core/util/controllers"
)

func Routes(router *gin.RouterGroup) {
	instructionPromptsRouter := router.Group("/instruction-templates")

	instructionPromptsRouter.GET("", func(c *gin.Context) {
		prompts, ok := AllInstructions()
		controllers.RespondList(c, ok, prompts)
	})

	instructionPromptsRouter.GET("/:templateId", func(c *gin.Context) {
		templateId, ok := controllers.GetParamAsID(c, "templateId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid prompt ID")
			return
		}

		prompts, ok := InstructionById(templateId)
		controllers.RespondSingle(c, ok, prompts)
	})

	instructionPromptsRouter.POST("", func(c *gin.Context) {
		var newPrompt InstructionTemplate
		if err := c.ShouldBind(&newPrompt); err != nil {
			controllers.RespondBadRequest(c, "Invalid prompt data")
			return
		}

		if !newPrompt.Type.IsValid() {
			controllers.RespondBadRequest(c, "Invalid template type")
			return
		}

		ok := CreateInstruction(&newPrompt)
		controllers.RespondSingle(c, ok, &newPrompt)
	})

	instructionPromptsRouter.PUT("/:templateId", func(c *gin.Context) {
		templateId, ok := controllers.GetParamAsID(c, "templateId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid prompt ID")
			return
		}
		var prompt InstructionTemplate
		if err := c.ShouldBind(&prompt); err != nil {
			controllers.RespondBadRequest(c, "Invalid prompt data")
			return
		}

		ok = UpdateInstruction(templateId, &prompt)
		controllers.RespondSingle(c, ok, &prompt)
	})

	instructionPromptsRouter.DELETE("/:templateId", func(c *gin.Context) {
		templateId, ok := controllers.GetParamAsID(c, "templateId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid prompt ID")
			return
		}

		ok = DeleteInstruction(templateId)
		controllers.RespondEmpty(c, ok)
	})
}
