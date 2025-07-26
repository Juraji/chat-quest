package routes

import (
	"chat-quest/backend/model"
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
)

func CharactersController(router *gin.RouterGroup, db *sql.DB) {
	charactersRouter := router.Group("/characters")

	charactersRouter.GET("/", func(c *gin.Context) {
		characters, err := model.AllCharacters(db)

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			log.Fatal(err)
		} else {
			c.JSON(200, characters)
		}
	})

	//charactersRouter.GET("/:id", func(c *gin.Context) {
	//	id := c.Param("id")
	//	character, err := model.CharacterById(db, id)
	//
	//  if character == nil {
	//    c.JSON(404, gin.H{"error": "Character not found"})
	//  } else {
	//    c.JSON(200, character)
	//  }
	//})
}
