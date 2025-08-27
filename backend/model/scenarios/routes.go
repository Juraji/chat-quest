package scenarios

import (
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/core/util/controllers"
)

func Routes(router *gin.RouterGroup) {
	scenariosRouter := router.Group("/scenarios")

	scenariosRouter.GET("", func(c *gin.Context) {
		scenarios, ok := AllScenarios()
		controllers.RespondList(c, ok, scenarios)
	})

	scenariosRouter.GET("/:scenarioId", func(c *gin.Context) {
		scenarioId, ok := controllers.GetParamAsID(c, "scenarioId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid scenario ID", nil)
			return
		}

		scenario, ok := ScenarioById(scenarioId)
		controllers.RespondSingle(c, ok, &scenario)
	})

	scenariosRouter.POST("", func(c *gin.Context) {
		var newScenario Scenario
		if err := c.ShouldBind(&newScenario); err != nil {
			controllers.RespondBadRequest(c, "Invalid scenario data", nil)
			return
		}

		ok := CreateScenario(&newScenario)
		controllers.RespondSingle(c, ok, &newScenario)
	})

	scenariosRouter.PUT("/:scenarioId", func(c *gin.Context) {
		scenarioId, ok := controllers.GetParamAsID(c, "scenarioId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid scenario ID", nil)
			return
		}

		var scenario Scenario
		if err := c.ShouldBind(&scenario); err != nil {
			controllers.RespondBadRequest(c, "Invalid scenario data", nil)
			return
		}

		ok = UpdateScenario(scenarioId, &scenario)
		controllers.RespondSingle(c, ok, &scenario)
	})

	scenariosRouter.DELETE("/:scenarioId", func(c *gin.Context) {
		scenarioId, ok := controllers.GetParamAsID(c, "scenarioId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid scenario ID", nil)
			return
		}

		ok = DeleteScenario(scenarioId)
		controllers.RespondEmpty(c, ok)
	})
}
