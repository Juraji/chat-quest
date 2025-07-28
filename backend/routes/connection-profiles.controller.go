package routes

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/model"
)

func ConnectionProfilesController(router *gin.RouterGroup, db *sql.DB) {
	connectionProfilesRouter := router.Group("/connection-profiles")

	connectionProfilesRouter.GET("", func(c *gin.Context) {
		profiles, err := model.AllConnectionProfiles(db)
		respondList(c, profiles, err)
	})

	connectionProfilesRouter.GET("/templates/:providerType", func(c *gin.Context) {
		providerType := c.Param("providerType")
		profile, err := model.GetConnectionProfileDefaults(providerType)
		if err != nil {
			respondSingle[model.ConnectionProfile](c, nil, nil)
		} else {
			respondSingle(c, &profile, nil)
		}
	})

	connectionProfilesRouter.GET("/:profileId", func(c *gin.Context) {
		profileId, err := getID(c, "profileId")
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid connection profile ID"})
			return
		}

		profile, err := model.ConnectionProfileById(db, profileId)
		respondSingle(c, profile, err)
	})

	connectionProfilesRouter.POST("", func(c *gin.Context) {
		var newProfile model.ConnectionProfile
		if err := c.ShouldBindJSON(&newProfile); err != nil {
			c.JSON(400, gin.H{"error": "Invalid connection profile data"})
			return
		}

		err := model.CreateConnectionProfile(db, &newProfile)
		// TODO: Test connection and refresh models
		respondSingle(c, &newProfile, err)
	})

	connectionProfilesRouter.PUT("/:profileId", func(c *gin.Context) {
		profileId, err := getID(c, "profileId")
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid connection profile ID"})
			return
		}
		var profile model.ConnectionProfile
		if err := c.ShouldBindJSON(&profile); err != nil {
			c.JSON(400, gin.H{"error": "Invalid connection profile data"})
			return
		}

		err = model.UpdateConnectionProfile(db, profileId, &profile)
		// TODO: Test connection and refresh models
		respondSingle(c, &profile, err)
	})

	connectionProfilesRouter.DELETE("/:profileId", func(c *gin.Context) {
		profileId, err := getID(c, "profileId")
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid connection profile ID"})
			return
		}

		err = model.DeleteConnectionProfileById(db, profileId)
		respondDeleted(c, err)
	})

	connectionProfilesRouter.GET("/:profileId/models", func(c *gin.Context) {
		profileId, err := getID(c, "profileId")
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid connection profile ID"})
			return
		}

		models, err := model.LlmModelsByConnectionProfileId(db, profileId)
		respondList(c, models, err)
	})

	connectionProfilesRouter.POST("/:profileId/models/refresh", func(c *gin.Context) {})

	connectionProfilesRouter.PUT("/:profileId/models/:modelId", func(c *gin.Context) {
		modelId, err := getID(c, "modelId")
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid model ID"})
			return
		}

		var llmModel model.LlmModel
		if err := c.ShouldBindJSON(&llmModel); err != nil {
			c.JSON(400, gin.H{"error": "Invalid model data"})
			return
		}

		err = model.UpdateLlmModel(db, modelId, &llmModel)
		respondSingle(c, &llmModel, err)
	})
}
