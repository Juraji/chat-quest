package routes

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/model"
	"juraji.nl/chat-quest/util"
)

func ChatSessionsController(router *gin.RouterGroup, db *sql.DB) {
	sessionRouter := router.Group("/worlds/:worldId/chat-sessions")

	sessionRouter.GET("", func(c *gin.Context) {
		worldId, err := util.GetIDParam(c, "worldId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid world ID")
			return
		}

		sessions, err := model.GetAllChatSessionsByWorldId(db, worldId)
		util.RespondList(c, sessions, err)
	})

	sessionRouter.GET("/:sessionId", func(c *gin.Context) {
		worldId, err := util.GetIDParam(c, "worldId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid world ID")
			return
		}
		sessionId, err := util.GetIDParam(c, "sessionId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid session ID")
			return
		}

		session, err := model.GetChatSessionById(db, worldId, sessionId)
		util.RespondSingle(c, session, err)
	})

	sessionRouter.POST("", func(c *gin.Context) {
		worldId, err := util.GetIDParam(c, "worldId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid world ID")
			return
		}

		var session model.ChatSession
		if err := c.ShouldBindJSON(&session); err != nil {
			util.RespondBadRequest(c, "Invalid session data")
			return
		}

		err = model.CreateChatSession(db, worldId, &session)
		util.RespondSingle(c, &session, err)
	})

	sessionRouter.PUT("/:sessionId", func(c *gin.Context) {
		worldId, err := util.GetIDParam(c, "worldId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid world ID")
			return
		}
		sessionId, err := util.GetIDParam(c, "sessionId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid session ID")
			return
		}

		var session model.ChatSession
		if err := c.ShouldBindJSON(&session); err != nil {
			util.RespondBadRequest(c, "Invalid session data")
			return
		}

		err = model.UpdateChatSession(db, worldId, sessionId, &session)
		util.RespondSingle(c, &session, err)
	})

	sessionRouter.DELETE("/:sessionId", func(c *gin.Context) {
		worldId, err := util.GetIDParam(c, "worldId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid world ID")
			return
		}
		sessionId, err := util.GetIDParam(c, "sessionId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid session ID")
			return
		}

		err = model.DeleteChatSessionById(db, worldId, sessionId)
		util.RespondDeleted(c, err)
	})

	sessionRouter.GET("/:sessionId/chat-messages", func(c *gin.Context) {
		sessionId, err := util.GetIDParam(c, "sessionId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid session ID")
			return
		}

		messages, err := model.GetChatMessages(db, sessionId)
		util.RespondList(c, messages, err)
	})

	sessionRouter.POST("/:sessionId/chat-messages", func(c *gin.Context) {
		sessionId, err := util.GetIDParam(c, "sessionId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid session ID")
			return
		}

		var message model.ChatMessage
		if err := c.ShouldBindJSON(&message); err != nil {
			util.RespondBadRequest(c, "Invalid chat message data")
			return
		}

		// TODO: Trigger response (embedding, fetching memories, prompt building, chat completion)
		//       go func...
		// TODO: Trigger chat truncation (creating memories)

		err = model.CreateChatMessage(db, sessionId, &message)
		util.RespondSingle(c, &message, err)
	})

	sessionRouter.PUT("/:sessionId/chat-messages/:messageId", func(c *gin.Context) {
		sessionId, err := util.GetIDParam(c, "sessionId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid session ID")
			return
		}

		messageId, err := util.GetIDParam(c, "messageId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid chat message ID")
			return
		}

		var message model.ChatMessage
		if err := c.ShouldBindJSON(&message); err != nil {
			util.RespondBadRequest(c, "Invalid chat message data")
			return
		}

		// TODO: Trigger updating of memories

		err = model.UpdateChatMessage(db, sessionId, messageId, &message)
		util.RespondSingle(c, &message, err)
	})

	sessionRouter.DELETE("/:sessionId/chat-messages/:messageId", func(c *gin.Context) {
		sessionId, err := util.GetIDParam(c, "sessionId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid session ID")
			return
		}
		messageId, err := util.GetIDParam(c, "messageId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid chat message ID")
			return
		}

		err = model.DeleteChatMessagesFrom(db, sessionId, messageId)
		util.RespondDeleted(c, err)
	})
}
