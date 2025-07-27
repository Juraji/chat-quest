package routes

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/model"
	"net/http"
)

func ScenariosController(router *gin.RouterGroup, db *sql.DB) {
	scenariosRouter := router.Group("/scenarios")

	scenariosRouter.GET("/", func(c *gin.Context) {
		scenarios, err := model.AllScenarios(db)
		respondList(c, scenarios, err)
	})

	scenariosRouter.GET("/:scenarioId", func(c *gin.Context) {
		scenarioId, err := getID(c, "scenarioId")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid scenario ID"})
			return
		}

		scenario, err := model.ScenarioById(db, scenarioId)
		respondSingle(c, scenario, err)
	})

	scenariosRouter.POST("/", func(c *gin.Context) {
		var newScenario model.Scenario
		if err := c.ShouldBind(&newScenario); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid scenario data"})
			return
		}

		err := model.CreateScenario(db, &newScenario)
		respondSingle(c, &newScenario, err)
	})

	scenariosRouter.PUT("/:scenarioId", func(c *gin.Context) {
		scenarioId, err := getID(c, "scenarioId")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid scenario ID"})
			return
		}

		var scenario model.Scenario
		if err := c.ShouldBind(&scenario); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid scenario data"})
			return
		}

		err = model.UpdateScenario(db, scenarioId, &scenario)
		respondSingle(c, &scenario, err)
	})

	scenariosRouter.DELETE("/:scenarioId", func(c *gin.Context) {
		scenarioId, err := getID(c, "scenarioId")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid scenario ID"})
			return
		}

		err = model.DeleteScenario(db, scenarioId)
		respondDeleted(c, err)
	})
}
