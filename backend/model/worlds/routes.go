package worlds

import (
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/core/util/controllers"
)

func Routes(router *gin.RouterGroup) {
	worldsRouter := router.Group("/worlds")

	worldsRouter.GET("", func(c *gin.Context) {
		worlds, err := GetAllWorlds()
		controllers.RespondList(c, worlds, err)
	})

	worldsRouter.GET("/:worldId", func(c *gin.Context) {
		worldId, ok := controllers.GetParamAsID(c, "worldId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid world ID", nil)
			return
		}

		world, err := WorldById(worldId)
		controllers.RespondSingle(c, world, err)
	})

	worldsRouter.POST("", func(c *gin.Context) {
		var newWorld World
		if err := c.ShouldBind(&newWorld); err != nil {
			controllers.RespondBadRequest(c, "Invalid world data", nil)
			return
		}

		err := CreateWorld(&newWorld)
		controllers.RespondSingle(c, &newWorld, err)
	})

	worldsRouter.PUT("/:worldId", func(c *gin.Context) {
		worldId, ok := controllers.GetParamAsID(c, "worldId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid world ID", nil)
			return
		}

		var world World
		if err := c.ShouldBind(&world); err != nil {
			controllers.RespondBadRequest(c, "Invalid world data", nil)
			return
		}

		err := UpdateWorld(worldId, &world)
		controllers.RespondSingle(c, &world, err)
	})

	worldsRouter.DELETE("/:worldId", func(c *gin.Context) {
		worldId, ok := controllers.GetParamAsID(c, "worldId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid world ID", nil)
			return
		}

		err := DeleteWorld(worldId)
		controllers.RespondEmpty(c, err)
	})
}
