package memories

import (
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/cq"
	"juraji.nl/chat-quest/util"
)

func Routes(cq *cq.ChatQuestContext, router *gin.RouterGroup) {
	memoriesRoutes(cq, router)
	prefsRoutes(cq, router)
}

func memoriesRoutes(cq *cq.ChatQuestContext, router *gin.RouterGroup) {
	memoriesRouter := router.Group("/worlds/:worldId/memories")

	memoriesRouter.GET("", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		worldId, err := util.GetIDParam(c, "worldId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid world ID")
			return
		}

		memories, err := GetMemoriesByWorldId(rcq, worldId)
		util.RespondList(rcq, c, memories, err)
	})

	memoriesRouter.GET("/by-character/:characterId", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		worldId, err := util.GetIDParam(c, "worldId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid world ID")
			return
		}

		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid character ID")
			return
		}

		memories, err := GetMemoriesByWorldAndCharacterId(rcq, worldId, characterId)
		util.RespondList(rcq, c, memories, err)
	})

	memoriesRouter.POST("", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		worldId, err := util.GetIDParam(c, "worldId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid world ID")
			return
		}

		var newMemory Memory
		if err := c.Bind(&newMemory); err != nil {
			util.RespondBadRequest(rcq, c, "Invalid memory data")
			return
		}

		newMemory.WorldId = worldId

		// TODO: Generate embeddings!
		err = CreateMemory(rcq, &newMemory)
		util.RespondSingle(rcq, c, &newMemory, err)
	})

	memoriesRouter.PUT("/:memoryId", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		worldId, err := util.GetIDParam(c, "worldId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid world ID")
			return
		}
		memoryId, err := util.GetIDParam(c, "memoryId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid memory ID")
			return
		}

		var memory Memory
		if err := c.Bind(&memory); err != nil {
			util.RespondBadRequest(rcq, c, "Invalid memory data")
			return
		}

		memory.WorldId = worldId

		// TODO: Generate embeddings!
		err = UpdateMemory(rcq, memoryId, &memory)
		util.RespondSingle(rcq, c, &memory, err)
	})

	memoriesRouter.DELETE("/:memoryId", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		memoryId, err := util.GetIDParam(c, "memoryId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid memory ID")
			return
		}

		err = DeleteMemory(rcq, memoryId)
		util.RespondDeleted(rcq, c, err)
	})
}

func prefsRoutes(cq *cq.ChatQuestContext, router *gin.RouterGroup) {
	prefsRouter := router.Group("/memories/preferences")

	prefsRouter.GET("", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		prefs, err := GetMemoryPreferences(cq)
		util.RespondSingle(rcq, c, prefs, err)
	})

	prefsRouter.PUT("", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		var update MemoryPreferences
		if err := c.ShouldBind(&update); err != nil {
			util.RespondBadRequest(rcq, c, "Invalid preference data")
			return
		}

		err := UpdateMemoryPreferences(rcq, &update)
		util.RespondSingle(rcq, c, &update, err)
	})

	prefsRouter.GET("/is-valid", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		var messages []string

		prefs, err := GetMemoryPreferences(cq)
		if err != nil {
			util.RespondInternalError(rcq, c, err)
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

		util.RespondSingle(rcq, c, &messages, nil)
	})
}
