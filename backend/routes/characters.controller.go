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

	charactersRouter.GET("/:id", func(c *gin.Context) {
		id, err := getID(c, "id")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		character, err := model.CharacterById(db, id)
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

	charactersRouter.PUT("/:id", func(c *gin.Context) {
		id, err := getID(c, "id")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
			return
		}

		var character model.Character
		if err := c.ShouldBind(&character); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character data"})
			return
		}

		err = model.UpdateCharacter(db, id, &character)
		respondEmpty(c, err)
	})

	charactersRouter.DELETE("/:id", func(c *gin.Context) {
		id, err := getID(c, "id")
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid character ID"})
			return
		}

		err = model.DeleteCharacterById(db, id)
		respondDeleted(c, err)
	})
}
