package api

import (
	"github.com/gin-gonic/gin"
	worlds2 "juraji.nl/chat-quest/model/worlds"
)

func WorldsRoutes(router *gin.RouterGroup) {
	worldsRouter := router.Group("/worlds")

	worldsRouter.GET("", func(c *gin.Context) {
		worlds, err := worlds2.GetAllWorlds()
		respondList(c, worlds, err)
	})

	worldsRouter.GET("/:worldId", func(c *gin.Context) {
		worldId, ok := getParamAsID(c, "worldId")
		if !ok {
			respondBadRequest(c, "Invalid world ID", nil)
			return
		}

		world, err := worlds2.WorldById(worldId)
		respondSingle(c, world, err)
	})

	worldsRouter.POST("", func(c *gin.Context) {
		var newWorld worlds2.World
		if err := c.ShouldBind(&newWorld); err != nil {
			respondBadRequest(c, "Invalid world data", nil)
			return
		}

		err := worlds2.CreateWorld(&newWorld)
		respondSingle(c, &newWorld, err)
	})

	worldsRouter.PUT("/:worldId", func(c *gin.Context) {
		worldId, ok := getParamAsID(c, "worldId")
		if !ok {
			respondBadRequest(c, "Invalid world ID", nil)
			return
		}

		var world worlds2.World
		if err := c.ShouldBind(&world); err != nil {
			respondBadRequest(c, "Invalid world data", nil)
			return
		}

		err := worlds2.UpdateWorld(worldId, &world)
		respondSingle(c, &world, err)
	})

	worldsRouter.DELETE("/:worldId", func(c *gin.Context) {
		worldId, ok := getParamAsID(c, "worldId")
		if !ok {
			respondBadRequest(c, "Invalid world ID", nil)
			return
		}

		err := worlds2.DeleteWorld(worldId)
		respondEmpty(c, err)
	})
}
