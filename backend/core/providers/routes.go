package providers

import (
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/core/util/controllers"
)

func Routes(router *gin.RouterGroup) {
	connectionProfilesRouter := router.Group("/connection-profiles")

	connectionProfilesRouter.GET("", func(c *gin.Context) {
		profiles, err := AllConnectionProfiles()
		controllers.RespondListE(c, profiles, err)
	})

	connectionProfilesRouter.GET("/:profileId", func(c *gin.Context) {
		profileId, ok := controllers.GetParamAsID(c, "profileId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid connection profile ID", nil)
			return
		}

		profile, err := ConnectionProfileById(profileId)
		controllers.RespondSingleE(c, profile, err)
	})

	connectionProfilesRouter.POST("", func(c *gin.Context) {
		var newProfile ConnectionProfile
		if err := c.ShouldBindJSON(&newProfile); err != nil {
			controllers.RespondBadRequest(c, "Invalid connection profile data", nil)
			return
		}

		llmModels, err := GetAvailableModels(&newProfile)
		if err != nil {
			controllers.RespondNotAcceptable(c, "Connection test failed (Failed to get available models)", err)
			return
		}

		err = CreateConnectionProfile(&newProfile, llmModels)
		controllers.RespondSingleE(c, &newProfile, err)
	})

	connectionProfilesRouter.PUT("/:profileId", func(c *gin.Context) {
		profileId, ok := controllers.GetParamAsID(c, "profileId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid connection profile ID", nil)
			return
		}
		var profile ConnectionProfile
		if err := c.ShouldBindJSON(&profile); err != nil {
			controllers.RespondBadRequest(c, "Invalid connection profile data", nil)
			return
		}
		if !profile.ProviderType.IsValid() {
			controllers.RespondBadRequest(c, "Invalid connection profile type", nil)
			return
		}

		err := UpdateConnectionProfile(profileId, &profile)
		controllers.RespondSingleE(c, &profile, err)
	})

	connectionProfilesRouter.DELETE("/:profileId", func(c *gin.Context) {
		profileId, ok := controllers.GetParamAsID(c, "profileId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid connection profile ID", nil)
			return
		}

		err := DeleteConnectionProfileById(profileId)
		controllers.RespondEmptyE(c, err)
	})

	connectionProfilesRouter.GET("/:profileId/models", func(c *gin.Context) {
		profileId, ok := controllers.GetParamAsID(c, "profileId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid connection profile ID", nil)
			return
		}

		models, err := LlmModelsByConnectionProfileId(profileId)
		controllers.RespondListE(c, models, err)
	})

	connectionProfilesRouter.POST("/:profileId/models/refresh", func(c *gin.Context) {
		profileId, ok := controllers.GetParamAsID(c, "profileId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid connection profile ID", nil)
			return
		}

		profile, err := ConnectionProfileById(profileId)
		if err != nil {
			controllers.RespondInternalError(c, err)
			return
		}

		llmModels, err := GetAvailableModels(profile)
		if err != nil {
			controllers.RespondNotAcceptable(c, "Connection test failed (Failed to get available models)", err)
			return
		}

		err = MergeLlmModels(profileId, llmModels)
		controllers.RespondEmptyE(c, err)
	})

	connectionProfilesRouter.PUT("/:profileId/models/:modelId", func(c *gin.Context) {
		modelId, ok := controllers.GetParamAsID(c, "modelId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid model ID", nil)
			return
		}

		var llmModel LlmModel
		if err := c.ShouldBindJSON(&llmModel); err != nil {
			controllers.RespondBadRequest(c, "Invalid model data", nil)
			return
		}

		err := UpdateLlmModel(modelId, &llmModel)
		controllers.RespondSingleE(c, &llmModel, err)
	})

	connectionProfilesRouter.DELETE("/:profileId/models/:modelId", func(c *gin.Context) {
		modelId, ok := controllers.GetParamAsID(c, "modelId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid model ID", nil)
			return
		}

		err := DeleteLlmModelById(modelId)
		controllers.RespondEmptyE(c, err)
	})

	connectionProfilesRouter.GET("/model-views", func(c *gin.Context) {
		views, err := GetAllLlmModelViews()
		controllers.RespondListE(c, views, err)
	})
}
