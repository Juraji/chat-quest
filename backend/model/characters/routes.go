package characters

import (
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/core/util/controllers"
)

func Routes(router *gin.RouterGroup) {
	charactersRoutes(router)
	tagsRoutes(router)
}

func charactersRoutes(router *gin.RouterGroup) {
	charactersRouter := router.Group("/characters")

	charactersRouter.GET("", func(c *gin.Context) {
		characters, err := AllCharacterListViews()
		controllers.RespondList(c, characters, err)
	})

	charactersRouter.GET("/:characterId", func(c *gin.Context) {
		characterId, ok := controllers.GetParamAsID(c, "characterId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid character ID", nil)
			return
		}

		character, err := CharacterById(characterId)
		controllers.RespondSingle(c, character, err)
	})

	charactersRouter.POST("", func(c *gin.Context) {
		var newCharacter Character
		if err := c.ShouldBind(&newCharacter); err != nil {
			controllers.RespondBadRequest(c, "Invalid character data", nil)
			return
		}

		err := CreateCharacter(&newCharacter)
		controllers.RespondSingle(c, &newCharacter, err)
	})

	charactersRouter.PUT("/:characterId", func(c *gin.Context) {
		characterId, ok := controllers.GetParamAsID(c, "characterId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid character ID", nil)
			return
		}

		var character Character
		if err := c.ShouldBind(&character); err != nil {
			controllers.RespondBadRequest(c, "Invalid character data", nil)
			return
		}

		err := UpdateCharacter(characterId, &character)
		controllers.RespondSingle(c, &character, err)
	})

	charactersRouter.DELETE("/:characterId", func(c *gin.Context) {
		characterId, ok := controllers.GetParamAsID(c, "characterId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid character ID", nil)
			return
		}

		err := DeleteCharacterById(characterId)
		controllers.RespondEmpty(c, err)
	})

	charactersRouter.GET("/:characterId/tags", func(c *gin.Context) {
		characterId, ok := controllers.GetParamAsID(c, "characterId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid character ID", nil)
			return
		}

		tags, err := TagsByCharacterId(characterId)
		controllers.RespondList(c, tags, err)
	})

	charactersRouter.POST("/:characterId/tags", func(c *gin.Context) {
		characterId, ok := controllers.GetParamAsID(c, "characterId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid character ID", nil)
			return
		}

		var tagIds []int
		if err := c.ShouldBind(&tagIds); err != nil {
			controllers.RespondBadRequest(c, "Invalid dialogue examples data", nil)
			return
		}

		err := SetCharacterTags(characterId, tagIds)
		controllers.RespondEmpty(c, err)
	})

	charactersRouter.POST("/:characterId/tags/:tagId", func(c *gin.Context) {
		characterId, ok := controllers.GetParamAsID(c, "characterId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid character ID", nil)
			return
		}
		tagId, ok := controllers.GetParamAsID(c, "tagId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid tag ID", nil)
			return
		}

		err := AddCharacterTag(characterId, tagId)
		controllers.RespondEmpty(c, err)
	})

	charactersRouter.DELETE("/:characterId/tags/:tagId", func(c *gin.Context) {
		characterId, ok := controllers.GetParamAsID(c, "characterId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid character ID", nil)
			return
		}
		tagId, ok := controllers.GetParamAsID(c, "tagId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid tag ID", nil)
			return
		}

		err := RemoveCharacterTag(characterId, tagId)
		controllers.RespondEmpty(c, err)
	})

	charactersRouter.GET("/:characterId/dialogue-examples", func(c *gin.Context) {
		characterId, ok := controllers.GetParamAsID(c, "characterId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid character ID", nil)
			return
		}

		examples, err := DialogueExamplesByCharacterId(characterId)
		controllers.RespondList(c, examples, err)
	})

	charactersRouter.POST("/:characterId/dialogue-examples", func(c *gin.Context) {
		characterId, ok := controllers.GetParamAsID(c, "characterId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid character ID", nil)
			return
		}

		var examples []string
		if err := c.ShouldBind(&examples); err != nil {
			controllers.RespondBadRequest(c, "Invalid dialogue examples data", nil)
			return
		}

		err := SetDialogueExamplesByCharacterId(characterId, examples)
		controllers.RespondEmpty(c, err)
	})

	charactersRouter.GET("/:characterId/greetings", func(c *gin.Context) {
		characterId, ok := controllers.GetParamAsID(c, "characterId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid character ID", nil)
			return
		}

		greetings, err := CharacterGreetingsByCharacterId(characterId)
		controllers.RespondList(c, greetings, err)
	})

	charactersRouter.POST("/:characterId/greetings", func(c *gin.Context) {
		characterId, ok := controllers.GetParamAsID(c, "characterId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid character ID", nil)
			return
		}

		var greetings []string
		if err := c.ShouldBind(&greetings); err != nil {
			controllers.RespondBadRequest(c, "Invalid greetings data", nil)
			return
		}

		err := SetGreetingsByCharacterId(characterId, greetings)
		controllers.RespondEmpty(c, err)
	})

	charactersRouter.GET("/:characterId/group-greetings", func(c *gin.Context) {
		characterId, ok := controllers.GetParamAsID(c, "characterId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid character ID", nil)
			return
		}

		greetings, err := CharacterGroupGreetingsByCharacterId(characterId)
		controllers.RespondList(c, greetings, err)
	})

	charactersRouter.POST("/:characterId/group-greetings", func(c *gin.Context) {
		characterId, ok := controllers.GetParamAsID(c, "characterId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid character ID", nil)
			return
		}

		var greetings []string
		if err := c.ShouldBind(&greetings); err != nil {
			controllers.RespondBadRequest(c, "Invalid greetings data", nil)
			return
		}

		err := SetGroupGreetingsByCharacterId(characterId, greetings)
		controllers.RespondEmpty(c, err)
	})
}

func tagsRoutes(router *gin.RouterGroup) {
	tagsRouter := router.Group("/tags")

	tagsRouter.GET("", func(c *gin.Context) {
		tags, err := AllTags()
		controllers.RespondList(c, tags, err)
	})

	tagsRouter.GET("/:id", func(c *gin.Context) {
		id, ok := controllers.GetParamAsID(c, "id")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid tag ID", nil)
			return
		}

		tag, err := TagById(id)
		controllers.RespondSingle(c, &tag, err)
	})

	tagsRouter.POST("", func(c *gin.Context) {
		var newTag Tag
		if err := c.ShouldBind(&newTag); err != nil {
			controllers.RespondBadRequest(c, "Invalid tag data", nil)
			return
		}

		err := CreateTag(&newTag)
		controllers.RespondSingle(c, &newTag, err)
	})

	tagsRouter.PUT("/:id", func(c *gin.Context) {
		id, ok := controllers.GetParamAsID(c, "id")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid tag ID", nil)
			return
		}

		var tag Tag
		if err := c.ShouldBind(&tag); err != nil {
			controllers.RespondBadRequest(c, "Invalid tag data", nil)
			return
		}

		err := UpdateTag(id, &tag)
		controllers.RespondSingle(c, &tag, err)
	})

	tagsRouter.DELETE("/:id", func(c *gin.Context) {
		id, ok := controllers.GetParamAsID(c, "id")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid tag ID", nil)
			return
		}

		err := DeleteTagById(id)
		controllers.RespondEmpty(c, err)
	})
}
