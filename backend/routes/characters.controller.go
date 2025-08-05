package routes

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/model"
	"juraji.nl/chat-quest/util"
)

func CharactersController(router *gin.RouterGroup, db *sql.DB) {
	charactersRouter := router.Group("/characters")

	charactersRouter.GET("", func(c *gin.Context) {
		characters, err := model.AllCharacters(db)
		util.RespondList(c, characters, err)
	})

	charactersRouter.GET("/with-tags", func(c *gin.Context) {
		characters, err := model.AllCharactersWithTags(db)
		util.RespondList(c, characters, err)
	})

	charactersRouter.GET("/:characterId", func(c *gin.Context) {
		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid character ID")
			return
		}

		character, err := model.CharacterById(db, characterId)
		util.RespondSingle(c, character, err)
	})

	charactersRouter.POST("", func(c *gin.Context) {
		var newCharacter model.Character
		if err := c.ShouldBind(&newCharacter); err != nil {
			util.RespondBadRequest(c, "Invalid character data")
			return
		}

		err := model.CreateCharacter(db, &newCharacter)
		util.RespondSingle(c, &newCharacter, err)
	})

	charactersRouter.PUT("/:characterId", func(c *gin.Context) {
		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid character ID")
			return
		}

		var character model.Character
		if err := c.ShouldBind(&character); err != nil {
			util.RespondBadRequest(c, "Invalid character data")
			return
		}

		err = model.UpdateCharacter(db, characterId, &character)
		util.RespondSingle(c, &character, err)
	})

	charactersRouter.DELETE("/:characterId", func(c *gin.Context) {
		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid character ID")
			return
		}

		err = model.DeleteCharacterById(db, characterId)
		util.RespondDeleted(c, err)
	})

	charactersRouter.GET("/:characterId/details", func(c *gin.Context) {
		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid character ID")
			return
		}

		details, err := model.CharacterDetailsByCharacterId(db, characterId)
		util.RespondSingle(c, details, err)
	})

	charactersRouter.PUT("/:characterId/details", func(c *gin.Context) {
		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid character ID")
			return
		}

		var details model.CharacterDetails
		if err := c.ShouldBind(&details); err != nil {
			util.RespondBadRequest(c, "Invalid character details data")
			return
		}

		err = model.UpdateCharacterDetails(db, characterId, &details)
		util.RespondSingle(c, &details, err)
	})

	charactersRouter.GET("/:characterId/tags", func(c *gin.Context) {
		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid character ID")
			return
		}

		tags, err := model.TagsByCharacterId(db, characterId)
		util.RespondList(c, tags, err)
	})

	charactersRouter.POST("/:characterId/tags", func(c *gin.Context) {
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

		err = model.SetCharacterTags(db, characterId, tagIds)
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

		err = model.AddCharacterTag(db, characterId, tagId)
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

		err = model.RemoveCharacterTag(db, characterId, tagId)
		util.RespondDeleted(c, err)
	})

	charactersRouter.GET("/:characterId/dialogue-examples", func(c *gin.Context) {
		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid character ID")
			return
		}

		examples, err := model.DialogueExamplesByCharacterId(db, characterId)
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

		err = model.SetDialogueExamplesByCharacterId(db, characterId, examples)
		util.RespondEmpty(c, err)
	})

	charactersRouter.GET("/:characterId/greetings", func(c *gin.Context) {
		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid character ID")
			return
		}

		greetings, err := model.CharacterGreetingsByCharacterId(db, characterId)
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

		err = model.SetGreetingsByCharacterId(db, characterId, greetings)
		util.RespondEmpty(c, err)
	})

	charactersRouter.GET("/:characterId/group-greetings", func(c *gin.Context) {
		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid character ID")
			return
		}

		greetings, err := model.CharacterGroupGreetingsByCharacterId(db, characterId)
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

		err = model.SetGroupGreetingsByCharacterId(db, characterId, greetings)
		util.RespondEmpty(c, err)
	})
}
