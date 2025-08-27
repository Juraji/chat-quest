package worlds

import (
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/core/util/controllers"
)

func Routes(router *gin.RouterGroup) {
	worldsRouter := router.Group("/worlds")

	worldsRouter.GET("", func(c *gin.Context) {
		worlds, ok := GetAllWorlds()
		controllers.RespondList(c, ok, worlds)
	})

	worldsRouter.GET("/:worldId", func(c *gin.Context) {
		worldId, ok := controllers.GetParamAsID(c, "worldId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid world ID", nil)
			return
		}

		world, ok := WorldById(worldId)
		controllers.RespondSingle(c, ok, world)
	})

	worldsRouter.POST("", func(c *gin.Context) {
		var newWorld World
		if err := c.ShouldBind(&newWorld); err != nil {
			controllers.RespondBadRequest(c, "Invalid world data", nil)
			return
		}

		ok := CreateWorld(&newWorld)
		controllers.RespondSingle(c, ok, &newWorld)
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

		ok = UpdateWorld(worldId, &world)
		controllers.RespondSingle(c, ok, &world)
	})

	worldsRouter.DELETE("/:worldId", func(c *gin.Context) {
		worldId, ok := controllers.GetParamAsID(c, "worldId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid world ID", nil)
			return
		}

		ok = DeleteWorld(worldId)
		controllers.RespondEmpty(c, ok)
	})
}
