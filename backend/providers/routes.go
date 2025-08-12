package providers

import (
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/cq"
	"juraji.nl/chat-quest/util"
)

func Routes(cq *cq.ChatQuestContext, router *gin.RouterGroup) {
	connectionProfilesRouter := router.Group("/connection-profiles")

	connectionProfilesRouter.GET("", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		profiles, err := AllConnectionProfiles(cq)
		util.RespondList(c, profiles, err)
	})

	connectionProfilesRouter.GET("/:profileId", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		profileId, err := util.GetIDParam(c, "profileId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid connection profile ID")
			return
		}

		profile, err := ConnectionProfileById(cq, profileId)
		util.RespondSingle(c, profile, err)
	})

	connectionProfilesRouter.POST("", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		var newProfile ConnectionProfile
		if err := c.ShouldBindJSON(&newProfile); err != nil {
			util.RespondBadRequest(c, "Invalid connection profile data")
			return
		}

		llmModels, err := newProfile.GetAvailableModels()
		if err != nil {
			util.RespondNotAcceptable(c, "Connection test failed (Failed to get available models)", err)
			return
		}

		err = CreateConnectionProfile(cq, &newProfile, llmModels)
		util.RespondSingle(c, &newProfile, err)
	})

	connectionProfilesRouter.PUT("/:profileId", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		profileId, err := util.GetIDParam(c, "profileId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid connection profile ID")
			return
		}
		var profile ConnectionProfile
		if err := c.ShouldBindJSON(&profile); err != nil {
			util.RespondBadRequest(c, "Invalid connection profile data")
			return
		}

		err = UpdateConnectionProfile(cq, profileId, &profile)
		util.RespondSingle(c, &profile, err)
	})

	connectionProfilesRouter.DELETE("/:profileId", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		profileId, err := util.GetIDParam(c, "profileId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid connection profile ID")
			return
		}

		err = DeleteConnectionProfileById(cq, profileId)
		util.RespondDeleted(c, err)
	})

	connectionProfilesRouter.GET("/:profileId/models", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		profileId, err := util.GetIDParam(c, "profileId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid connection profile ID")
			return
		}

		models, err := LlmModelsByConnectionProfileId(cq, profileId)
		util.RespondList(c, models, err)
	})

	connectionProfilesRouter.POST("/:profileId/models/refresh", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		profileId, err := util.GetIDParam(c, "profileId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid connection profile ID")
			return
		}

		profile, err := ConnectionProfileById(cq, profileId)
		if err != nil {
			util.RespondEmpty(c, err)
			return
		}

		llmModels, err := profile.GetAvailableModels()
		if err != nil {
			util.RespondNotAcceptable(c, "Connection test failed (Failed to get available models)", err)
			return
		}

		err = MergeLlmModels(cq, profileId, llmModels)
		util.RespondEmpty(c, err)
	})

	connectionProfilesRouter.POST("/:profileId/models", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		profileId, err := util.GetIDParam(c, "profileId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid connection profile ID")
			return
		}

		var newLlmModel LlmModel
		if err := c.ShouldBindJSON(&newLlmModel); err != nil {
			util.RespondBadRequest(c, "Invalid llm model data")
			return
		}

		err = CreateLlmModel(cq, profileId, &newLlmModel)
		util.RespondSingle(c, &newLlmModel, err)
	})

	connectionProfilesRouter.PUT("/:profileId/models/:modelId", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		modelId, err := util.GetIDParam(c, "modelId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid model ID")
			return
		}

		var llmModel LlmModel
		if err := c.ShouldBindJSON(&llmModel); err != nil {
			util.RespondBadRequest(c, "Invalid model data")
			return
		}

		err = UpdateLlmModel(cq, modelId, &llmModel)
		util.RespondSingle(c, &llmModel, err)
	})

	connectionProfilesRouter.DELETE("/:profileId/models/:modelId", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		modelId, err := util.GetIDParam(c, "modelId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid model ID")
			return
		}

		err = DeleteLlmModelById(cq, modelId)
		util.RespondDeleted(c, err)
	})

	connectionProfilesRouter.GET("/model-views", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		views, err := GetAllLlmModelViews(cq)
		util.RespondList(c, views, err)
	})
}
