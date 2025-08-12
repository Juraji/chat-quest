package scenarios

import (
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/cq"
	"juraji.nl/chat-quest/util"
)

func Routes(cq *cq.ChatQuestContext, router *gin.RouterGroup) {
	scenariosRouter := router.Group("/scenarios")

	scenariosRouter.GET("", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		scenarios, err := AllScenarios(cq)
		util.RespondList(c, scenarios, err)
	})

	scenariosRouter.GET("/:scenarioId", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		scenarioId, err := util.GetIDParam(c, "scenarioId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid scenario ID")
			return
		}

		scenarios, err := ScenarioById(cq, scenarioId)
		util.RespondSingle(c, scenarios, err)
	})

	scenariosRouter.POST("", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		var newScenario Scenario
		if err := c.ShouldBind(&newScenario); err != nil {
			util.RespondBadRequest(c, "Invalid scenario data")
			return
		}

		err := CreateScenario(cq, &newScenario)
		util.RespondSingle(c, &newScenario, err)
	})

	scenariosRouter.PUT("/:scenarioId", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		scenarioId, err := util.GetIDParam(c, "scenarioId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid scenario ID")
			return
		}

		var scenario Scenario
		if err = c.ShouldBind(&scenario); err != nil {
			util.RespondBadRequest(c, "Invalid scenario data")
			return
		}

		err = UpdateScenario(cq, scenarioId, &scenario)
		util.RespondSingle(c, &scenario, err)
	})

	scenariosRouter.DELETE("/:scenarioId", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		scenarioId, err := util.GetIDParam(c, "scenarioId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid scenario ID")
			return
		}

		err = DeleteScenario(cq, scenarioId)
		util.RespondDeleted(c, err)
	})
}
