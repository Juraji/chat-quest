package worlds

import (
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/core/util/controllers"
)

func Routes(router *gin.RouterGroup) {
	preferencesRoutes(router)
	worldsRoutes(router)
}

func preferencesRoutes(router *gin.RouterGroup) {
	prefsRouter := router.Group("/worlds/preferences")

	prefsRouter.GET("", func(c *gin.Context) {
		prefs, ok := GetChatPreferences()
		controllers.RespondSingle(c, ok, prefs)
	})

	prefsRouter.PUT("", func(c *gin.Context) {
		var update ChatPreferences
		if err := c.ShouldBind(&update); err != nil {
			controllers.RespondBadRequest(c, "Invalid preference data")
			return
		}

		ok := UpdateChatPreferences(&update)
		controllers.RespondSingle(c, ok, &update)
	})

	prefsRouter.GET("/is-valid", func(c *gin.Context) {
		var messages []string

		prefs, ok := GetChatPreferences()
		if !ok {
			controllers.RespondInternalError(c, nil)
			return
		}

		if prefs.ChatModelID == nil {
			messages = append(messages, "No chat model set")
		}
		if prefs.ChatInstructionID == nil {
			messages = append(messages, "No chat instruction set")
		}

		controllers.RespondSingle(c, true, &messages)
	})
}

func worldsRoutes(router *gin.RouterGroup) {
	worldsRouter := router.Group("/worlds")

	worldsRouter.GET("", func(c *gin.Context) {
		worlds, ok := GetAllWorlds()
		controllers.RespondList(c, ok, worlds)
	})

	worldsRouter.GET("/:worldId", func(c *gin.Context) {
		worldId, ok := controllers.GetParamAsID(c, "worldId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid world ID")
			return
		}

		world, ok := WorldById(worldId)
		controllers.RespondSingle(c, ok, world)
	})

	worldsRouter.POST("", func(c *gin.Context) {
		var newWorld World
		if err := c.ShouldBind(&newWorld); err != nil {
			controllers.RespondBadRequest(c, "Invalid world data")
			return
		}

		ok := CreateWorld(&newWorld)
		controllers.RespondSingle(c, ok, &newWorld)
	})

	worldsRouter.PUT("/:worldId", func(c *gin.Context) {
		worldId, ok := controllers.GetParamAsID(c, "worldId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid world ID")
			return
		}

		var world World
		if err := c.ShouldBind(&world); err != nil {
			controllers.RespondBadRequest(c, "Invalid world data")
			return
		}

		ok = UpdateWorld(worldId, &world)
		controllers.RespondSingle(c, ok, &world)
	})

	worldsRouter.DELETE("/:worldId", func(c *gin.Context) {
		worldId, ok := controllers.GetParamAsID(c, "worldId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid world ID")
			return
		}

		ok = DeleteWorld(worldId)
		controllers.RespondEmpty(c, ok)
	})
}
