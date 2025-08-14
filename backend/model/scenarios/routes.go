package scenarios

import (
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/core/util"
)

func Routes(router *gin.RouterGroup) {
	scenariosRouter := router.Group("/scenarios")

	scenariosRouter.GET("", func(c *gin.Context) {
		scenarios, err := AllScenarios()
		util.RespondList(c, scenarios, err)
	})

	scenariosRouter.GET("/:scenarioId", func(c *gin.Context) {
		scenarioId, err := util.GetIDParam(c, "scenarioId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid scenario ID")
			return
		}

		scenarios, err := ScenarioById(scenarioId)
		util.RespondSingle(c, scenarios, err)
	})

	scenariosRouter.POST("", func(c *gin.Context) {
		var newScenario Scenario
		if err := c.ShouldBind(&newScenario); err != nil {
			util.RespondBadRequest(c, "Invalid scenario data")
			return
		}

		err := CreateScenario(&newScenario)
		util.RespondSingle(c, &newScenario, err)
	})

	scenariosRouter.PUT("/:scenarioId", func(c *gin.Context) {
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

		err = UpdateScenario(scenarioId, &scenario)
		util.RespondSingle(c, &scenario, err)
	})

	scenariosRouter.DELETE("/:scenarioId", func(c *gin.Context) {
		scenarioId, err := util.GetIDParam(c, "scenarioId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid scenario ID")
			return
		}

		err = DeleteScenario(scenarioId)
		util.RespondDeleted(c, err)
	})
}
