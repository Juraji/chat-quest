package characters

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/util"
)

func Routes(router *gin.RouterGroup, db *sql.DB) {
	charactersRouter := router.Group("/characters")
	tagsRouter := router.Group("/tags")

	charactersRouter.GET("", func(c *gin.Context) {
		characters, err := AllCharacters(db)
		util.RespondList(c, characters, err)
	})

	charactersRouter.GET("/with-tags", func(c *gin.Context) {
		characters, err := AllCharactersWithTags(db)
		util.RespondList(c, characters, err)
	})

	charactersRouter.GET("/:characterId", func(c *gin.Context) {
		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid character ID")
			return
		}

		character, err := CharacterById(db, characterId)
		util.RespondSingle(c, character, err)
	})

	charactersRouter.POST("", func(c *gin.Context) {
		var newCharacter Character
		if err := c.ShouldBind(&newCharacter); err != nil {
			util.RespondBadRequest(c, "Invalid character data")
			return
		}

		err := CreateCharacter(db, &newCharacter)
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

		err = UpdateCharacter(db, characterId, &character)
		util.RespondSingle(c, &character, err)
	})

	charactersRouter.DELETE("/:characterId", func(c *gin.Context) {
		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid character ID")
			return
		}

		err = DeleteCharacterById(db, characterId)
		util.RespondDeleted(c, err)
	})

	charactersRouter.GET("/:characterId/details", func(c *gin.Context) {
		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid character ID")
			return
		}

		details, err := CharacterDetailsByCharacterId(db, characterId)
		util.RespondSingle(c, details, err)
	})

	charactersRouter.PUT("/:characterId/details", func(c *gin.Context) {
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

		err = UpdateCharacterDetails(db, characterId, &details)
		util.RespondSingle(c, &details, err)
	})

	charactersRouter.GET("/:characterId/tags", func(c *gin.Context) {
		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid character ID")
			return
		}

		tags, err := TagsByCharacterId(db, characterId)
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

		err = SetCharacterTags(db, characterId, tagIds)
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

		err = AddCharacterTag(db, characterId, tagId)
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

		err = RemoveCharacterTag(db, characterId, tagId)
		util.RespondDeleted(c, err)
	})

	charactersRouter.GET("/:characterId/dialogue-examples", func(c *gin.Context) {
		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid character ID")
			return
		}

		examples, err := DialogueExamplesByCharacterId(db, characterId)
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

		err = SetDialogueExamplesByCharacterId(db, characterId, examples)
		util.RespondEmpty(c, err)
	})

	charactersRouter.GET("/:characterId/greetings", func(c *gin.Context) {
		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid character ID")
			return
		}

		greetings, err := CharacterGreetingsByCharacterId(db, characterId)
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

		err = SetGreetingsByCharacterId(db, characterId, greetings)
		util.RespondEmpty(c, err)
	})

	charactersRouter.GET("/:characterId/group-greetings", func(c *gin.Context) {
		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid character ID")
			return
		}

		greetings, err := CharacterGroupGreetingsByCharacterId(db, characterId)
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

		err = SetGroupGreetingsByCharacterId(db, characterId, greetings)
		util.RespondEmpty(c, err)
	})

	tagsRouter.GET("", func(c *gin.Context) {
		tags, err := AllTags(db)
		util.RespondList(c, tags, err)
	})

	tagsRouter.GET("/:id", func(c *gin.Context) {
		id, err := util.GetIDParam(c, "id")
		if err != nil {
			util.RespondBadRequest(c, "Invalid tag ID")
			return
		}

		tag, err := TagById(db, id)
		util.RespondSingle(c, tag, err)
	})

	tagsRouter.POST("", func(c *gin.Context) {
		var newTag Tag
		if err := c.ShouldBind(&newTag); err != nil {
			util.RespondBadRequest(c, "Invalid tag data")
			return
		}

		err := CreateTag(db, &newTag)
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

		err = UpdateTag(db, id, &tag)
		util.RespondSingle(c, &tag, err)
	})

	tagsRouter.DELETE("/:id", func(c *gin.Context) {
		id, err := util.GetIDParam(c, "id")
		if err != nil {
			util.RespondBadRequest(c, "Invalid tag ID")
			return
		}

		err = DeleteTagById(db, id)
		util.RespondDeleted(c, err)
	})
}
