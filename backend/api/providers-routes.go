package api

import (
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/core/providers"
)

func ProvidersRoutes(router *gin.RouterGroup) {
	connectionProfilesRouter := router.Group("/connection-profiles")

	connectionProfilesRouter.GET("", func(c *gin.Context) {
		profiles, err := providers.AllConnectionProfiles()
		respondList(c, profiles, err)
	})

	connectionProfilesRouter.GET("/:profileId", func(c *gin.Context) {
		profileId, ok := getParamAsID(c, "profileId")
		if !ok {
			respondBadRequest(c, "Invalid connection profile ID", nil)
			return
		}

		profile, err := providers.ConnectionProfileById(profileId)
		respondSingle(c, profile, err)
	})

	connectionProfilesRouter.POST("", func(c *gin.Context) {
		var newProfile providers.ConnectionProfile
		if err := c.ShouldBindJSON(&newProfile); err != nil {
			respondBadRequest(c, "Invalid connection profile data", nil)
			return
		}

		llmModels, err := providers.GetAvailableModels(&newProfile)
		if err != nil {
			respondNotAcceptable(c, "Connection test failed (Failed to get available models)", err)
			return
		}

		err = providers.CreateConnectionProfile(&newProfile, llmModels)
		respondSingle(c, &newProfile, err)
	})

	connectionProfilesRouter.PUT("/:profileId", func(c *gin.Context) {
		profileId, ok := getParamAsID(c, "profileId")
		if !ok {
			respondBadRequest(c, "Invalid connection profile ID", nil)
			return
		}
		var profile providers.ConnectionProfile
		if err := c.ShouldBindJSON(&profile); err != nil {
			respondBadRequest(c, "Invalid connection profile data", nil)
			return
		}
		if !profile.ProviderType.IsValid() {
			respondBadRequest(c, "Invalid connection profile type", nil)
			return
		}

		err := providers.UpdateConnectionProfile(profileId, &profile)
		respondSingle(c, &profile, err)
	})

	connectionProfilesRouter.DELETE("/:profileId", func(c *gin.Context) {
		profileId, ok := getParamAsID(c, "profileId")
		if !ok {
			respondBadRequest(c, "Invalid connection profile ID", nil)
			return
		}

		err := providers.DeleteConnectionProfileById(profileId)
		respondEmpty(c, err)
	})

	connectionProfilesRouter.GET("/:profileId/models", func(c *gin.Context) {
		profileId, ok := getParamAsID(c, "profileId")
		if !ok {
			respondBadRequest(c, "Invalid connection profile ID", nil)
			return
		}

		models, err := providers.LlmModelsByConnectionProfileId(profileId)
		respondList(c, models, err)
	})

	connectionProfilesRouter.POST("/:profileId/models/refresh", func(c *gin.Context) {
		profileId, ok := getParamAsID(c, "profileId")
		if !ok {
			respondBadRequest(c, "Invalid connection profile ID", nil)
			return
		}

		profile, err := providers.ConnectionProfileById(profileId)
		if err != nil {
			respondInternalError(c, err)
			return
		}

		llmModels, err := providers.GetAvailableModels(profile)
		if err != nil {
			respondNotAcceptable(c, "Connection test failed (Failed to get available models)", err)
			return
		}

		err = providers.MergeLlmModels(profileId, llmModels)
		respondEmpty(c, err)
	})

	connectionProfilesRouter.PUT("/:profileId/models/:modelId", func(c *gin.Context) {
		modelId, ok := getParamAsID(c, "modelId")
		if !ok {
			respondBadRequest(c, "Invalid model ID", nil)
			return
		}

		var llmModel providers.LlmModel
		if err := c.ShouldBindJSON(&llmModel); err != nil {
			respondBadRequest(c, "Invalid model data", nil)
			return
		}

		if !llmModel.ModelType.IsValid() {
			respondBadRequest(c, "Invalid model type", nil)
			return
		}

		err := providers.UpdateLlmModel(modelId, &llmModel)
		respondSingle(c, &llmModel, err)
	})

	connectionProfilesRouter.DELETE("/:profileId/models/:modelId", func(c *gin.Context) {
		modelId, ok := getParamAsID(c, "modelId")
		if !ok {
			respondBadRequest(c, "Invalid model ID", nil)
			return
		}

		err := providers.DeleteLlmModelById(modelId)
		respondEmpty(c, err)
	})

	connectionProfilesRouter.GET("/model-views", func(c *gin.Context) {
		views, err := providers.GetAllLlmModelViews()
		respondList(c, views, err)
	})
}
