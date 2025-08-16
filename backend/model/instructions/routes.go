package instructions

import (
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/core/util"
)

func Routes(router *gin.RouterGroup) {
	instructionPromptsRouter := router.Group("/instruction-templates")

	instructionPromptsRouter.GET("", func(c *gin.Context) {
		prompts, err := AllInstructionPrompts()
		util.RespondList(c, prompts, err)
	})

	instructionPromptsRouter.GET("/:templateId", func(c *gin.Context) {
		templateId, err := util.GetIDParam(c, "templateId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid prompt ID")
			return
		}

		prompts, err := InstructionById(templateId)
		util.RespondSingle(c, prompts, err)
	})

	instructionPromptsRouter.POST("", func(c *gin.Context) {
		var newPrompt InstructionTemplate
		if err := c.ShouldBind(&newPrompt); err != nil {
			util.RespondBadRequest(c, "Invalid prompt data")
			return
		}

		if !newPrompt.Type.IsValid() {
			util.RespondBadRequest(c, "Invalid template type")
			return
		}

		err := CreateInstructionPrompt(&newPrompt)
		util.RespondSingle(c, &newPrompt, err)
	})

	instructionPromptsRouter.PUT("/:templateId", func(c *gin.Context) {
		templateId, err := util.GetIDParam(c, "templateId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid prompt ID")
			return
		}
		var prompt InstructionTemplate
		if err := c.ShouldBind(&prompt); err != nil {
			util.RespondBadRequest(c, "Invalid prompt data")
			return
		}

		err = UpdateInstructionPrompt(templateId, &prompt)
		util.RespondSingle(c, &prompt, err)
	})

	instructionPromptsRouter.DELETE("/:templateId", func(c *gin.Context) {
		templateId, err := util.GetIDParam(c, "templateId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid prompt ID")
			return
		}

		err = DeleteInstructionPrompt(templateId)
		util.RespondDeleted(c, err)
	})
}
