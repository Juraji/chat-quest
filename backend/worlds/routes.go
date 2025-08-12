package worlds

import (
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/cq"
	"juraji.nl/chat-quest/util"
)

func Routes(cq *cq.ChatQuestContext, router *gin.RouterGroup) {
	preferencesRoutes(cq, router)
	worldsRoutes(cq, router)
}

func preferencesRoutes(cq *cq.ChatQuestContext, router *gin.RouterGroup) {
	prefsRouter := router.Group("/worlds/preferences")

	prefsRouter.GET("", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		prefs, err := GetChatPreferences(cq)
		util.RespondSingle(cq, c, prefs, err)
	})

	prefsRouter.PUT("", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		var update ChatPreferences
		if err := c.ShouldBind(&update); err != nil {
			util.RespondBadRequest(cq, c, "Invalid preference data")
			return
		}

		err := UpdateChatPreferences(cq, &update)
		util.RespondSingle(cq, c, &update, err)
	})

	prefsRouter.GET("/is-valid", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		var messages []string

		prefs, err := GetChatPreferences(cq)
		if err != nil {
			util.RespondInternalError(cq, c, err)
			return
		}

		if prefs.ChatModelID == nil {
			messages = append(messages, "No chat model set")
		}
		if prefs.ChatInstructionID == nil {
			messages = append(messages, "No chat instruction set")
		}

		util.RespondSingle(cq, c, &messages, nil)
	})
}

func worldsRoutes(cq *cq.ChatQuestContext, router *gin.RouterGroup) {
	worldsRouter := router.Group("/worlds")

	worldsRouter.GET("", func(c *gin.Context) {
		worlds, err := GetAllWorlds(cq)
		util.RespondList(cq, c, worlds, err)
	})

	worldsRouter.GET("/:worldId", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		worldId, err := util.GetIDParam(c, "worldId")
		if err != nil {
			util.RespondBadRequest(cq, c, "Invalid world ID")
			return
		}

		world, err := WorldById(cq, worldId)
		util.RespondSingle(cq, c, world, err)
	})

	worldsRouter.POST("", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		var newWorld World
		if err := c.ShouldBind(&newWorld); err != nil {
			util.RespondBadRequest(cq, c, "Invalid world data")
			return
		}

		err := CreateWorld(cq, &newWorld)
		util.RespondSingle(cq, c, &newWorld, err)
	})

	worldsRouter.PUT("/:worldId", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		worldId, err := util.GetIDParam(c, "worldId")
		if err != nil {
			util.RespondBadRequest(cq, c, "Invalid world ID")
			return
		}

		var world World
		if err := c.ShouldBind(&world); err != nil {
			util.RespondBadRequest(cq, c, "Invalid world data")
			return
		}

		err = UpdateWorld(cq, worldId, &world)
		util.RespondSingle(cq, c, &world, err)
	})

	worldsRouter.DELETE("/:worldId", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		worldId, err := util.GetIDParam(c, "worldId")
		if err != nil {
			util.RespondBadRequest(cq, c, "Invalid world ID")
			return
		}

		err = DeleteWorld(cq, worldId)
		util.RespondDeleted(cq, c, err)
	})
}
