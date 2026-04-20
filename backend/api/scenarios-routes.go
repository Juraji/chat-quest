package api

import (
	"github.com/gin-gonic/gin"
	scenarios2 "juraji.nl/chat-quest/model/scenarios"
)

func ScenariosRoutes(router *gin.RouterGroup) {
	scenariosRouter := router.Group("/scenarios")

	scenariosRouter.GET("", func(c *gin.Context) {
		scenarios, err := scenarios2.AllScenarios()
		respondList(c, scenarios, err)
	})

	scenariosRouter.GET("/:scenarioId", func(c *gin.Context) {
		scenarioId, ok := getParamAsID(c, "scenarioId")
		if !ok {
			respondBadRequest(c, "Invalid scenario ID", nil)
			return
		}

		scenario, err := scenarios2.ScenarioById(scenarioId)
		respondSingle(c, &scenario, err)
	})

	scenariosRouter.POST("", func(c *gin.Context) {
		var newScenario scenarios2.Scenario
		if err := c.ShouldBind(&newScenario); err != nil {
			respondBadRequest(c, "Invalid scenario data", nil)
			return
		}

		err := scenarios2.CreateScenario(&newScenario)
		respondSingle(c, &newScenario, err)
	})

	scenariosRouter.PUT("/:scenarioId", func(c *gin.Context) {
		scenarioId, ok := getParamAsID(c, "scenarioId")
		if !ok {
			respondBadRequest(c, "Invalid scenario ID", nil)
			return
		}

		var scenario scenarios2.Scenario
		if err := c.ShouldBind(&scenario); err != nil {
			respondBadRequest(c, "Invalid scenario data", nil)
			return
		}

		err := scenarios2.UpdateScenario(scenarioId, &scenario)
		respondSingle(c, &scenario, err)
	})

	scenariosRouter.DELETE("/:scenarioId", func(c *gin.Context) {
		scenarioId, ok := getParamAsID(c, "scenarioId")
		if !ok {
			respondBadRequest(c, "Invalid scenario ID", nil)
			return
		}

		err := scenarios2.DeleteScenario(scenarioId)
		respondEmpty(c, err)
	})
}
