package routes

import (
	"chat-quest/backend/model"
	"database/sql"
	"github.com/gin-gonic/gin"
)

func CharactersController(router *gin.RouterGroup, db *sql.DB) {
	charactersRouter := router.Group("/characters")

	charactersRouter.GET("/", func(c *gin.Context) {
		characters, err := model.AllCharacters(db)
		respondList(c, characters, err)
	})

	charactersRouter.GET("/:id", func(c *gin.Context) {
		id, err := getID(c, "id")
		if id == 0 {
			c.JSON(400, gin.H{"error": "Invalid character ID"})
			return
		}

		character, err := model.CharacterById(db, id)
		respondSingle(c, character, err)
	})
}
