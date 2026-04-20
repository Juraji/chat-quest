package api

import (
	"github.com/gin-gonic/gin"
	preferences2 "juraji.nl/chat-quest/model/preferences"
)

func PreferencesRoutes(router *gin.RouterGroup) {
	preferencesRouter := router.Group("/preferences")

	preferencesRouter.GET("", func(c *gin.Context) {
		preferences, err := preferences2.GetPreferences(false)
		respondSingle(c, preferences, err)
	})

	preferencesRouter.PUT("", func(c *gin.Context) {
		var update *preferences2.Preferences
		if err := c.ShouldBindJSON(&update); err != nil {
			respondBadRequest(c, "Invalid preferences", err)
			return
		}

		err := preferences2.UpdatePreferences(update)
		respondSingle(c, update, err)
	})

	preferencesRouter.GET("/validate", func(c *gin.Context) {
		prefs, err := preferences2.GetPreferences(false)
		if err != nil {
			respondInternalError(c, err)
			return
		}

		errs := prefs.Validate()
		respondList(c, errs, nil)
	})
}
