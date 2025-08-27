package preferences

import (
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/core/util/controllers"
)

func Routes(router *gin.RouterGroup) {
	preferencesRouter := router.Group("/preferences")

	preferencesRouter.GET("", func(c *gin.Context) {
		preferences, err := GetPreferences()
		controllers.RespondSingleE(c, preferences, err)
	})

	preferencesRouter.PUT("", func(c *gin.Context) {
		var update *Preferences
		if err := c.ShouldBindJSON(&update); err != nil {
			controllers.RespondBadRequestE(c, err)
			return
		}

		err := UpdatePreferences(update)
		controllers.RespondSingleE(c, update, err)
	})
}
