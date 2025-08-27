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

		memories, err := GetMemoriesByWorldId(worldId)
		controllers.RespondListE(c, memories, err)
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

		memories, err := GetMemoriesByWorldAndCharacterId(worldId, characterId)
		controllers.RespondListE(c, memories, err)
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

		err := CreateMemory(worldId, &newMemory)
		controllers.RespondSingleE(c, &newMemory, err)
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

		err := UpdateMemory(memoryId, &memory)
		controllers.RespondSingleE(c, &memory, err)
	})

	memoriesRouter.DELETE("/:memoryId", func(c *gin.Context) {
		memoryId, ok := controllers.GetParamAsID(c, "memoryId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid memory ID", nil)
			return
		}

		err := DeleteMemory(memoryId)
		controllers.RespondEmptyE(c, err)
	})
}
