package routes

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/model"
	"net/http"
)

func CharactersController(router *gin.RouterGroup, db *sql.DB) {
	charactersRouter := router.Group("/characters")

	charactersRouter.GET("/", func(c *gin.Context) {
		characters, err := model.AllCharacters(db)
		respondList(c, characters, err)
	})

	charactersRouter.GET("/:characterId", func(c *gin.Context) {
		characterId, err := getID(c, "characterId")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		character, err := model.CharacterById(db, characterId)
		respondSingle(c, character, err)
	})

	charactersRouter.POST("/", func(c *gin.Context) {
		var newCharacter model.Character
		if err := c.ShouldBind(&newCharacter); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character data"})
			return
		}

		err := model.CreateCharacter(db, &newCharacter)
		respondSingle(c, &newCharacter, err)
	})

	charactersRouter.PUT("/:characterId", func(c *gin.Context) {
		characterId, err := getID(c, "characterId")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		var character model.Character
		if err := c.ShouldBind(&character); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character data"})
			return
		}

		err = model.UpdateCharacter(db, characterId, &character)
		respondSingle(c, &character, err)
	})

	charactersRouter.DELETE("/:characterId", func(c *gin.Context) {
		characterId, err := getID(c, "characterId")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		err = model.DeleteCharacterById(db, characterId)
		respondDeleted(c, err)
	})

	charactersRouter.GET("/:characterId/details", func(c *gin.Context) {
		characterId, err := getID(c, "characterId")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		details, err := model.CharacterDetailsByCharacterId(db, characterId)
		respondSingle(c, details, err)
	})

	charactersRouter.PUT("/:characterId/details", func(c *gin.Context) {
		characterId, err := getID(c, "characterId")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		var details model.CharacterDetails
		if err := c.ShouldBind(&details); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character details data"})
			return
		}

		err = model.UpdateCharacterDetails(db, characterId, &details)
		respondSingle(c, &details, err)
	})

	charactersRouter.GET("/:characterId/tags", func(c *gin.Context) {
		characterId, err := getID(c, "characterId")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		tags, err := model.TagsByCharacterId(db, characterId)
		respondList(c, tags, err)
	})

	charactersRouter.POST("/:characterId/tags/:tagId", func(c *gin.Context) {
		characterId, err := getID(c, "characterId")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}
		tagId, err := getID(c, "tagId")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
			return
		}

		err = model.AddCharacterTag(db, characterId, tagId)
		respondEmpty(c, err)
	})

	charactersRouter.DELETE("/:characterId/tags/:tagId", func(c *gin.Context) {
		characterId, err := getID(c, "characterId")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}
		tagId, err := getID(c, "tagId")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
			return
		}

		err = model.RemoveCharacterTag(db, characterId, tagId)
		respondDeleted(c, err)
	})

	charactersRouter.GET("/:characterId/dialogue-examples", func(c *gin.Context) {
		characterId, err := getID(c, "characterId")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		examples, err := model.DialogueExamplesByCharacterId(db, characterId)
		respondList(c, examples, err)
	})

	charactersRouter.POST("/:characterId/dialogue-examples", func(c *gin.Context) {
		characterId, err := getID(c, "characterId")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		var examples []string
		if err := c.ShouldBind(&examples); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid dialogue examples data"})
			return
		}

		err = model.ReplaceDialogueExamplesByCharacterId(db, characterId, examples)
		respondEmpty(c, err)
	})

	charactersRouter.GET("/:characterId/greetings", func(c *gin.Context) {
		characterId, err := getID(c, "characterId")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		greetings, err := model.CharacterGreetingsByCharacterId(db, characterId)
		respondList(c, greetings, err)
	})

	charactersRouter.POST("/:characterId/greetings", func(c *gin.Context) {
		characterId, err := getID(c, "characterId")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		var greetings []string
		if err := c.ShouldBind(&greetings); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid greetings data"})
			return
		}

		err = model.ReplaceGreetingsByCharacterId(db, characterId, greetings)
		respondEmpty(c, err)
	})

	charactersRouter.GET("/:characterId/group-greetings", func(c *gin.Context) {
		characterId, err := getID(c, "characterId")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		greetings, err := model.CharacterGreetingsByCharacterId(db, characterId)
		respondList(c, greetings, err)
	})

	charactersRouter.POST("/:characterId/group-greetings", func(c *gin.Context) {
		characterId, err := getID(c, "characterId")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		var greetings []string
		if err := c.ShouldBind(&greetings); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid greetings data"})
			return
		}

		err = model.ReplaceGroupGreetingsByCharacterId(db, characterId, greetings)
		respondEmpty(c, err)
	})
}
