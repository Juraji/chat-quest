package scenarios

import (
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/core"
	"juraji.nl/chat-quest/core/util"
)

func Routes(cq *core.ChatQuestContext, router *gin.RouterGroup) {
	scenariosRouter := router.Group("/scenarios")

	scenariosRouter.GET("", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		scenarios, err := AllScenarios(cq)
		util.RespondList(rcq, c, scenarios, err)
	})

	scenariosRouter.GET("/:scenarioId", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		scenarioId, err := util.GetIDParam(c, "scenarioId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid scenario ID")
			return
		}

		scenarios, err := ScenarioById(rcq, scenarioId)
		util.RespondSingle(rcq, c, scenarios, err)
	})

	scenariosRouter.POST("", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		var newScenario Scenario
		if err := c.ShouldBind(&newScenario); err != nil {
			util.RespondBadRequest(rcq, c, "Invalid scenario data")
			return
		}

		err := CreateScenario(rcq, &newScenario)
		util.RespondSingle(rcq, c, &newScenario, err)
	})

	scenariosRouter.PUT("/:scenarioId", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		scenarioId, err := util.GetIDParam(c, "scenarioId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid scenario ID")
			return
		}

		var scenario Scenario
		if err = c.ShouldBind(&scenario); err != nil {
			util.RespondBadRequest(rcq, c, "Invalid scenario data")
			return
		}

		err = UpdateScenario(rcq, scenarioId, &scenario)
		util.RespondSingle(rcq, c, &scenario, err)
	})

	scenariosRouter.DELETE("/:scenarioId", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		scenarioId, err := util.GetIDParam(c, "scenarioId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid scenario ID")
			return
		}

		err = DeleteScenario(rcq, scenarioId)
		util.RespondDeleted(rcq, c, err)
	})
}
