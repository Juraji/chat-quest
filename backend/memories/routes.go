package memories

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/util"
)

func Routes(router *gin.RouterGroup, db *sql.DB) {
	memoriesRoutes(router, db)
	prefsRoutes(router, db)
}

func memoriesRoutes(router *gin.RouterGroup, db *sql.DB) {
	memoriesRouter := router.Group("/worlds/:worldId/memories")

	memoriesRouter.GET("", func(c *gin.Context) {
		worldId, err := util.GetIDParam(c, "worldId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid world ID")
			return
		}

		memories, err := GetMemoriesByWorldId(db, worldId)
		util.RespondList(c, memories, err)
	})

	memoriesRouter.GET("/by-character/:characterId", func(c *gin.Context) {
		worldId, err := util.GetIDParam(c, "worldId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid world ID")
			return
		}

		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid character ID")
			return
		}

		memories, err := GetMemoriesByWorldAndCharacterId(db, worldId, characterId)
		util.RespondList(c, memories, err)
	})

	memoriesRouter.POST("", func(c *gin.Context) {
		worldId, err := util.GetIDParam(c, "worldId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid world ID")
			return
		}

		var newMemory Memory
		if err := c.Bind(&newMemory); err != nil {
			util.RespondBadRequest(c, "Invalid memory data")
			return
		}

		newMemory.WorldId = worldId

		// TODO: Generate embeddings!
		err = CreateMemory(db, &newMemory)
		util.RespondSingle(c, &newMemory, err)
	})

	memoriesRouter.PUT("/:memoryId", func(c *gin.Context) {
		worldId, err := util.GetIDParam(c, "worldId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid world ID")
			return
		}
		memoryId, err := util.GetIDParam(c, "memoryId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid memory ID")
			return
		}

		var memory Memory
		if err := c.Bind(&memory); err != nil {
			util.RespondBadRequest(c, "Invalid memory data")
			return
		}

		memory.WorldId = worldId

		// TODO: Generate embeddings!
		err = UpdateMemory(db, memoryId, &memory)
		util.RespondSingle(c, &memory, err)
	})

	memoriesRouter.DELETE("/:memoryId", func(c *gin.Context) {
		memoryId, err := util.GetIDParam(c, "memoryId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid memory ID")
			return
		}

		err = DeleteMemory(db, memoryId)
		util.RespondDeleted(c, err)
	})
}

func prefsRoutes(router *gin.RouterGroup, db *sql.DB) {
	prefsRouter := router.Group("/memories/preferences")

	prefsRouter.GET("", func(c *gin.Context) {
		prefs, err := GetMemoryPreferences(db)
		util.RespondSingle(c, prefs, err)
	})

	prefsRouter.PUT("", func(c *gin.Context) {
		var update MemoryPreferences
		if err := c.ShouldBind(&update); err != nil {
			util.RespondBadRequest(c, "Invalid preference data")
			return
		}

		err := UpdateMemoryPreferences(db, &update)
		util.RespondSingle(c, &update, err)
	})

	prefsRouter.GET("/is-valid", func(c *gin.Context) {
		var messages []string

		prefs, err := GetMemoryPreferences(db)
		if err != nil {
			util.RespondInternalError(c, err)
			return
		}

		if prefs.MemoriesModelID == nil {
			messages = append(messages, "No memory model set")
		}
		if prefs.MemoriesInstructionID == nil {
			messages = append(messages, "No memory instruction set")
		}
		if prefs.EmbeddingModelID == nil {
			messages = append(messages, "No memory embedding model set")
		}

		util.RespondSingle(c, &messages, nil)
	})
}
