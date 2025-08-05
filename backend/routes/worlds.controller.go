package routes

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/model"
	"juraji.nl/chat-quest/util"
)

func WorldsController(router *gin.RouterGroup, db *sql.DB) {
	worldsRouter := router.Group("/worlds")

	worldsRouter.GET("", func(c *gin.Context) {
		worlds, err := model.GetAllWorlds(db)
		util.RespondList(c, worlds, err)
	})

	worldsRouter.GET("/:worldId", func(c *gin.Context) {
		worldId, err := util.GetIDParam(c, "worldId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid world ID")
			return
		}

		world, err := model.WorldById(db, worldId)
		util.RespondSingle(c, world, err)
	})

	worldsRouter.POST("", func(c *gin.Context) {
		var newWorld model.World
		if err := c.ShouldBind(&newWorld); err != nil {
			util.RespondBadRequest(c, "Invalid world data")
			return
		}

		err := model.CreateWorld(db, &newWorld)
		util.RespondSingle(c, &newWorld, err)
	})

	worldsRouter.PUT("/:worldId", func(c *gin.Context) {
		worldId, err := util.GetIDParam(c, "worldId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid world ID")
			return
		}

		var world model.World
		if err := c.ShouldBind(&world); err != nil {
			util.RespondBadRequest(c, "Invalid world data")
			return
		}

		err = model.UpdateWorld(db, worldId, &world)
		util.RespondSingle(c, &world, err)
	})

	worldsRouter.DELETE("/:worldId", func(c *gin.Context) {
		worldId, err := util.GetIDParam(c, "worldId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid world ID")
			return
		}

		err = model.DeleteWorld(db, worldId)
		util.RespondDeleted(c, err)
	})
}
