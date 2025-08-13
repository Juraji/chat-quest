package instructions

import (
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/core"
	"juraji.nl/chat-quest/core/util"
)

func Routes(cq *core.ChatQuestContext, router *gin.RouterGroup) {
	instructionPromptsRouter := router.Group("/instruction-templates")

	instructionPromptsRouter.GET("", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		prompts, err := AllInstructionPrompts(cq)
		util.RespondList(rcq, c, prompts, err)
	})

	instructionPromptsRouter.GET("/:templateId", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		templateId, err := util.GetIDParam(c, "templateId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid prompt ID")
			return
		}

		prompts, err := InstructionPromptById(rcq, templateId)
		util.RespondSingle(rcq, c, prompts, err)
	})

	instructionPromptsRouter.POST("", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		var newPrompt InstructionTemplate
		if err := c.ShouldBind(&newPrompt); err != nil {
			util.RespondBadRequest(rcq, c, "Invalid prompt data")
			return
		}

		err := CreateInstructionPrompt(rcq, &newPrompt)
		util.RespondSingle(rcq, c, &newPrompt, err)
	})

	instructionPromptsRouter.PUT("/:templateId", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		templateId, err := util.GetIDParam(c, "templateId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid prompt ID")
			return
		}
		var prompt InstructionTemplate
		if err := c.ShouldBind(&prompt); err != nil {
			util.RespondBadRequest(rcq, c, "Invalid prompt data")
			return
		}

		err = UpdateInstructionPrompt(rcq, templateId, &prompt)
		util.RespondSingle(rcq, c, &prompt, err)
	})

	instructionPromptsRouter.DELETE("/:templateId", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		templateId, err := util.GetIDParam(c, "templateId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid prompt ID")
			return
		}

		err = DeleteInstructionPrompt(rcq, templateId)
		util.RespondDeleted(rcq, c, err)
	})
}
