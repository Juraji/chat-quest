package characters

import (
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/cq"
	"juraji.nl/chat-quest/util"
)

func Routes(cq *cq.ChatQuestContext, router *gin.RouterGroup) {
	charactersRoutes(cq, router)
	tagsRoutes(cq, router)
}

func charactersRoutes(cq *cq.ChatQuestContext, router *gin.RouterGroup) {
	charactersRouter := router.Group("/characters")

	charactersRouter.GET("", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		characters, err := AllCharacters(cq)
		util.RespondList(c, characters, err)
	})

	charactersRouter.GET("/with-tags", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		characters, err := AllCharactersWithTags(cq)
		util.RespondList(c, characters, err)
	})

	charactersRouter.GET("/:characterId", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid character ID")
			return
		}

		character, err := CharacterById(cq, characterId)
		util.RespondSingle(c, character, err)
	})

	charactersRouter.POST("", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		var newCharacter Character
		if err := c.ShouldBind(&newCharacter); err != nil {
			util.RespondBadRequest(c, "Invalid character data")
			return
		}

		err := CreateCharacter(cq, &newCharacter)
		util.RespondSingle(c, &newCharacter, err)
	})

	charactersRouter.PUT("/:characterId", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid character ID")
			return
		}

		var character Character
		if err := c.ShouldBind(&character); err != nil {
			util.RespondBadRequest(c, "Invalid character data")
			return
		}

		err = UpdateCharacter(cq, characterId, &character)
		util.RespondSingle(c, &character, err)
	})

	charactersRouter.DELETE("/:characterId", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid character ID")
			return
		}

		err = DeleteCharacterById(cq, characterId)
		util.RespondDeleted(c, err)
	})

	charactersRouter.GET("/:characterId/details", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid character ID")
			return
		}

		details, err := CharacterDetailsByCharacterId(cq, characterId)
		util.RespondSingle(c, details, err)
	})

	charactersRouter.PUT("/:characterId/details", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid character ID")
			return
		}

		var details CharacterDetails
		if err := c.ShouldBind(&details); err != nil {
			util.RespondBadRequest(c, "Invalid character details data")
			return
		}

		err = UpdateCharacterDetails(cq, characterId, &details)
		util.RespondSingle(c, &details, err)
	})

	charactersRouter.GET("/:characterId/tags", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid character ID")
			return
		}

		tags, err := TagsByCharacterId(cq, characterId)
		util.RespondList(c, tags, err)
	})

	charactersRouter.POST("/:characterId/tags", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid character ID")
			return
		}

		var tagIds []int64
		if err := c.ShouldBind(&tagIds); err != nil {
			util.RespondBadRequest(c, "Invalid dialogue examples data")
			return
		}

		err = SetCharacterTags(cq, characterId, tagIds)
		util.RespondEmpty(c, err)
	})

	charactersRouter.POST("/:characterId/tags/:tagId", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid character ID")
			return
		}
		tagId, err := util.GetIDParam(c, "tagId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid tag ID")
			return
		}

		err = AddCharacterTag(cq, characterId, tagId)
		util.RespondEmpty(c, err)
	})

	charactersRouter.DELETE("/:characterId/tags/:tagId", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid character ID")
			return
		}
		tagId, err := util.GetIDParam(c, "tagId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid tag ID")
			return
		}

		err = RemoveCharacterTag(cq, characterId, tagId)
		util.RespondDeleted(c, err)
	})

	charactersRouter.GET("/:characterId/dialogue-examples", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid character ID")
			return
		}

		examples, err := DialogueExamplesByCharacterId(cq, characterId)
		util.RespondList(c, examples, err)
	})

	charactersRouter.POST("/:characterId/dialogue-examples", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid character ID")
			return
		}

		var examples []string
		if err := c.ShouldBind(&examples); err != nil {
			util.RespondBadRequest(c, "Invalid dialogue examples data")
			return
		}

		err = SetDialogueExamplesByCharacterId(cq, characterId, examples)
		util.RespondEmpty(c, err)
	})

	charactersRouter.GET("/:characterId/greetings", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid character ID")
			return
		}

		greetings, err := CharacterGreetingsByCharacterId(cq, characterId)
		util.RespondList(c, greetings, err)
	})

	charactersRouter.POST("/:characterId/greetings", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid character ID")
			return
		}

		var greetings []string
		if err := c.ShouldBind(&greetings); err != nil {
			util.RespondBadRequest(c, "Invalid greetings data")
			return
		}

		err = SetGreetingsByCharacterId(cq, characterId, greetings)
		util.RespondEmpty(c, err)
	})

	charactersRouter.GET("/:characterId/group-greetings", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid character ID")
			return
		}

		greetings, err := CharacterGroupGreetingsByCharacterId(cq, characterId)
		util.RespondList(c, greetings, err)
	})

	charactersRouter.POST("/:characterId/group-greetings", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid character ID")
			return
		}

		var greetings []string
		if err := c.ShouldBind(&greetings); err != nil {
			util.RespondBadRequest(c, "Invalid greetings data")
			return
		}

		err = SetGroupGreetingsByCharacterId(cq, characterId, greetings)
		util.RespondEmpty(c, err)
	})
}

func tagsRoutes(cq *cq.ChatQuestContext, router *gin.RouterGroup) {
	tagsRouter := router.Group("/tags")

	tagsRouter.GET("", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		tags, err := AllTags(cq)
		util.RespondList(c, tags, err)
	})

	tagsRouter.GET("/:id", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		id, err := util.GetIDParam(c, "id")
		if err != nil {
			util.RespondBadRequest(c, "Invalid tag ID")
			return
		}

		tag, err := TagById(cq, id)
		util.RespondSingle(c, tag, err)
	})

	tagsRouter.POST("", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		var newTag Tag
		if err := c.ShouldBind(&newTag); err != nil {
			util.RespondBadRequest(c, "Invalid tag data")
			return
		}

		err := CreateTag(cq, &newTag)
		util.RespondSingle(c, &newTag, err)
	})

	tagsRouter.PUT("/:id", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		id, err := util.GetIDParam(c, "id")
		if err != nil {
			util.RespondBadRequest(c, "Invalid tag ID")
			return
		}

		var tag Tag
		if err := c.ShouldBind(&tag); err != nil {
			util.RespondBadRequest(c, "Invalid tag data")
			return
		}

		err = UpdateTag(cq, id, &tag)
		util.RespondSingle(c, &tag, err)
	})

	tagsRouter.DELETE("/:id", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		id, err := util.GetIDParam(c, "id")
		if err != nil {
			util.RespondBadRequest(c, "Invalid tag ID")
			return
		}

		err = DeleteTagById(cq, id)
		util.RespondDeleted(c, err)
	})
}
