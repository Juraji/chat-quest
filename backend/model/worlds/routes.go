package worlds

import (
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/core"
	"juraji.nl/chat-quest/core/util"
)

func Routes(cq *core.ChatQuestContext, router *gin.RouterGroup) {
	preferencesRoutes(cq, router)
	worldsRoutes(cq, router)
}

func preferencesRoutes(cq *core.ChatQuestContext, router *gin.RouterGroup) {
	prefsRouter := router.Group("/worlds/preferences")

	prefsRouter.GET("", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		prefs, err := GetChatPreferences(cq)
		util.RespondSingle(rcq, c, prefs, err)
	})

	prefsRouter.PUT("", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		var update ChatPreferences
		if err := c.ShouldBind(&update); err != nil {
			util.RespondBadRequest(rcq, c, "Invalid preference data")
			return
		}

		err := UpdateChatPreferences(rcq, &update)
		util.RespondSingle(rcq, c, &update, err)
	})

	prefsRouter.GET("/is-valid", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		var messages []string

		prefs, err := GetChatPreferences(cq)
		if err != nil {
			util.RespondInternalError(rcq, c, err)
			return
		}

		if prefs.ChatModelID == nil {
			messages = append(messages, "No chat model set")
		}
		if prefs.ChatInstructionID == nil {
			messages = append(messages, "No chat instruction set")
		}

		util.RespondSingle(rcq, c, &messages, nil)
	})
}

func worldsRoutes(cq *core.ChatQuestContext, router *gin.RouterGroup) {
	worldsRouter := router.Group("/worlds")

	worldsRouter.GET("", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		worlds, err := GetAllWorlds(cq)
		util.RespondList(rcq, c, worlds, err)
	})

	worldsRouter.GET("/:worldId", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		worldId, err := util.GetIDParam(c, "worldId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid world ID")
			return
		}

		world, err := WorldById(rcq, worldId)
		util.RespondSingle(rcq, c, world, err)
	})

	worldsRouter.POST("", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		var newWorld World
		if err := c.ShouldBind(&newWorld); err != nil {
			util.RespondBadRequest(rcq, c, "Invalid world data")
			return
		}

		err := CreateWorld(rcq, &newWorld)
		util.RespondSingle(rcq, c, &newWorld, err)
	})

	worldsRouter.PUT("/:worldId", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		worldId, err := util.GetIDParam(c, "worldId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid world ID")
			return
		}

		var world World
		if err := c.ShouldBind(&world); err != nil {
			util.RespondBadRequest(rcq, c, "Invalid world data")
			return
		}

		err = UpdateWorld(rcq, worldId, &world)
		util.RespondSingle(rcq, c, &world, err)
	})

	worldsRouter.DELETE("/:worldId", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		worldId, err := util.GetIDParam(c, "worldId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid world ID")
			return
		}

		err = DeleteWorld(rcq, worldId)
		util.RespondDeleted(rcq, c, err)
	})
}
