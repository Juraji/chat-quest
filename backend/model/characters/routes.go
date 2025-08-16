package characters

import (
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/core/util"
)

func Routes(router *gin.RouterGroup) {
	charactersRoutes(router)
	tagsRoutes(router)
}

func charactersRoutes(router *gin.RouterGroup) {
	charactersRouter := router.Group("/characters")

	charactersRouter.GET("", func(c *gin.Context) {
		characters, err := AllCharacterListViews()
		util.RespondList(c, characters, err)
	})

	charactersRouter.GET("/:characterId", func(c *gin.Context) {
		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid character ID")
			return
		}

		character, err := CharacterById(characterId)
		util.RespondSingle(c, character, err)
	})

	charactersRouter.POST("", func(c *gin.Context) {
		var newCharacter Character
		if err := c.ShouldBind(&newCharacter); err != nil {
			util.RespondBadRequest(c, "Invalid character data")
			return
		}

		err := CreateCharacter(&newCharacter)
		util.RespondSingle(c, &newCharacter, err)
	})

	charactersRouter.PUT("/:characterId", func(c *gin.Context) {
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

		err = UpdateCharacter(characterId, &character)
		util.RespondSingle(c, &character, err)
	})

	charactersRouter.DELETE("/:characterId", func(c *gin.Context) {
		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid character ID")
			return
		}

		err = DeleteCharacterById(characterId)
		util.RespondDeleted(c, err)
	})

	charactersRouter.GET("/:characterId/tags", func(c *gin.Context) {
		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid character ID")
			return
		}

		tags, err := TagsByCharacterId(characterId)
		util.RespondList(c, tags, err)
	})

	charactersRouter.POST("/:characterId/tags", func(c *gin.Context) {
		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid character ID")
			return
		}

		var tagIds []int
		if err := c.ShouldBind(&tagIds); err != nil {
			util.RespondBadRequest(c, "Invalid dialogue examples data")
			return
		}

		err = SetCharacterTags(characterId, tagIds)
		util.RespondEmpty(c, err)
	})

	charactersRouter.POST("/:characterId/tags/:tagId", func(c *gin.Context) {
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

		err = AddCharacterTag(characterId, tagId)
		util.RespondEmpty(c, err)
	})

	charactersRouter.DELETE("/:characterId/tags/:tagId", func(c *gin.Context) {
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

		err = RemoveCharacterTag(characterId, tagId)
		util.RespondDeleted(c, err)
	})

	charactersRouter.GET("/:characterId/dialogue-examples", func(c *gin.Context) {
		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid character ID")
			return
		}

		examples, err := DialogueExamplesByCharacterId(characterId)
		util.RespondList(c, examples, err)
	})

	charactersRouter.POST("/:characterId/dialogue-examples", func(c *gin.Context) {
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

		err = SetDialogueExamplesByCharacterId(characterId, examples)
		util.RespondEmpty(c, err)
	})

	charactersRouter.GET("/:characterId/greetings", func(c *gin.Context) {
		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid character ID")
			return
		}

		greetings, err := CharacterGreetingsByCharacterId(characterId)
		util.RespondList(c, greetings, err)
	})

	charactersRouter.POST("/:characterId/greetings", func(c *gin.Context) {
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

		err = SetGreetingsByCharacterId(characterId, greetings)
		util.RespondEmpty(c, err)
	})

	charactersRouter.GET("/:characterId/group-greetings", func(c *gin.Context) {
		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid character ID")
			return
		}

		greetings, err := CharacterGroupGreetingsByCharacterId(characterId)
		util.RespondList(c, greetings, err)
	})

	charactersRouter.POST("/:characterId/group-greetings", func(c *gin.Context) {
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

		err = SetGroupGreetingsByCharacterId(characterId, greetings)
		util.RespondEmpty(c, err)
	})
}

func tagsRoutes(router *gin.RouterGroup) {
	tagsRouter := router.Group("/tags")

	tagsRouter.GET("", func(c *gin.Context) {
		tags, err := AllTags()
		util.RespondList(c, tags, err)
	})

	tagsRouter.GET("/:id", func(c *gin.Context) {
		id, err := util.GetIDParam(c, "id")
		if err != nil {
			util.RespondBadRequest(c, "Invalid tag ID")
			return
		}

		tag, err := TagById(id)
		util.RespondSingle(c, tag, err)
	})

	tagsRouter.POST("", func(c *gin.Context) {
		var newTag Tag
		if err := c.ShouldBind(&newTag); err != nil {
			util.RespondBadRequest(c, "Invalid tag data")
			return
		}

		err := CreateTag(&newTag)
		util.RespondSingle(c, &newTag, err)
	})

	tagsRouter.PUT("/:id", func(c *gin.Context) {
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

		err = UpdateTag(id, &tag)
		util.RespondSingle(c, &tag, err)
	})

	tagsRouter.DELETE("/:id", func(c *gin.Context) {
		id, err := util.GetIDParam(c, "id")
		if err != nil {
			util.RespondBadRequest(c, "Invalid tag ID")
			return
		}

		err = DeleteTagById(id)
		util.RespondDeleted(c, err)
	})
}
