package scenarios

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/util"
)

func Routes(router *gin.RouterGroup, db *sql.DB) {
	scenariosRouter := router.Group("/scenarios")

	scenariosRouter.GET("", func(c *gin.Context) {
		scenarios, err := AllScenarios(db)
		util.RespondList(c, scenarios, err)
	})

	scenariosRouter.GET("/:scenarioId", func(c *gin.Context) {
		scenarioId, err := util.GetIDParam(c, "scenarioId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid scenario ID")
			return
		}

		scenarios, err := ScenarioById(db, scenarioId)
		util.RespondSingle(c, scenarios, err)
	})

	scenariosRouter.POST("", func(c *gin.Context) {
		var newScenario Scenario
		if err := c.ShouldBind(&newScenario); err != nil {
			util.RespondBadRequest(c, "Invalid scenario data")
			return
		}

		err := CreateScenario(db, &newScenario)
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

		err = UpdateScenario(db, scenarioId, &scenario)
		util.RespondSingle(c, &scenario, err)
	})

	scenariosRouter.DELETE("/:scenarioId", func(c *gin.Context) {
		scenarioId, err := util.GetIDParam(c, "scenarioId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid scenario ID")
			return
		}

		err = DeleteScenario(db, scenarioId)
		util.RespondDeleted(c, err)
	})
}
