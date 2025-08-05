package routes

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/model"
	"juraji.nl/chat-quest/util"
)

func ChatPreferencesController(router *gin.RouterGroup, db *sql.DB) {
	prefsRouter := router.Group("/chat-preferences")

	prefsRouter.GET("", func(c *gin.Context) {
		prefs, err := model.GetChatPreferences(db)
		util.RespondSingle(c, prefs, err)
	})

	prefsRouter.PUT("", func(c *gin.Context) {
		var update model.ChatPreferences
		if err := c.ShouldBind(&update); err != nil {
			util.RespondBadRequest(c, "Invalid preference data")
			return
		}

		err := model.UpdateChatPreferences(db, &update)
		util.RespondSingle(c, &update, err)
	})
}
