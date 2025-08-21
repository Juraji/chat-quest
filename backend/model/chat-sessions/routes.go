package chat_sessions

import (
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/core/util/controllers"
)

func Routes(router *gin.RouterGroup) {
	sessionRouter := router.Group("/worlds/:worldId/chat-sessions")

	sessionRouter.GET("", func(c *gin.Context) {
		worldId, ok := controllers.GetParamAsID(c, "worldId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid world ID")
			return
		}

		sessions, ok := GetAllByWorldId(worldId)
		controllers.RespondList(c, ok, sessions)
	})

	sessionRouter.GET("/:sessionId", func(c *gin.Context) {
		worldId, ok := controllers.GetParamAsID(c, "worldId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid world ID")
			return
		}
		sessionId, ok := controllers.GetParamAsID(c, "sessionId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid session ID")
			return
		}

		session, ok := GetByWorldIdAndId(worldId, sessionId)
		controllers.RespondSingle(c, ok, session)
	})

	sessionRouter.POST("", func(c *gin.Context) {
		worldId, ok := controllers.GetParamAsID(c, "worldId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid world ID")
			return
		}

		characterIds, ok := controllers.GetQueryParamsAsIDs(c, "characterId")
		if !ok {
			controllers.RespondBadRequest(c, "invalid characterId in query")
			return
		}

		var session ChatSession
		if err := c.ShouldBindJSON(&session); err != nil {
			controllers.RespondBadRequest(c, "Invalid session data")
			return
		}

		ok = Create(worldId, &session, characterIds)
		controllers.RespondSingle(c, ok, &session)
	})

	sessionRouter.PUT("/:sessionId", func(c *gin.Context) {
		worldId, ok := controllers.GetParamAsID(c, "worldId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid world ID")
			return
		}
		sessionId, ok := controllers.GetParamAsID(c, "sessionId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid session ID")
			return
		}

		var session ChatSession
		if err := c.ShouldBindJSON(&session); err != nil {
			controllers.RespondBadRequest(c, "Invalid session data")
			return
		}

		ok = Update(worldId, sessionId, &session)
		controllers.RespondSingle(c, ok, &session)
	})

	sessionRouter.DELETE("/:sessionId", func(c *gin.Context) {
		worldId, ok := controllers.GetParamAsID(c, "worldId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid world ID")
			return
		}
		sessionId, ok := controllers.GetParamAsID(c, "sessionId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid session ID")
			return
		}

		ok = Delete(worldId, sessionId)
		controllers.RespondEmpty(c, ok)
	})

	sessionRouter.GET("/:sessionId/chat-messages", func(c *gin.Context) {
		sessionId, ok := controllers.GetParamAsID(c, "sessionId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid session ID")
			return
		}

		messages, ok := GetChatMessages(sessionId)
		controllers.RespondList(c, ok, messages)
	})

	sessionRouter.POST("/:sessionId/chat-messages", func(c *gin.Context) {
		sessionId, ok := controllers.GetParamAsID(c, "sessionId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid session ID")
			return
		}

		var message ChatMessage
		if err := c.ShouldBindJSON(&message); err != nil {
			controllers.RespondBadRequest(c, "Invalid chat message data")
			return
		}

		message.IsSystem = false

		ok = CreateChatMessage(sessionId, &message)
		controllers.RespondSingle(c, ok, &message)
	})

	sessionRouter.PUT("/:sessionId/chat-messages/:messageId", func(c *gin.Context) {
		sessionId, ok := controllers.GetParamAsID(c, "sessionId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid session ID")
			return
		}

		messageId, ok := controllers.GetParamAsID(c, "messageId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid chat message ID")
			return
		}

		var message ChatMessage
		if err := c.ShouldBindJSON(&message); err != nil {
			controllers.RespondBadRequest(c, "Invalid chat message data")
			return
		}

		ok = UpdateChatMessage(sessionId, messageId, &message)
		controllers.RespondSingle(c, ok, &message)
	})

	sessionRouter.DELETE("/:sessionId/chat-messages/:messageId", func(c *gin.Context) {
		sessionId, ok := controllers.GetParamAsID(c, "sessionId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid session ID")
			return
		}
		messageId, ok := controllers.GetParamAsID(c, "messageId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid chat message ID")
			return
		}

		ok = DeleteChatMessagesFrom(sessionId, messageId)
		controllers.RespondEmpty(c, ok)
	})

	sessionRouter.GET("/:sessionId/participants", func(c *gin.Context) {
		sessionId, ok := controllers.GetParamAsID(c, "sessionId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid session ID")
			return
		}

		participants, ok := GetParticipants(sessionId)
		controllers.RespondList(c, ok, participants)
	})

	sessionRouter.POST("/:sessionId/participants/:characterId", func(c *gin.Context) {
		sessionId, ok := controllers.GetParamAsID(c, "sessionId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid session ID")
			return
		}
		characterId, ok := controllers.GetParamAsID(c, "characterId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid character ID")
			return
		}

		ok = AddParticipant(sessionId, characterId)
		controllers.RespondEmpty(c, ok)
	})

	sessionRouter.DELETE("/:sessionId/participants/:characterId", func(c *gin.Context) {
		sessionId, ok := controllers.GetParamAsID(c, "sessionId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid session ID")
			return
		}
		characterId, ok := controllers.GetParamAsID(c, "characterId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid character ID")
			return
		}

		ok = RemoveParticipant(sessionId, characterId)
		controllers.RespondEmpty(c, ok)
	})

	sessionRouter.POST("/:sessionId/participants/:characterId/trigger-response", func(c *gin.Context) {
		sessionId, ok := controllers.GetParamAsID(c, "sessionId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid session ID")
			return
		}
		characterId, ok := controllers.GetParamAsID(c, "characterId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid character ID")
			return
		}

		participant, ok := GetParticipant(sessionId, characterId)
		if !ok {
			controllers.RespondBadRequest(c, "Character is not a participant")
			return
		}

		ChatParticipantResponseRequestedSignal.EmitBG(participant)
		controllers.RespondEmpty(c, true)
	})
}
