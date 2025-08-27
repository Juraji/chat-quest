package memories

import (
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/core/util/controllers"
)

func Routes(router *gin.RouterGroup) {
	memoriesRouter := router.Group("/worlds/:worldId/memories")

	memoriesRouter.GET("", func(c *gin.Context) {
		worldId, ok := controllers.GetParamAsID(c, "worldId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid world ID", nil)
			return
		}

		memories, ok := GetMemoriesByWorldId(worldId)
		controllers.RespondList(c, ok, memories)
	})

	memoriesRouter.GET("/by-character/:characterId", func(c *gin.Context) {
		worldId, ok := controllers.GetParamAsID(c, "worldId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid world ID", nil)
			return
		}

		characterId, ok := controllers.GetParamAsID(c, "characterId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid character ID", nil)
			return
		}

		memories, ok := GetMemoriesByWorldAndCharacterId(worldId, characterId)
		controllers.RespondList(c, ok, memories)
	})

	memoriesRouter.POST("", func(c *gin.Context) {
		worldId, ok := controllers.GetParamAsID(c, "worldId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid world ID", nil)
			return
		}

		var newMemory Memory
		if err := c.Bind(&newMemory); err != nil {
			controllers.RespondBadRequest(c, "Invalid memory data", nil)
			return
		}

		ok = CreateMemory(worldId, &newMemory)
		controllers.RespondSingle(c, ok, &newMemory)
	})

	memoriesRouter.PUT("/:memoryId", func(c *gin.Context) {
		worldId, ok := controllers.GetParamAsID(c, "worldId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid world ID", nil)
			return
		}
		memoryId, ok := controllers.GetParamAsID(c, "memoryId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid memory ID", nil)
			return
		}

		var memory Memory
		if err := c.Bind(&memory); err != nil {
			controllers.RespondBadRequest(c, "Invalid memory data", nil)
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
			controllers.RespondBadRequest(c, "Invalid memory ID", nil)
			return
		}

		ok = DeleteMemory(memoryId)
		controllers.RespondEmpty(c, ok)
	})
}
