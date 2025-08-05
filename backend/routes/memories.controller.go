package routes

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/model"
	"juraji.nl/chat-quest/util"
)

func MemoriesController(router *gin.RouterGroup, db *sql.DB) {
	memoriesRouter := router.Group("/worlds/:worldId/memories")

	memoriesRouter.GET("", func(c *gin.Context) {
		worldId, err := util.GetIDParam(c, "worldId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid world ID")
			return
		}

		memories, err := model.GetMemoriesByWorldId(db, worldId)
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

		memories, err := model.GetMemoriesByWorldAndCharacterId(db, worldId, characterId)
		util.RespondList(c, memories, err)
	})

	memoriesRouter.POST("", func(c *gin.Context) {
		worldId, err := util.GetIDParam(c, "worldId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid world ID")
			return
		}

		var newMemory model.Memory
		if err := c.Bind(&newMemory); err != nil {
			util.RespondBadRequest(c, "Invalid memory data")
			return
		}

		newMemory.WorldId = worldId

		// TODO: Generate embeddings!
		err = model.CreateMemory(db, &newMemory)
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

		var memory model.Memory
		if err := c.Bind(&memory); err != nil {
			util.RespondBadRequest(c, "Invalid memory data")
			return
		}

		memory.WorldId = worldId

		// TODO: Generate embeddings!
		err = model.UpdateMemory(db, memoryId, &memory)
		util.RespondSingle(c, &memory, err)
	})

	memoriesRouter.DELETE("/:memoryId", func(c *gin.Context) {
		memoryId, err := util.GetIDParam(c, "memoryId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid memory ID")
			return
		}

		err = model.DeleteMemory(db, memoryId)
		util.RespondDeleted(c, err)
	})
}
