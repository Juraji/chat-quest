package routes

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/model"
	"juraji.nl/chat-quest/util"
)

func TagsController(router *gin.RouterGroup, db *sql.DB) {
	tagsRouter := router.Group("/tags")

	tagsRouter.GET("", func(c *gin.Context) {
		tags, err := model.AllTags(db)
		util.RespondList(c, tags, err)
	})

	tagsRouter.GET("/:id", func(c *gin.Context) {
		id, err := util.GetIDParam(c, "id")
		if err != nil {
			util.RespondBadRequest(c, "Invalid tag ID")
			return
		}

		tag, err := model.TagById(db, id)
		util.RespondSingle(c, tag, err)
	})

	tagsRouter.POST("", func(c *gin.Context) {
		var newTag model.Tag
		if err := c.ShouldBind(&newTag); err != nil {
			util.RespondBadRequest(c, "Invalid tag data")
			return
		}

		err := model.CreateTag(db, &newTag)
		util.RespondSingle(c, &newTag, err)
	})

	tagsRouter.PUT("/:id", func(c *gin.Context) {
		id, err := util.GetIDParam(c, "id")
		if err != nil {
			util.RespondBadRequest(c, "Invalid tag ID")
			return
		}

		var tag model.Tag
		if err := c.ShouldBind(&tag); err != nil {
			util.RespondBadRequest(c, "Invalid tag data")
			return
		}

		err = model.UpdateTag(db, id, &tag)
		util.RespondSingle(c, &tag, err)
	})

	tagsRouter.DELETE("/:id", func(c *gin.Context) {
		id, err := util.GetIDParam(c, "id")
		if err != nil {
			util.RespondBadRequest(c, "Invalid tag ID")
			return
		}

		err = model.DeleteTagById(db, id)
		util.RespondDeleted(c, err)
	})
}
