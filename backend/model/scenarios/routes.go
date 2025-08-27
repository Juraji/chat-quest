package scenarios

import (
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/core/util/controllers"
)

func Routes(router *gin.RouterGroup) {
	scenariosRouter := router.Group("/scenarios")

	scenariosRouter.GET("", func(c *gin.Context) {
		scenarios, err := AllScenarios()
		controllers.RespondList(c, scenarios, err)
	})

	scenariosRouter.GET("/:scenarioId", func(c *gin.Context) {
		scenarioId, ok := controllers.GetParamAsID(c, "scenarioId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid scenario ID", nil)
			return
		}

		scenario, err := ScenarioById(scenarioId)
		controllers.RespondSingle(c, &scenario, err)
	})

	scenariosRouter.POST("", func(c *gin.Context) {
		var newScenario Scenario
		if err := c.ShouldBind(&newScenario); err != nil {
			controllers.RespondBadRequest(c, "Invalid scenario data", nil)
			return
		}

		err := CreateScenario(&newScenario)
		controllers.RespondSingle(c, &newScenario, err)
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

		err := UpdateScenario(scenarioId, &scenario)
		controllers.RespondSingle(c, &scenario, err)
	})

	scenariosRouter.DELETE("/:scenarioId", func(c *gin.Context) {
		scenarioId, ok := controllers.GetParamAsID(c, "scenarioId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid scenario ID", nil)
			return
		}

		err := DeleteScenario(scenarioId)
		controllers.RespondEmpty(c, err)
	})
}
