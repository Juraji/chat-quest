package memories

import (
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/core/util/controllers"
)

func Routes(router *gin.RouterGroup) {
	memoriesRoutes(router)
	prefsRoutes(router)
}

func memoriesRoutes(router *gin.RouterGroup) {
	memoriesRouter := router.Group("/worlds/:worldId/memories")

	memoriesRouter.GET("", func(c *gin.Context) {
		worldId, ok := controllers.GetParamAsID(c, "worldId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid world ID")
			return
		}

		memories, ok := GetMemoriesByWorldId(worldId)
		controllers.RespondList(c, ok, memories)
	})

	memoriesRouter.GET("/by-character/:characterId", func(c *gin.Context) {
		worldId, ok := controllers.GetParamAsID(c, "worldId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid world ID")
			return
		}

		characterId, ok := controllers.GetParamAsID(c, "characterId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid character ID")
			return
		}

		memories, ok := GetMemoriesByWorldAndCharacterId(worldId, characterId)
		controllers.RespondList(c, ok, memories)
	})

	memoriesRouter.POST("", func(c *gin.Context) {
		worldId, ok := controllers.GetParamAsID(c, "worldId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid world ID")
			return
		}

		var newMemory Memory
		if err := c.Bind(&newMemory); err != nil {
			controllers.RespondBadRequest(c, "Invalid memory data")
			return
		}

		newMemory.WorldId = worldId

		// TODO: Generate embeddings!
		ok = CreateMemory(&newMemory)
		controllers.RespondSingle(c, ok, &newMemory)
	})

	memoriesRouter.PUT("/:memoryId", func(c *gin.Context) {
		worldId, ok := controllers.GetParamAsID(c, "worldId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid world ID")
			return
		}
		memoryId, ok := controllers.GetParamAsID(c, "memoryId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid memory ID")
			return
		}

		var memory Memory
		if err := c.Bind(&memory); err != nil {
			controllers.RespondBadRequest(c, "Invalid memory data")
			return
		}

		memory.WorldId = worldId

		// TODO: Generate embeddings!
		ok = UpdateMemory(memoryId, &memory)
		controllers.RespondSingle(c, ok, &memory)
	})

	memoriesRouter.DELETE("/:memoryId", func(c *gin.Context) {
		memoryId, ok := controllers.GetParamAsID(c, "memoryId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid memory ID")
			return
		}

		ok = DeleteMemory(memoryId)
		controllers.RespondEmpty(c, ok)
	})
}

func prefsRoutes(router *gin.RouterGroup) {
	prefsRouter := router.Group("/memories/preferences")

	prefsRouter.GET("", func(c *gin.Context) {
		prefs, ok := GetMemoryPreferences()
		controllers.RespondSingle(c, ok, prefs)
	})

	prefsRouter.PUT("", func(c *gin.Context) {
		var update MemoryPreferences
		if err := c.ShouldBind(&update); err != nil {
			controllers.RespondBadRequest(c, "Invalid preference data")
			return
		}

		ok := UpdateMemoryPreferences(&update)
		controllers.RespondSingle(c, ok, &update)
	})

	prefsRouter.GET("/is-valid", func(c *gin.Context) {
		var messages []string

		prefs, ok := GetMemoryPreferences()
		if !ok {
			controllers.RespondInternalError(c, nil)
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

		controllers.RespondSingle(c, true, &messages)
	})
}
