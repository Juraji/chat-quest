package characters

import (
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/core"
	"juraji.nl/chat-quest/core/util"
)

func Routes(cq *core.ChatQuestContext, router *gin.RouterGroup) {
	charactersRoutes(cq, router)
	tagsRoutes(cq, router)
}

func charactersRoutes(cq *core.ChatQuestContext, router *gin.RouterGroup) {
	charactersRouter := router.Group("/characters")

	charactersRouter.GET("", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		characters, err := AllCharacters(cq)
		util.RespondList(rcq, c, characters, err)
	})

	charactersRouter.GET("/with-tags", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		characters, err := AllCharactersWithTags(cq)
		util.RespondList(rcq, c, characters, err)
	})

	charactersRouter.GET("/:characterId", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid character ID")
			return
		}

		character, err := CharacterById(rcq, characterId)
		util.RespondSingle(rcq, c, character, err)
	})

	charactersRouter.POST("", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		var newCharacter Character
		if err := c.ShouldBind(&newCharacter); err != nil {
			util.RespondBadRequest(rcq, c, "Invalid character data")
			return
		}

		err := CreateCharacter(rcq, &newCharacter)
		util.RespondSingle(rcq, c, &newCharacter, err)
	})

	charactersRouter.PUT("/:characterId", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid character ID")
			return
		}

		var character Character
		if err := c.ShouldBind(&character); err != nil {
			util.RespondBadRequest(rcq, c, "Invalid character data")
			return
		}

		err = UpdateCharacter(rcq, characterId, &character)
		util.RespondSingle(rcq, c, &character, err)
	})

	charactersRouter.DELETE("/:characterId", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid character ID")
			return
		}

		err = DeleteCharacterById(rcq, characterId)
		util.RespondDeleted(rcq, c, err)
	})

	charactersRouter.GET("/:characterId/details", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid character ID")
			return
		}

		details, err := CharacterDetailsByCharacterId(rcq, characterId)
		util.RespondSingle(rcq, c, details, err)
	})

	charactersRouter.PUT("/:characterId/details", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid character ID")
			return
		}

		var details CharacterDetails
		if err := c.ShouldBind(&details); err != nil {
			util.RespondBadRequest(rcq, c, "Invalid character details data")
			return
		}

		err = UpdateCharacterDetails(rcq, characterId, &details)
		util.RespondSingle(rcq, c, &details, err)
	})

	charactersRouter.GET("/:characterId/tags", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid character ID")
			return
		}

		tags, err := TagsByCharacterId(rcq, characterId)
		util.RespondList(rcq, c, tags, err)
	})

	charactersRouter.POST("/:characterId/tags", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid character ID")
			return
		}

		var tagIds []int
		if err := c.ShouldBind(&tagIds); err != nil {
			util.RespondBadRequest(rcq, c, "Invalid dialogue examples data")
			return
		}

		err = SetCharacterTags(rcq, characterId, tagIds)
		util.RespondEmpty(rcq, c, err)
	})

	charactersRouter.POST("/:characterId/tags/:tagId", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid character ID")
			return
		}
		tagId, err := util.GetIDParam(c, "tagId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid tag ID")
			return
		}

		err = AddCharacterTag(rcq, characterId, tagId)
		util.RespondEmpty(rcq, c, err)
	})

	charactersRouter.DELETE("/:characterId/tags/:tagId", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid character ID")
			return
		}
		tagId, err := util.GetIDParam(c, "tagId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid tag ID")
			return
		}

		err = RemoveCharacterTag(rcq, characterId, tagId)
		util.RespondDeleted(rcq, c, err)
	})

	charactersRouter.GET("/:characterId/dialogue-examples", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid character ID")
			return
		}

		examples, err := DialogueExamplesByCharacterId(rcq, characterId)
		util.RespondList(rcq, c, examples, err)
	})

	charactersRouter.POST("/:characterId/dialogue-examples", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid character ID")
			return
		}

		var examples []string
		if err := c.ShouldBind(&examples); err != nil {
			util.RespondBadRequest(rcq, c, "Invalid dialogue examples data")
			return
		}

		err = SetDialogueExamplesByCharacterId(rcq, characterId, examples)
		util.RespondEmpty(rcq, c, err)
	})

	charactersRouter.GET("/:characterId/greetings", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid character ID")
			return
		}

		greetings, err := CharacterGreetingsByCharacterId(rcq, characterId)
		util.RespondList(rcq, c, greetings, err)
	})

	charactersRouter.POST("/:characterId/greetings", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid character ID")
			return
		}

		var greetings []string
		if err := c.ShouldBind(&greetings); err != nil {
			util.RespondBadRequest(rcq, c, "Invalid greetings data")
			return
		}

		err = SetGreetingsByCharacterId(rcq, characterId, greetings)
		util.RespondEmpty(rcq, c, err)
	})

	charactersRouter.GET("/:characterId/group-greetings", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid character ID")
			return
		}

		greetings, err := CharacterGroupGreetingsByCharacterId(rcq, characterId)
		util.RespondList(rcq, c, greetings, err)
	})

	charactersRouter.POST("/:characterId/group-greetings", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid character ID")
			return
		}

		var greetings []string
		if err := c.ShouldBind(&greetings); err != nil {
			util.RespondBadRequest(rcq, c, "Invalid greetings data")
			return
		}

		err = SetGroupGreetingsByCharacterId(rcq, characterId, greetings)
		util.RespondEmpty(rcq, c, err)
	})
}

func tagsRoutes(cq *core.ChatQuestContext, router *gin.RouterGroup) {
	tagsRouter := router.Group("/tags")

	tagsRouter.GET("", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		tags, err := AllTags(cq)
		util.RespondList(rcq, c, tags, err)
	})

	tagsRouter.GET("/:id", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		id, err := util.GetIDParam(c, "id")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid tag ID")
			return
		}

		tag, err := TagById(rcq, id)
		util.RespondSingle(rcq, c, tag, err)
	})

	tagsRouter.POST("", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		var newTag Tag
		if err := c.ShouldBind(&newTag); err != nil {
			util.RespondBadRequest(rcq, c, "Invalid tag data")
			return
		}

		err := CreateTag(rcq, &newTag)
		util.RespondSingle(rcq, c, &newTag, err)
	})

	tagsRouter.PUT("/:id", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		id, err := util.GetIDParam(c, "id")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid tag ID")
			return
		}

		var tag Tag
		if err := c.ShouldBind(&tag); err != nil {
			util.RespondBadRequest(rcq, c, "Invalid tag data")
			return
		}

		err = UpdateTag(rcq, id, &tag)
		util.RespondSingle(rcq, c, &tag, err)
	})

	tagsRouter.DELETE("/:id", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		id, err := util.GetIDParam(c, "id")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid tag ID")
			return
		}

		err = DeleteTagById(rcq, id)
		util.RespondDeleted(rcq, c, err)
	})
}
