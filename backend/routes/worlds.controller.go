package routes

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/model"
)

func WorldsController(router *gin.RouterGroup, db *sql.DB) {
	worldsRouter := router.Group("/worlds")

	worldsRouter.GET("", func(c *gin.Context) {
		worlds, err := model.GetAllWorlds(db)
		respondList(c, worlds, err)
	})

	worldsRouter.GET("/:worldId", func(c *gin.Context) {
		worldId, err := getIDParam(c, "worldId")
		if err != nil {
			respondBadRequest(c, "Invalid world ID")
			return
		}

		world, err := model.WorldById(db, worldId)
		respondSingle(c, world, err)
	})

	worldsRouter.POST("", func(c *gin.Context) {
		var newWorld model.World
		if err := c.ShouldBind(&newWorld); err != nil {
			respondBadRequest(c, "Invalid world data")
			return
		}

		err := model.CreateWorld(db, &newWorld)
		respondSingle(c, &newWorld, err)
	})

	worldsRouter.PUT("/:worldId", func(c *gin.Context) {
		worldId, err := getIDParam(c, "worldId")
		if err != nil {
			respondBadRequest(c, "Invalid world ID")
			return
		}

		var world model.World
		if err := c.ShouldBind(&world); err != nil {
			respondBadRequest(c, "Invalid world data")
			return
		}

		err = model.UpdateWorld(db, worldId, &world)
		respondSingle(c, &world, err)
	})

	worldsRouter.DELETE("/:worldId", func(c *gin.Context) {
		worldId, err := getIDParam(c, "worldId")
		if err != nil {
			respondBadRequest(c, "Invalid world ID")
			return
		}

		err = model.DeleteWorld(db, worldId)
		respondDeleted(c, err)
	})
}
