package worlds

import (
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/core/util"
)

func Routes(router *gin.RouterGroup) {
	preferencesRoutes(router)
	worldsRoutes(router)
}

func preferencesRoutes(router *gin.RouterGroup) {
	prefsRouter := router.Group("/worlds/preferences")

	prefsRouter.GET("", func(c *gin.Context) {
		prefs, err := GetChatPreferences()
		util.RespondSingle(c, prefs, err)
	})

	prefsRouter.PUT("", func(c *gin.Context) {
		var update ChatPreferences
		if err := c.ShouldBind(&update); err != nil {
			util.RespondBadRequest(c, "Invalid preference data")
			return
		}

		err := UpdateChatPreferences(&update)
		util.RespondSingle(c, &update, err)
	})

	prefsRouter.GET("/is-valid", func(c *gin.Context) {
		var messages []string

		prefs, err := GetChatPreferences()
		if err != nil {
			util.RespondInternalError(c, err)
			return
		}

		if prefs.ChatModelID == nil {
			messages = append(messages, "No chat model set")
		}
		if prefs.ChatInstructionID == nil {
			messages = append(messages, "No chat instruction set")
		}

		util.RespondSingle(c, &messages, nil)
	})
}

func worldsRoutes(router *gin.RouterGroup) {
	worldsRouter := router.Group("/worlds")

	worldsRouter.GET("", func(c *gin.Context) {
		worlds, err := GetAllWorlds()
		util.RespondList(c, worlds, err)
	})

	worldsRouter.GET("/:worldId", func(c *gin.Context) {
		worldId, err := util.GetIDParam(c, "worldId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid world ID")
			return
		}

		world, err := WorldById(worldId)
		util.RespondSingle(c, world, err)
	})

	worldsRouter.POST("", func(c *gin.Context) {
		var newWorld World
		if err := c.ShouldBind(&newWorld); err != nil {
			util.RespondBadRequest(c, "Invalid world data")
			return
		}

		err := CreateWorld(&newWorld)
		util.RespondSingle(c, &newWorld, err)
	})

	worldsRouter.PUT("/:worldId", func(c *gin.Context) {
		worldId, err := util.GetIDParam(c, "worldId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid world ID")
			return
		}

		var world World
		if err := c.ShouldBind(&world); err != nil {
			util.RespondBadRequest(c, "Invalid world data")
			return
		}

		err = UpdateWorld(worldId, &world)
		util.RespondSingle(c, &world, err)
	})

	worldsRouter.DELETE("/:worldId", func(c *gin.Context) {
		worldId, err := util.GetIDParam(c, "worldId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid world ID")
			return
		}

		err = DeleteWorld(worldId)
		util.RespondDeleted(c, err)
	})
}
