package routes

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/model"
)

func MemoriesController(router *gin.RouterGroup, db *sql.DB) {
	memoriesRouter := router.Group("/worlds/:worldId/memories")

	memoriesRouter.GET("", func(c *gin.Context) {
		worldId, err := getIDParam(c, "worldId")
		if err != nil {
			respondBadRequest(c, "Invalid world ID")
			return
		}

		memories, err := model.GetMemoriesByWorldId(db, worldId)
		respondList(c, memories, err)
	})

	memoriesRouter.GET("/by-character/:characterId", func(c *gin.Context) {
		worldId, err := getIDParam(c, "worldId")
		if err != nil {
			respondBadRequest(c, "Invalid world ID")
			return
		}

		characterId, err := getIDParam(c, "characterId")
		if err != nil {
			respondBadRequest(c, "Invalid character ID")
			return
		}

		memories, err := model.GetMemoriesByWorldAndCharacterId(db, worldId, characterId)
		respondList(c, memories, err)
	})

	memoriesRouter.POST("", func(c *gin.Context) {
		worldId, err := getIDParam(c, "worldId")
		if err != nil {
			respondBadRequest(c, "Invalid world ID")
			return
		}

		var newMemory model.Memory
		if err := c.Bind(&newMemory); err != nil {
			respondBadRequest(c, "Invalid memory data")
			return
		}

		newMemory.WorldId = worldId

		// TODO: Generate embeddings!
		err = model.CreateMemory(db, &newMemory)
		respondSingle(c, &newMemory, err)
	})

	memoriesRouter.PUT("/:memoryId", func(c *gin.Context) {
		worldId, err := getIDParam(c, "worldId")
		if err != nil {
			respondBadRequest(c, "Invalid world ID")
			return
		}
		memoryId, err := getIDParam(c, "memoryId")
		if err != nil {
			respondBadRequest(c, "Invalid memory ID")
			return
		}

		var memory model.Memory
		if err := c.Bind(&memory); err != nil {
			respondBadRequest(c, "Invalid memory data")
			return
		}

		memory.WorldId = worldId

		// TODO: Generate embeddings!
		err = model.UpdateMemory(db, memoryId, &memory)
		respondSingle(c, &memory, err)
	})

	memoriesRouter.DELETE("/:memoryId", func(c *gin.Context) {
		memoryId, err := getIDParam(c, "memoryId")
		if err != nil {
			respondBadRequest(c, "Invalid memory ID")
			return
		}

		err = model.DeleteMemory(db, memoryId)
		respondDeleted(c, err)
	})
}
