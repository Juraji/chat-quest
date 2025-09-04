package instructions

import (
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/core/util/controllers"
)

func Routes(router *gin.RouterGroup) {
	instructionsRouter := router.Group("/instruction")

	instructionsRouter.GET("", func(c *gin.Context) {
		prompts, err := AllInstructions()
		controllers.RespondList(c, prompts, err)
	})

	instructionsRouter.GET("/:templateId", func(c *gin.Context) {
		templateId, ok := controllers.GetParamAsID(c, "templateId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid prompt ID", nil)
			return
		}

		prompts, err := InstructionById(templateId)
		controllers.RespondSingle(c, prompts, err)
	})

	instructionsRouter.POST("", func(c *gin.Context) {
		var newPrompt Instruction
		if err := c.ShouldBind(&newPrompt); err != nil {
			controllers.RespondBadRequest(c, "Invalid prompt data", nil)
			return
		}

		if !newPrompt.Type.IsValid() {
			controllers.RespondBadRequest(c, "Invalid template type", nil)
			return
		}

		err := CreateInstruction(&newPrompt)
		controllers.RespondSingle(c, &newPrompt, err)
	})

	instructionsRouter.PUT("/:templateId", func(c *gin.Context) {
		templateId, ok := controllers.GetParamAsID(c, "templateId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid prompt ID", nil)
			return
		}
		var prompt Instruction
		if err := c.ShouldBind(&prompt); err != nil {
			controllers.RespondBadRequest(c, "Invalid prompt data", nil)
			return
		}

		err := UpdateInstruction(templateId, &prompt)
		controllers.RespondSingle(c, &prompt, err)
	})

	instructionsRouter.DELETE("/:templateId", func(c *gin.Context) {
		templateId, ok := controllers.GetParamAsID(c, "templateId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid prompt ID", nil)
			return
		}

		err := DeleteInstruction(templateId)
		controllers.RespondEmpty(c, err)
	})
}
