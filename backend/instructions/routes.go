package instructions

import (
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/cq"
	"juraji.nl/chat-quest/util"
)

func Routes(cq *cq.ChatQuestContext, router *gin.RouterGroup) {
	instructionPromptsRouter := router.Group("/instruction-templates")

	instructionPromptsRouter.GET("", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		prompts, err := AllInstructionPrompts(cq)
		util.RespondList(cq, c, prompts, err)
	})

	instructionPromptsRouter.GET("/:templateId", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		templateId, err := util.GetIDParam(c, "templateId")
		if err != nil {
			util.RespondBadRequest(cq, c, "Invalid prompt ID")
			return
		}

		prompts, err := InstructionPromptById(cq, templateId)
		util.RespondSingle(cq, c, prompts, err)
	})

	instructionPromptsRouter.POST("", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		var newPrompt InstructionTemplate
		if err := c.ShouldBind(&newPrompt); err != nil {
			util.RespondBadRequest(cq, c, "Invalid prompt data")
			return
		}

		err := CreateInstructionPrompt(cq, &newPrompt)
		util.RespondSingle(cq, c, &newPrompt, err)
	})

	instructionPromptsRouter.PUT("/:templateId", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		templateId, err := util.GetIDParam(c, "templateId")
		if err != nil {
			util.RespondBadRequest(cq, c, "Invalid prompt ID")
			return
		}
		var prompt InstructionTemplate
		if err := c.ShouldBind(&prompt); err != nil {
			util.RespondBadRequest(cq, c, "Invalid prompt data")
			return
		}

		err = UpdateInstructionPrompt(cq, templateId, &prompt)
		util.RespondSingle(cq, c, &prompt, err)
	})

	instructionPromptsRouter.DELETE("/:templateId", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		templateId, err := util.GetIDParam(c, "templateId")
		if err != nil {
			util.RespondBadRequest(cq, c, "Invalid prompt ID")
			return
		}

		err = DeleteInstructionPrompt(cq, templateId)
		util.RespondDeleted(cq, c, err)
	})
}
