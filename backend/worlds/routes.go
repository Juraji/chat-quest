package worlds

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/util"
)

func Routes(router *gin.RouterGroup, db *sql.DB) {
	preferencesRoutes(router, db)
	worldsRoutes(router, db)
}

func preferencesRoutes(router *gin.RouterGroup, db *sql.DB) {
	prefsRouter := router.Group("/worlds/preferences")

	prefsRouter.GET("", func(c *gin.Context) {
		prefs, err := GetChatPreferences(db)
		util.RespondSingle(c, prefs, err)
	})

	prefsRouter.PUT("", func(c *gin.Context) {
		var update ChatPreferences
		if err := c.ShouldBind(&update); err != nil {
			util.RespondBadRequest(c, "Invalid preference data")
			return
		}

		err := UpdateChatPreferences(db, &update)
		util.RespondSingle(c, &update, err)
	})

	prefsRouter.GET("/is-valid", func(c *gin.Context) {
		var messages []string

		prefs, err := GetChatPreferences(db)
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

func worldsRoutes(router *gin.RouterGroup, db *sql.DB) {
	worldsRouter := router.Group("/worlds")

	worldsRouter.GET("", func(c *gin.Context) {
		worlds, err := GetAllWorlds(db)
		util.RespondList(c, worlds, err)
	})

	worldsRouter.GET("/:worldId", func(c *gin.Context) {
		worldId, err := util.GetIDParam(c, "worldId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid world ID")
			return
		}

		world, err := WorldById(db, worldId)
		util.RespondSingle(c, world, err)
	})

	worldsRouter.POST("", func(c *gin.Context) {
		var newWorld World
		if err := c.ShouldBind(&newWorld); err != nil {
			util.RespondBadRequest(c, "Invalid world data")
			return
		}

		err := CreateWorld(db, &newWorld)
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

		err = UpdateWorld(db, worldId, &world)
		util.RespondSingle(c, &world, err)
	})

	worldsRouter.DELETE("/:worldId", func(c *gin.Context) {
		worldId, err := util.GetIDParam(c, "worldId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid world ID")
			return
		}

		err = DeleteWorld(db, worldId)
		util.RespondDeleted(c, err)
	})
}
