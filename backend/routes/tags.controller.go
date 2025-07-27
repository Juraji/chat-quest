package routes

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/model"
	"net/http"
)

func TagsController(router *gin.RouterGroup, db *sql.DB) {
	tagsRouter := router.Group("/tags")

	tagsRouter.GET("", func(c *gin.Context) {
		tags, err := model.AllTags(db)
		respondList(c, tags, err)
	})

	tagsRouter.GET("/:id", func(c *gin.Context) {
		id, err := getID(c, "id")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
			return
		}

		tag, err := model.TagById(db, id)
		respondSingle(c, tag, err)
	})

	tagsRouter.POST("/", func(c *gin.Context) {
		var newTag model.Tag
		if err := c.ShouldBind(&newTag); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag data"})
			return
		}

		err := model.CreateTag(db, &newTag)
		respondSingle(c, &newTag, err)
	})

	tagsRouter.PUT("/:id", func(c *gin.Context) {
		id, err := getID(c, "id")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
			return
		}

		var tag model.Tag
		if err := c.ShouldBind(&tag); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag data"})
			return
		}

		err = model.UpdateTag(db, id, &tag)
		respondSingle(c, &tag, err)
	})

	tagsRouter.DELETE("/:id", func(c *gin.Context) {
		id, err := getID(c, "id")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
			return
		}

		err = model.DeleteTagById(db, id)
		respondDeleted(c, err)
	})
}
