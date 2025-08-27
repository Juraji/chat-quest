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
		controllers.RespondListE(c, characters, err)
	})

	charactersRouter.GET("/:characterId", func(c *gin.Context) {
		characterId, ok := controllers.GetParamAsID(c, "characterId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid character ID")
			return
		}

		character, err := CharacterById(characterId)
		controllers.RespondSingleE(c, character, err)
	})

	charactersRouter.POST("", func(c *gin.Context) {
		var newCharacter Character
		if err := c.ShouldBind(&newCharacter); err != nil {
			controllers.RespondBadRequest(c, "Invalid character data")
			return
		}

		err := CreateCharacter(&newCharacter)
		controllers.RespondSingleE(c, &newCharacter, err)
	})

	charactersRouter.PUT("/:characterId", func(c *gin.Context) {
		characterId, ok := controllers.GetParamAsID(c, "characterId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid character ID")
			return
		}

		var character Character
		if err := c.ShouldBind(&character); err != nil {
			controllers.RespondBadRequest(c, "Invalid character data")
			return
		}

		err := UpdateCharacter(characterId, &character)
		controllers.RespondSingleE(c, &character, err)
	})

	charactersRouter.DELETE("/:characterId", func(c *gin.Context) {
		characterId, ok := controllers.GetParamAsID(c, "characterId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid character ID")
			return
		}

		err := DeleteCharacterById(characterId)
		controllers.RespondEmptyE(c, err)
	})

	charactersRouter.GET("/:characterId/tags", func(c *gin.Context) {
		characterId, ok := controllers.GetParamAsID(c, "characterId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid character ID")
			return
		}

		tags, err := TagsByCharacterId(characterId)
		controllers.RespondListE(c, tags, err)
	})

	charactersRouter.POST("/:characterId/tags", func(c *gin.Context) {
		characterId, ok := controllers.GetParamAsID(c, "characterId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid character ID")
			return
		}

		var tagIds []int
		if err := c.ShouldBind(&tagIds); err != nil {
			controllers.RespondBadRequest(c, "Invalid dialogue examples data")
			return
		}

		err := SetCharacterTags(characterId, tagIds)
		controllers.RespondEmptyE(c, err)
	})

	charactersRouter.POST("/:characterId/tags/:tagId", func(c *gin.Context) {
		characterId, ok := controllers.GetParamAsID(c, "characterId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid character ID")
			return
		}
		tagId, ok := controllers.GetParamAsID(c, "tagId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid tag ID")
			return
		}

		err := AddCharacterTag(characterId, tagId)
		controllers.RespondEmptyE(c, err)
	})

	charactersRouter.DELETE("/:characterId/tags/:tagId", func(c *gin.Context) {
		characterId, ok := controllers.GetParamAsID(c, "characterId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid character ID")
			return
		}
		tagId, ok := controllers.GetParamAsID(c, "tagId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid tag ID")
			return
		}

		err := RemoveCharacterTag(characterId, tagId)
		controllers.RespondEmptyE(c, err)
	})

	charactersRouter.GET("/:characterId/dialogue-examples", func(c *gin.Context) {
		characterId, ok := controllers.GetParamAsID(c, "characterId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid character ID")
			return
		}

		examples, err := DialogueExamplesByCharacterId(characterId)
		controllers.RespondListE(c, examples, err)
	})

	charactersRouter.POST("/:characterId/dialogue-examples", func(c *gin.Context) {
		characterId, ok := controllers.GetParamAsID(c, "characterId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid character ID")
			return
		}

		var examples []string
		if err := c.ShouldBind(&examples); err != nil {
			controllers.RespondBadRequest(c, "Invalid dialogue examples data")
			return
		}

		err := SetDialogueExamplesByCharacterId(characterId, examples)
		controllers.RespondEmptyE(c, err)
	})

	charactersRouter.GET("/:characterId/greetings", func(c *gin.Context) {
		characterId, ok := controllers.GetParamAsID(c, "characterId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid character ID")
			return
		}

		greetings, err := CharacterGreetingsByCharacterId(characterId)
		controllers.RespondListE(c, greetings, err)
	})

	charactersRouter.POST("/:characterId/greetings", func(c *gin.Context) {
		characterId, ok := controllers.GetParamAsID(c, "characterId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid character ID")
			return
		}

		var greetings []string
		if err := c.ShouldBind(&greetings); err != nil {
			controllers.RespondBadRequest(c, "Invalid greetings data")
			return
		}

		err := SetGreetingsByCharacterId(characterId, greetings)
		controllers.RespondEmptyE(c, err)
	})

	charactersRouter.GET("/:characterId/group-greetings", func(c *gin.Context) {
		characterId, ok := controllers.GetParamAsID(c, "characterId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid character ID")
			return
		}

		greetings, err := CharacterGroupGreetingsByCharacterId(characterId)
		controllers.RespondListE(c, greetings, err)
	})

	charactersRouter.POST("/:characterId/group-greetings", func(c *gin.Context) {
		characterId, ok := controllers.GetParamAsID(c, "characterId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid character ID")
			return
		}

		var greetings []string
		if err := c.ShouldBind(&greetings); err != nil {
			controllers.RespondBadRequest(c, "Invalid greetings data")
			return
		}

		err := SetGroupGreetingsByCharacterId(characterId, greetings)
		controllers.RespondEmptyE(c, err)
	})
}

func tagsRoutes(router *gin.RouterGroup) {
	tagsRouter := router.Group("/tags")

	tagsRouter.GET("", func(c *gin.Context) {
		tags, err := AllTags()
		controllers.RespondListE(c, tags, err)
	})

	tagsRouter.GET("/:id", func(c *gin.Context) {
		id, ok := controllers.GetParamAsID(c, "id")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid tag ID")
			return
		}

		tag, err := TagById(id)
		controllers.RespondSingleE(c, &tag, err)
	})

	tagsRouter.POST("", func(c *gin.Context) {
		var newTag Tag
		if err := c.ShouldBind(&newTag); err != nil {
			controllers.RespondBadRequest(c, "Invalid tag data")
			return
		}

		err := CreateTag(&newTag)
		controllers.RespondSingleE(c, &newTag, err)
	})

	tagsRouter.PUT("/:id", func(c *gin.Context) {
		id, ok := controllers.GetParamAsID(c, "id")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid tag ID")
			return
		}

		var tag Tag
		if err := c.ShouldBind(&tag); err != nil {
			controllers.RespondBadRequest(c, "Invalid tag data")
			return
		}

		err := UpdateTag(id, &tag)
		controllers.RespondSingleE(c, &tag, err)
	})

	tagsRouter.DELETE("/:id", func(c *gin.Context) {
		id, ok := controllers.GetParamAsID(c, "id")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid tag ID")
			return
		}

		err := DeleteTagById(id)
		controllers.RespondEmptyE(c, err)
	})
}
