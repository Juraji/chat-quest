package providers

import (
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/core/util/controllers"
)

func Routes(router *gin.RouterGroup) {
	connectionProfilesRouter := router.Group("/connection-profiles")

	connectionProfilesRouter.GET("", func(c *gin.Context) {
		profiles, ok := AllConnectionProfiles()
		controllers.RespondList(c, ok, profiles)
	})

	connectionProfilesRouter.GET("/:profileId", func(c *gin.Context) {
		profileId, ok := controllers.GetParamAsID(c, "profileId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid connection profile ID")
			return
		}

		profile, ok := ConnectionProfileById(profileId)
		controllers.RespondSingle(c, ok, profile)
	})

	connectionProfilesRouter.POST("", func(c *gin.Context) {
		var newProfile ConnectionProfile
		if err := c.ShouldBindJSON(&newProfile); err != nil {
			controllers.RespondBadRequest(c, "Invalid connection profile data")
			return
		}

		llmModels, err := newProfile.GetAvailableModels()
		if err != nil {
			controllers.RespondNotAcceptable(c, "Connection test failed (Failed to get available models)", err)
			return
		}

		ok := CreateConnectionProfile(&newProfile, llmModels)
		controllers.RespondSingle(c, ok, &newProfile)
	})

	connectionProfilesRouter.PUT("/:profileId", func(c *gin.Context) {
		profileId, ok := controllers.GetParamAsID(c, "profileId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid connection profile ID")
			return
		}
		var profile ConnectionProfile
		if err := c.ShouldBindJSON(&profile); err != nil {
			controllers.RespondBadRequest(c, "Invalid connection profile data")
			return
		}
		if !profile.ProviderType.IsValid() {
			controllers.RespondBadRequest(c, "Invalid connection profile type")
			return
		}

		ok = UpdateConnectionProfile(profileId, &profile)
		controllers.RespondSingle(c, ok, &profile)
	})

	connectionProfilesRouter.DELETE("/:profileId", func(c *gin.Context) {
		profileId, ok := controllers.GetParamAsID(c, "profileId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid connection profile ID")
			return
		}

		ok = DeleteConnectionProfileById(profileId)
		controllers.RespondEmpty(c, ok)
	})

	connectionProfilesRouter.GET("/:profileId/models", func(c *gin.Context) {
		profileId, ok := controllers.GetParamAsID(c, "profileId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid connection profile ID")
			return
		}

		models, ok := LlmModelsByConnectionProfileId(profileId)
		controllers.RespondList(c, ok, models)
	})

	connectionProfilesRouter.POST("/:profileId/models/refresh", func(c *gin.Context) {
		profileId, ok := controllers.GetParamAsID(c, "profileId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid connection profile ID")
			return
		}

		profile, ok := ConnectionProfileById(profileId)
		if !ok {
			controllers.RespondInternalError(c, nil)
			return
		}

		llmModels, err := profile.GetAvailableModels()
		if err != nil {
			controllers.RespondNotAcceptable(c, "Connection test failed (Failed to get available models)", err)
			return
		}

		ok = MergeLlmModels(profileId, llmModels)
		controllers.RespondEmpty(c, ok)
	})

	connectionProfilesRouter.PUT("/:profileId/models/:modelId", func(c *gin.Context) {
		modelId, ok := controllers.GetParamAsID(c, "modelId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid model ID")
			return
		}

		var llmModel LlmModel
		if err := c.ShouldBindJSON(&llmModel); err != nil {
			controllers.RespondBadRequest(c, "Invalid model data")
			return
		}

		ok = UpdateLlmModel(modelId, &llmModel)
		controllers.RespondSingle(c, ok, &llmModel)
	})

	connectionProfilesRouter.DELETE("/:profileId/models/:modelId", func(c *gin.Context) {
		modelId, ok := controllers.GetParamAsID(c, "modelId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid model ID")
			return
		}

		ok = DeleteLlmModelById(modelId)
		controllers.RespondEmpty(c, ok)
	})

	connectionProfilesRouter.GET("/model-views", func(c *gin.Context) {
		views, ok := GetAllLlmModelViews()
		controllers.RespondList(c, ok, views)
	})
}
