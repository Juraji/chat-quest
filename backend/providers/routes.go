package providers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/util"
)

func Routes(router *gin.RouterGroup, db *sql.DB) {
	connectionProfilesRouter := router.Group("/connection-profiles")

	connectionProfilesRouter.GET("", func(c *gin.Context) {
		profiles, err := AllConnectionProfiles(db)
		util.RespondList(c, profiles, err)
	})

	connectionProfilesRouter.GET("/:profileId", func(c *gin.Context) {
		profileId, err := util.GetIDParam(c, "profileId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid connection profile ID")
			return
		}

		profile, err := ConnectionProfileById(db, profileId)
		util.RespondSingle(c, profile, err)
	})

	connectionProfilesRouter.POST("", func(c *gin.Context) {
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

		err = CreateConnectionProfile(db, &newProfile, llmModels)
		util.RespondSingle(c, &newProfile, err)
	})

	connectionProfilesRouter.PUT("/:profileId", func(c *gin.Context) {
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

		err = UpdateConnectionProfile(db, profileId, &profile)
		util.RespondSingle(c, &profile, err)
	})

	connectionProfilesRouter.DELETE("/:profileId", func(c *gin.Context) {
		profileId, err := util.GetIDParam(c, "profileId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid connection profile ID")
			return
		}

		err = DeleteConnectionProfileById(db, profileId)
		util.RespondDeleted(c, err)
	})

	connectionProfilesRouter.GET("/:profileId/models", func(c *gin.Context) {
		profileId, err := util.GetIDParam(c, "profileId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid connection profile ID")
			return
		}

		models, err := LlmModelsByConnectionProfileId(db, profileId)
		util.RespondList(c, models, err)
	})

	connectionProfilesRouter.POST("/:profileId/models/refresh", func(c *gin.Context) {
		profileId, err := util.GetIDParam(c, "profileId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid connection profile ID")
			return
		}

		profile, err := ConnectionProfileById(db, profileId)
		if err != nil {
			util.RespondEmpty(c, err)
			return
		}

		llmModels, err := profile.GetAvailableModels()
		if err != nil {
			util.RespondNotAcceptable(c, "Connection test failed (Failed to get available models)", err)
			return
		}

		err = MergeLlmModels(db, profileId, llmModels)
		util.RespondEmpty(c, err)
	})

	connectionProfilesRouter.POST("/:profileId/models", func(c *gin.Context) {
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

		err = CreateLlmModel(db, profileId, &newLlmModel)
		util.RespondSingle(c, &newLlmModel, err)
	})

	connectionProfilesRouter.PUT("/:profileId/models/:modelId", func(c *gin.Context) {
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

		err = UpdateLlmModel(db, modelId, &llmModel)
		util.RespondSingle(c, &llmModel, err)
	})

	connectionProfilesRouter.DELETE("/:profileId/models/:modelId", func(c *gin.Context) {
		modelId, err := util.GetIDParam(c, "modelId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid model ID")
			return
		}

		err = DeleteLlmModelById(db, modelId)
		util.RespondDeleted(c, err)
	})

	connectionProfilesRouter.GET("/model-views", func(c *gin.Context) {
		views, err := GetAllLlmModelViews(db)
		util.RespondList(c, views, err)
	})
}
