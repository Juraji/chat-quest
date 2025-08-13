package providers

import (
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/core"
	"juraji.nl/chat-quest/core/util"
)

func Routes(cq *core.ChatQuestContext, router *gin.RouterGroup) {
	connectionProfilesRouter := router.Group("/connection-profiles")

	connectionProfilesRouter.GET("", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		profiles, err := AllConnectionProfiles(cq)
		util.RespondList(rcq, c, profiles, err)
	})

	connectionProfilesRouter.GET("/:profileId", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		profileId, err := util.GetIDParam(c, "profileId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid connection profile ID")
			return
		}

		profile, err := ConnectionProfileById(rcq, profileId)
		util.RespondSingle(rcq, c, profile, err)
	})

	connectionProfilesRouter.POST("", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		var newProfile ConnectionProfile
		if err := c.ShouldBindJSON(&newProfile); err != nil {
			util.RespondBadRequest(rcq, c, "Invalid connection profile data")
			return
		}

		llmModels, err := newProfile.GetAvailableModels()
		if err != nil {
			util.RespondNotAcceptable(rcq, c, "Connection test failed (Failed to get available models)", err)
			return
		}

		err = CreateConnectionProfile(rcq, &newProfile, llmModels)
		util.RespondSingle(rcq, c, &newProfile, err)
	})

	connectionProfilesRouter.PUT("/:profileId", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		profileId, err := util.GetIDParam(c, "profileId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid connection profile ID")
			return
		}
		var profile ConnectionProfile
		if err := c.ShouldBindJSON(&profile); err != nil {
			util.RespondBadRequest(rcq, c, "Invalid connection profile data")
			return
		}
		if !profile.ProviderType.IsValid() {
			util.RespondBadRequest(rcq, c, "Invalid connection profile type")
			return
		}

		err = UpdateConnectionProfile(rcq, profileId, &profile)
		util.RespondSingle(rcq, c, &profile, err)
	})

	connectionProfilesRouter.DELETE("/:profileId", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		profileId, err := util.GetIDParam(c, "profileId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid connection profile ID")
			return
		}

		err = DeleteConnectionProfileById(rcq, profileId)
		util.RespondDeleted(rcq, c, err)
	})

	connectionProfilesRouter.GET("/:profileId/models", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		profileId, err := util.GetIDParam(c, "profileId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid connection profile ID")
			return
		}

		models, err := LlmModelsByConnectionProfileId(rcq, profileId)
		util.RespondList(rcq, c, models, err)
	})

	connectionProfilesRouter.POST("/:profileId/models/refresh", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		profileId, err := util.GetIDParam(c, "profileId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid connection profile ID")
			return
		}

		profile, err := ConnectionProfileById(rcq, profileId)
		if err != nil {
			util.RespondEmpty(rcq, c, err)
			return
		}

		llmModels, err := profile.GetAvailableModels()
		if err != nil {
			util.RespondNotAcceptable(rcq, c, "Connection test failed (Failed to get available models)", err)
			return
		}

		err = MergeLlmModels(rcq, profileId, llmModels)
		util.RespondEmpty(rcq, c, err)
	})

	connectionProfilesRouter.POST("/:profileId/models", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		profileId, err := util.GetIDParam(c, "profileId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid connection profile ID")
			return
		}

		var newLlmModel LlmModel
		if err := c.ShouldBindJSON(&newLlmModel); err != nil {
			util.RespondBadRequest(rcq, c, "Invalid llm model data")
			return
		}

		err = CreateLlmModel(rcq, profileId, &newLlmModel)
		util.RespondSingle(rcq, c, &newLlmModel, err)
	})

	connectionProfilesRouter.PUT("/:profileId/models/:modelId", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		modelId, err := util.GetIDParam(c, "modelId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid model ID")
			return
		}

		var llmModel LlmModel
		if err := c.ShouldBindJSON(&llmModel); err != nil {
			util.RespondBadRequest(rcq, c, "Invalid model data")
			return
		}

		err = UpdateLlmModel(rcq, modelId, &llmModel)
		util.RespondSingle(rcq, c, &llmModel, err)
	})

	connectionProfilesRouter.DELETE("/:profileId/models/:modelId", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		modelId, err := util.GetIDParam(c, "modelId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid model ID")
			return
		}

		err = DeleteLlmModelById(rcq, modelId)
		util.RespondDeleted(rcq, c, err)
	})

	connectionProfilesRouter.GET("/model-views", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		views, err := GetAllLlmModelViews(cq)
		util.RespondList(rcq, c, views, err)
	})
}
