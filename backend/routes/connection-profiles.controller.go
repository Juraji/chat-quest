package routes

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/ai"
	"juraji.nl/chat-quest/model"
)

func ConnectionProfilesController(router *gin.RouterGroup, db *sql.DB) {
	connectionProfilesRouter := router.Group("/connection-profiles")

	connectionProfilesRouter.GET("", func(c *gin.Context) {
		profiles, err := model.AllConnectionProfiles(db)
		respondList(c, profiles, err)
	})

	connectionProfilesRouter.GET("/:profileId", func(c *gin.Context) {
		profileId, err := getIDParam(c, "profileId")
		if err != nil {
			respondBadRequest(c, "Invalid connection profile ID")
			return
		}

		profile, err := model.ConnectionProfileById(db, profileId)
		respondSingle(c, profile, err)
	})

	connectionProfilesRouter.POST("", func(c *gin.Context) {
		var newProfile model.ConnectionProfile
		if err := c.ShouldBindJSON(&newProfile); err != nil {
			respondBadRequest(c, "Invalid connection profile data")
			return
		}

		llmModels, err := ai.GetAvailableModels(newProfile)
		if err != nil {
			respondNotAcceptable(c, "Connection test failed (Failed to get available models)", err)
			return
		}

		err = model.CreateConnectionProfile(db, &newProfile, llmModels)
		respondSingle(c, &newProfile, err)
	})

	connectionProfilesRouter.PUT("/:profileId", func(c *gin.Context) {
		profileId, err := getIDParam(c, "profileId")
		if err != nil {
			respondBadRequest(c, "Invalid connection profile ID")
			return
		}
		var profile model.ConnectionProfile
		if err := c.ShouldBindJSON(&profile); err != nil {
			respondBadRequest(c, "Invalid connection profile data")
			return
		}

		err = model.UpdateConnectionProfile(db, profileId, &profile)
		respondSingle(c, &profile, err)
	})

	connectionProfilesRouter.DELETE("/:profileId", func(c *gin.Context) {
		profileId, err := getIDParam(c, "profileId")
		if err != nil {
			respondBadRequest(c, "Invalid connection profile ID")
			return
		}

		err = model.DeleteConnectionProfileById(db, profileId)
		respondDeleted(c, err)
	})

	connectionProfilesRouter.GET("/:profileId/models", func(c *gin.Context) {
		profileId, err := getIDParam(c, "profileId")
		if err != nil {
			respondBadRequest(c, "Invalid connection profile ID")
			return
		}

		models, err := model.LlmModelsByConnectionProfileId(db, profileId)
		respondList(c, models, err)
	})

	connectionProfilesRouter.POST("/:profileId/models/refresh", func(c *gin.Context) {
		profileId, err := getIDParam(c, "profileId")
		if err != nil {
			respondBadRequest(c, "Invalid connection profile ID")
			return
		}

		profile, err := model.ConnectionProfileById(db, profileId)
		if err != nil {
			respondEmpty(c, err)
			return
		}

		llmModels, err := ai.GetAvailableModels(*profile)
		if err != nil {
			respondNotAcceptable(c, "Connection test failed (Failed to get available models)", err)
			return
		}

		err = model.MergeLlmModels(db, profileId, llmModels)
		respondEmpty(c, err)
	})

	connectionProfilesRouter.POST("/:profileId/models", func(c *gin.Context) {
		profileId, err := getIDParam(c, "profileId")
		if err != nil {
			respondBadRequest(c, "Invalid connection profile ID")
			return
		}

		var newLlmModel model.LlmModel
		if err := c.ShouldBindJSON(&newLlmModel); err != nil {
			respondBadRequest(c, "Invalid llm model data")
			return
		}

		err = model.CreateLlmModel(db, profileId, &newLlmModel)
		respondSingle(c, &newLlmModel, err)
	})

	connectionProfilesRouter.PUT("/:profileId/models/:modelId", func(c *gin.Context) {
		modelId, err := getIDParam(c, "modelId")
		if err != nil {
			respondBadRequest(c, "Invalid model ID")
			return
		}

		var llmModel model.LlmModel
		if err := c.ShouldBindJSON(&llmModel); err != nil {
			respondBadRequest(c, "Invalid model data")
			return
		}

		err = model.UpdateLlmModel(db, modelId, &llmModel)
		respondSingle(c, &llmModel, err)
	})

	connectionProfilesRouter.DELETE("/:profileId/models/:modelId", func(c *gin.Context) {
		modelId, err := getIDParam(c, "modelId")
		if err != nil {
			respondBadRequest(c, "Invalid model ID")
			return
		}

		err = model.DeleteLlmModelById(db, modelId)
		respondDeleted(c, err)
	})
}
