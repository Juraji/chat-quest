package characters

import (
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/core/util/controllers"
)

func Routes(router *gin.RouterGroup) {
	charactersRouter := router.Group("/characters")

	charactersRouter.GET("", func(c *gin.Context) {
		characters, err := AllCharacters()
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

	charactersRouter.POST("/:characterId/duplicate", func(c *gin.Context) {
		characterId, ok := controllers.GetParamAsID(c, "characterId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid character ID", nil)
			return
		}

		newCharacter, err := DuplicateCharacter(characterId)
		controllers.RespondSingle(c, newCharacter, err)
	})
}
