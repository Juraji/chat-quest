package instructions

import (
	"net/http"

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

	instructionsRouter.GET("/default-templates", func(c *gin.Context) {
		var tplMap map[int]*Instruction
		for idx, tpl := range defaultInstructionTemplates {
			tplMap[idx] = tpl
		}

		c.JSON(http.StatusOK, tplMap)
	})

	instructionsRouter.POST("/default-templates/use/:templateIndex", func(c *gin.Context) {
		templateIndex, ok := controllers.GetParamAsID(c, "templateIndex")
		if !ok || templateIndex < 0 || templateIndex >= len(defaultInstructionTemplates) {
			controllers.RespondBadRequest(c, "Invalid template index", nil)
			return
		}

		tpl := defaultInstructionTemplates[templateIndex]
		instruction, err := reifyInstructionTemplate(tpl)
		if err != nil {
			controllers.RespondInternalError(c, err)
			return
		}

		err = CreateInstruction(instruction)
		controllers.RespondSingle(c, &instruction, err)
	})
}
