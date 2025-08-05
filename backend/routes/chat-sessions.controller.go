package routes

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/model"
)

func ChatSessionsController(router *gin.RouterGroup, db *sql.DB) {
	sessionRouter := router.Group("/worlds/:worldId/chat-sessions")

	sessionRouter.GET("", func(c *gin.Context) {
		worldId, err := getIDParam(c, "worldId")
		if err != nil {
			respondBadRequest(c, "Invalid world ID")
			return
		}

		sessions, err := model.GetAllChatSessionsByWorldId(db, worldId)
		respondList(c, sessions, err)
	})

	sessionRouter.GET("/:sessionId", func(c *gin.Context) {
		worldId, err := getIDParam(c, "worldId")
		if err != nil {
			respondBadRequest(c, "Invalid world ID")
			return
		}
		sessionId, err := getIDParam(c, "sessionId")
		if err != nil {
			respondBadRequest(c, "Invalid session ID")
			return
		}

		session, err := model.GetChatSessionById(db, worldId, sessionId)
		respondSingle(c, session, err)
	})

	sessionRouter.POST("", func(c *gin.Context) {
		worldId, err := getIDParam(c, "worldId")
		if err != nil {
			respondBadRequest(c, "Invalid world ID")
			return
		}

		var session model.ChatSession
		if err := c.ShouldBindJSON(&session); err != nil {
			respondBadRequest(c, "Invalid session data")
			return
		}

		err = model.CreateChatSession(db, worldId, &session)
		respondSingle(c, &session, err)
	})

	sessionRouter.PUT("/:sessionId", func(c *gin.Context) {
		worldId, err := getIDParam(c, "worldId")
		if err != nil {
			respondBadRequest(c, "Invalid world ID")
			return
		}
		sessionId, err := getIDParam(c, "sessionId")
		if err != nil {
			respondBadRequest(c, "Invalid session ID")
			return
		}

		var session model.ChatSession
		if err := c.ShouldBindJSON(&session); err != nil {
			respondBadRequest(c, "Invalid session data")
			return
		}

		err = model.UpdateChatSession(db, worldId, sessionId, &session)
		respondSingle(c, &session, err)
	})

	sessionRouter.DELETE("/:sessionId", func(c *gin.Context) {
		worldId, err := getIDParam(c, "worldId")
		if err != nil {
			respondBadRequest(c, "Invalid world ID")
			return
		}
		sessionId, err := getIDParam(c, "sessionId")
		if err != nil {
			respondBadRequest(c, "Invalid session ID")
			return
		}

		err = model.DeleteChatSessionById(db, worldId, sessionId)
		respondDeleted(c, err)
	})

	sessionRouter.GET("/:sessionId/chat-messages", func(c *gin.Context) {
		sessionId, err := getIDParam(c, "sessionId")
		if err != nil {
			respondBadRequest(c, "Invalid session ID")
			return
		}

		messages, err := model.GetChatMessages(db, sessionId)
		respondList(c, messages, err)
	})

	sessionRouter.POST("/:sessionId/chat-messages", func(c *gin.Context) {
		sessionId, err := getIDParam(c, "sessionId")
		if err != nil {
			respondBadRequest(c, "Invalid session ID")
			return
		}

		var message model.ChatMessage
		if err := c.ShouldBindJSON(&message); err != nil {
			respondBadRequest(c, "Invalid chat message data")
			return
		}

		// TODO: Trigger response (embedding, fetching memories, prompt building, chat completion)
		//       go func...
		// TODO: Trigger chat truncation (creating memories)

		err = model.CreateChatMessage(db, sessionId, &message)
		respondSingle(c, &message, err)
	})

	sessionRouter.PUT("/:sessionId/chat-messages/:messageId", func(c *gin.Context) {
		sessionId, err := getIDParam(c, "sessionId")
		if err != nil {
			respondBadRequest(c, "Invalid session ID")
			return
		}

		messageId, err := getIDParam(c, "messageId")
		if err != nil {
			respondBadRequest(c, "Invalid chat message ID")
			return
		}

		var message model.ChatMessage
		if err := c.ShouldBindJSON(&message); err != nil {
			respondBadRequest(c, "Invalid chat message data")
			return
		}

		// TODO: Trigger updating of memories

		err = model.UpdateChatMessage(db, sessionId, messageId, &message)
		respondSingle(c, &message, err)
	})

	sessionRouter.DELETE("/:sessionId/chat-messages/:messageId", func(c *gin.Context) {
		sessionId, err := getIDParam(c, "sessionId")
		if err != nil {
			respondBadRequest(c, "Invalid session ID")
			return
		}
		messageId, err := getIDParam(c, "messageId")
		if err != nil {
			respondBadRequest(c, "Invalid chat message ID")
			return
		}

		err = model.DeleteChatMessagesFrom(db, sessionId, messageId)
		respondDeleted(c, err)
	})
}
