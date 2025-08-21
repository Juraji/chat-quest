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
		characters, ok := AllCharacterListViews()
		controllers.RespondList(c, ok, characters)
	})

	charactersRouter.GET("/:characterId", func(c *gin.Context) {
		characterId, ok := controllers.GetParamAsID(c, "characterId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid character ID")
			return
		}

		character, ok := CharacterById(characterId)
		controllers.RespondSingle(c, ok, character)
	})

	charactersRouter.POST("", func(c *gin.Context) {
		var newCharacter Character
		if err := c.ShouldBind(&newCharacter); err != nil {
			controllers.RespondBadRequest(c, "Invalid character data")
			return
		}

		ok := CreateCharacter(&newCharacter)
		controllers.RespondSingle(c, ok, &newCharacter)
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

		ok = UpdateCharacter(characterId, &character)
		controllers.RespondSingle(c, ok, &character)
	})

	charactersRouter.DELETE("/:characterId", func(c *gin.Context) {
		characterId, ok := controllers.GetParamAsID(c, "characterId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid character ID")
			return
		}

		ok = DeleteCharacterById(characterId)
		controllers.RespondEmpty(c, ok)
	})

	charactersRouter.GET("/:characterId/tags", func(c *gin.Context) {
		characterId, ok := controllers.GetParamAsID(c, "characterId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid character ID")
			return
		}

		tags, ok := TagsByCharacterId(characterId)
		controllers.RespondList(c, ok, tags)
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

		ok = SetCharacterTags(characterId, tagIds)
		controllers.RespondEmpty(c, ok)
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

		ok = AddCharacterTag(characterId, tagId)
		controllers.RespondEmpty(c, ok)
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

		ok = RemoveCharacterTag(characterId, tagId)
		controllers.RespondEmpty(c, ok)
	})

	charactersRouter.GET("/:characterId/dialogue-examples", func(c *gin.Context) {
		characterId, ok := controllers.GetParamAsID(c, "characterId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid character ID")
			return
		}

		examples, ok := DialogueExamplesByCharacterId(characterId)
		controllers.RespondList(c, ok, examples)
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

		ok = SetDialogueExamplesByCharacterId(characterId, examples)
		controllers.RespondEmpty(c, ok)
	})

	charactersRouter.GET("/:characterId/greetings", func(c *gin.Context) {
		characterId, ok := controllers.GetParamAsID(c, "characterId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid character ID")
			return
		}

		greetings, ok := CharacterGreetingsByCharacterId(characterId)
		controllers.RespondList(c, ok, greetings)
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

		ok = SetGreetingsByCharacterId(characterId, greetings)
		controllers.RespondEmpty(c, ok)
	})

	charactersRouter.GET("/:characterId/group-greetings", func(c *gin.Context) {
		characterId, ok := controllers.GetParamAsID(c, "characterId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid character ID")
			return
		}

		greetings, ok := CharacterGroupGreetingsByCharacterId(characterId)
		controllers.RespondList(c, ok, greetings)
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

		ok = SetGroupGreetingsByCharacterId(characterId, greetings)
		controllers.RespondEmpty(c, ok)
	})
}

func tagsRoutes(router *gin.RouterGroup) {
	tagsRouter := router.Group("/tags")

	tagsRouter.GET("", func(c *gin.Context) {
		tags, ok := AllTags()
		controllers.RespondList(c, ok, tags)
	})

	tagsRouter.GET("/:id", func(c *gin.Context) {
		id, ok := controllers.GetParamAsID(c, "id")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid tag ID")
			return
		}

		tag, ok := TagById(id)
		controllers.RespondSingle(c, ok, &tag)
	})

	tagsRouter.POST("", func(c *gin.Context) {
		var newTag Tag
		if err := c.ShouldBind(&newTag); err != nil {
			controllers.RespondBadRequest(c, "Invalid tag data")
			return
		}

		ok := CreateTag(&newTag)
		controllers.RespondSingle(c, ok, &newTag)
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

		ok = UpdateTag(id, &tag)
		controllers.RespondSingle(c, ok, &tag)
	})

	tagsRouter.DELETE("/:id", func(c *gin.Context) {
		id, ok := controllers.GetParamAsID(c, "id")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid tag ID")
			return
		}

		ok = DeleteTagById(id)
		controllers.RespondEmpty(c, ok)
	})
}
