package api

import (
	"github.com/gin-gonic/gin"
	ch "juraji.nl/chat-quest/model/characters"
)

func CharactersRoutes(router *gin.RouterGroup) {
	charactersRouter := router.Group("/characters")

	charactersRouter.GET("", func(c *gin.Context) {
		characters, err := ch.AllCharacters()
		respondList(c, characters, err)
	})

	charactersRouter.GET("/:characterId", func(c *gin.Context) {
		characterId, ok := getParamAsID(c, "characterId")
		if !ok {
			respondBadRequest(c, "Invalid character ID", nil)
			return
		}

		character, err := ch.CharacterById(characterId)
		respondSingle(c, character, err)
	})

	charactersRouter.POST("", func(c *gin.Context) {
		var newCharacter ch.Character
		if err := c.ShouldBind(&newCharacter); err != nil {
			respondBadRequest(c, "Invalid character data", nil)
			return
		}

		err := ch.CreateCharacter(&newCharacter)
		respondSingle(c, &newCharacter, err)
	})

	charactersRouter.PUT("/:characterId", func(c *gin.Context) {
		characterId, ok := getParamAsID(c, "characterId")
		if !ok {
			respondBadRequest(c, "Invalid character ID", nil)
			return
		}

		var character ch.Character
		if err := c.ShouldBind(&character); err != nil {
			respondBadRequest(c, "Invalid character data", nil)
			return
		}

		err := ch.UpdateCharacter(characterId, &character)
		respondSingle(c, &character, err)
	})

	charactersRouter.DELETE("/:characterId", func(c *gin.Context) {
		characterId, ok := getParamAsID(c, "characterId")
		if !ok {
			respondBadRequest(c, "Invalid character ID", nil)
			return
		}

		err := ch.DeleteCharacterById(characterId)
		respondEmpty(c, err)
	})

	charactersRouter.GET("/:characterId/dialogue-examples", func(c *gin.Context) {
		characterId, ok := getParamAsID(c, "characterId")
		if !ok {
			respondBadRequest(c, "Invalid character ID", nil)
			return
		}

		examples, err := ch.DialogueExamplesByCharacterId(characterId)
		respondList(c, examples, err)
	})

	charactersRouter.POST("/:characterId/dialogue-examples", func(c *gin.Context) {
		characterId, ok := getParamAsID(c, "characterId")
		if !ok {
			respondBadRequest(c, "Invalid character ID", nil)
			return
		}

		var examples []string
		if err := c.ShouldBind(&examples); err != nil {
			respondBadRequest(c, "Invalid dialogue examples data", nil)
			return
		}

		err := ch.SetDialogueExamplesByCharacterId(characterId, examples)
		respondEmpty(c, err)
	})

	charactersRouter.GET("/:characterId/greetings", func(c *gin.Context) {
		characterId, ok := getParamAsID(c, "characterId")
		if !ok {
			respondBadRequest(c, "Invalid character ID", nil)
			return
		}

		greetings, err := ch.CharacterGreetingsByCharacterId(characterId)
		respondList(c, greetings, err)
	})

	charactersRouter.POST("/:characterId/greetings", func(c *gin.Context) {
		characterId, ok := getParamAsID(c, "characterId")
		if !ok {
			respondBadRequest(c, "Invalid character ID", nil)
			return
		}

		var greetings []string
		if err := c.ShouldBind(&greetings); err != nil {
			respondBadRequest(c, "Invalid greetings data", nil)
			return
		}

		err := ch.SetGreetingsByCharacterId(characterId, greetings)
		respondEmpty(c, err)
	})

	charactersRouter.POST("/:characterId/duplicate", func(c *gin.Context) {
		characterId, ok := getParamAsID(c, "characterId")
		if !ok {
			respondBadRequest(c, "Invalid character ID", nil)
			return
		}

		newCharacter, err := ch.DuplicateCharacter(characterId)
		respondSingle(c, newCharacter, err)
	})
}
