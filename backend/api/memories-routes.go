package api

import (
	"github.com/gin-gonic/gin"
	m "juraji.nl/chat-quest/model/memories"
	"juraji.nl/chat-quest/processing"
)

func MemoriesRoutes(router *gin.RouterGroup) {
	memoriesRouter := router.Group("/worlds/:worldId/memories")

	memoriesRouter.GET("", func(c *gin.Context) {
		worldId, ok := getParamAsID(c, "worldId")
		if !ok {
			respondBadRequest(c, "Invalid world ID", nil)
			return
		}

		memories, err := m.GetMemoriesByWorldId(worldId)
		respondList(c, memories, err)
	})

	memoriesRouter.GET("/by-character/:characterId", func(c *gin.Context) {
		worldId, ok := getParamAsID(c, "worldId")
		if !ok {
			respondBadRequest(c, "Invalid world ID", nil)
			return
		}

		characterId, ok := getParamAsID(c, "characterId")
		if !ok {
			respondBadRequest(c, "Invalid character ID", nil)
			return
		}

		memories, err := m.GetMemoriesByWorldAndCharacterId(worldId, characterId)
		respondList(c, memories, err)
	})

	memoriesRouter.POST("", func(c *gin.Context) {
		worldId, ok := getParamAsID(c, "worldId")
		if !ok {
			respondBadRequest(c, "Invalid world ID", nil)
			return
		}

		var newMemory m.Memory
		if err := c.Bind(&newMemory); err != nil {
			respondBadRequest(c, "Invalid memory data", nil)
			return
		}

		err := m.CreateMemory(worldId, &newMemory)
		respondSingle(c, &newMemory, err)
	})

	memoriesRouter.PUT("/:memoryId", func(c *gin.Context) {
		worldId, ok := getParamAsID(c, "worldId")
		if !ok {
			respondBadRequest(c, "Invalid world ID", nil)
			return
		}
		memoryId, ok := getParamAsID(c, "memoryId")
		if !ok {
			respondBadRequest(c, "Invalid memory ID", nil)
			return
		}

		var memory m.Memory
		if err := c.Bind(&memory); err != nil {
			respondBadRequest(c, "Invalid memory data", nil)
			return
		}

		memory.WorldId = worldId

		err := m.UpdateMemory(memoryId, &memory)
		respondSingle(c, &memory, err)
	})

	memoriesRouter.DELETE("/:memoryId", func(c *gin.Context) {
		memoryId, ok := getParamAsID(c, "memoryId")
		if !ok {
			respondBadRequest(c, "Invalid memory ID", nil)
			return
		}

		err := m.DeleteMemory(memoryId)
		respondEmpty(c, err)
	})

	memoriesRouter.GET("/bookmarks/:chatSessionId", func(c *gin.Context) {
		chatSessionId, ok := getParamAsID(c, "chatSessionId")
		if !ok {
			respondBadRequest(c, "Invalid chat session ID", nil)
			return
		}

		bookmarkID, err := m.GetMemoryBookmark(chatSessionId)
		if bookmarkID == nil {
			respondEmpty(c, err)
		} else {
			respondSingle(c, bookmarkID, err)
		}
	})

	memoriesRouter.POST("/generate/:messageId", func(c *gin.Context) {
		messageId, ok := getParamAsID(c, "messageId")
		if !ok {
			respondBadRequest(c, "Invalid message ID", nil)
			return
		}

		includeNPreceding := getQueryParamAsIntOr(c, "includeNPreceding", 0)

		err := processing.GenerateMemoriesForMessageID(c, messageId, includeNPreceding)
		respondEmpty(c, err)
	})
}
