package routes

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/model"
)

func ScenariosController(router *gin.RouterGroup, db *sql.DB) {
	scenariosRouter := router.Group("/scenarios")

	scenariosRouter.GET("", func(c *gin.Context) {
		scenarios, err := model.AllScenarios(db)
		respondList(c, scenarios, err)
	})

	scenariosRouter.GET("/:scenarioId", func(c *gin.Context) {
		scenarioId, err := getIDParam(c, "scenarioId")
		if err != nil {
			respondBadRequest(c, "Invalid scenario ID")
			return
		}

		scenarios, err := model.ScenarioById(db, scenarioId)
		respondSingle(c, scenarios, err)
	})

	scenariosRouter.POST("", func(c *gin.Context) {
		var newScenario model.Scenario
		if err := c.ShouldBind(&newScenario); err != nil {
			respondBadRequest(c, "Invalid scenario data")
			return
		}

		err := model.CreateScenario(db, &newScenario)
		respondSingle(c, &newScenario, err)
	})

	scenariosRouter.PUT("/:scenarioId", func(c *gin.Context) {
		scenarioId, err := getIDParam(c, "scenarioId")
		if err != nil {
			respondBadRequest(c, "Invalid scenario ID")
			return
		}

		var scenario model.Scenario
		if err = c.ShouldBind(&scenario); err != nil {
			respondBadRequest(c, "Invalid scenario data")
			return
		}

		err = model.UpdateScenario(db, scenarioId, &scenario)
		respondSingle(c, &scenario, err)
	})

	scenariosRouter.DELETE("/:scenarioId", func(c *gin.Context) {
		scenarioId, err := getIDParam(c, "scenarioId")
		if err != nil {
			respondBadRequest(c, "Invalid scenario ID")
			return
		}

		err = model.DeleteScenario(db, scenarioId)
		respondDeleted(c, err)
	})
}
