package preferences

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/util"
)

func Routes(router *gin.RouterGroup, db *sql.DB) {
	prefsRouter := router.Group("/chat-preferences")

	prefsRouter.GET("", func(c *gin.Context) {
		prefs, err := GetChatPreferences(db)
		util.RespondSingle(c, prefs, err)
	})

	prefsRouter.PUT("", func(c *gin.Context) {
		var update ChatPreferences
		if err := c.ShouldBind(&update); err != nil {
			util.RespondBadRequest(c, "Invalid preference data")
			return
		}

		err := UpdateChatPreferences(db, &update)
		util.RespondSingle(c, &update, err)
	})
}
