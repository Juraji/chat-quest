package chat_sessions

import (
	"strings"

	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/core/util/controllers"
)

func Routes(router *gin.RouterGroup) {
	sessionRouter := router.Group("/worlds/:worldId/chat-sessions")

	sessionRouter.GET("", func(c *gin.Context) {
		worldId, ok := controllers.GetParamAsID(c, "worldId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid world ID", nil)
			return
		}

		sessions, err := GetAllByWorldId(worldId)
		controllers.RespondList(c, sessions, err)
	})

	sessionRouter.GET("/:sessionId", func(c *gin.Context) {
		worldId, ok := controllers.GetParamAsID(c, "worldId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid world ID", nil)
			return
		}
		sessionId, ok := controllers.GetParamAsID(c, "sessionId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid session ID", nil)
			return
		}

		session, err := GetByWorldIdAndId(worldId, sessionId)
		controllers.RespondSingle(c, session, err)
	})

	sessionRouter.POST("", func(c *gin.Context) {
		worldId, ok := controllers.GetParamAsID(c, "worldId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid world ID", nil)
			return
		}

		characterIds, _ := controllers.GetQueryParamsAsIDs(c, "characterId")

		var session ChatSession
		if err := c.ShouldBindJSON(&session); err != nil {
			controllers.RespondBadRequest(c, "Invalid session data", nil)
			return
		}

		err := Create(worldId, &session, characterIds)
		controllers.RespondSingle(c, &session, err)
	})

	sessionRouter.PUT("/:sessionId", func(c *gin.Context) {
		worldId, ok := controllers.GetParamAsID(c, "worldId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid world ID", nil)
			return
		}
		sessionId, ok := controllers.GetParamAsID(c, "sessionId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid session ID", nil)
			return
		}

		var session ChatSession
		if err := c.ShouldBindJSON(&session); err != nil {
			controllers.RespondBadRequest(c, "Invalid session data", nil)
			return
		}

		err := Update(worldId, sessionId, &session)
		controllers.RespondSingle(c, &session, err)
	})

	sessionRouter.DELETE("/:sessionId", func(c *gin.Context) {
		worldId, ok := controllers.GetParamAsID(c, "worldId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid world ID", nil)
			return
		}
		sessionId, ok := controllers.GetParamAsID(c, "sessionId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid session ID", nil)
			return
		}

		err := Delete(worldId, sessionId)
		controllers.RespondEmpty(c, err)
	})

	sessionRouter.GET("/:sessionId/chat-messages", func(c *gin.Context) {
		sessionId, ok := controllers.GetParamAsID(c, "sessionId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid session ID", nil)
			return
		}

		messages, err := GetAllChatMessages(sessionId)
		controllers.RespondList(c, messages, err)
	})

	sessionRouter.POST("/:sessionId/chat-messages", func(c *gin.Context) {
		sessionId, ok := controllers.GetParamAsID(c, "sessionId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid session ID", nil)
			return
		}

		var message ChatMessage
		if err := c.ShouldBindJSON(&message); err != nil {
			controllers.RespondBadRequest(c, "Invalid chat message data", nil)
			return
		}

		err := CreateChatMessage(sessionId, &message)
		controllers.RespondSingle(c, &message, err)
	})

	sessionRouter.PUT("/:sessionId/chat-messages/:messageId", func(c *gin.Context) {
		sessionId, ok := controllers.GetParamAsID(c, "sessionId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid session ID", nil)
			return
		}

		messageId, ok := controllers.GetParamAsID(c, "messageId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid chat message ID", nil)
			return
		}

		var message ChatMessage
		if err := c.ShouldBindJSON(&message); err != nil {
			controllers.RespondBadRequest(c, "Invalid chat message data", nil)
			return
		}

		err := UpdateChatMessage(sessionId, messageId, &message)
		controllers.RespondSingle(c, &message, err)
	})

	sessionRouter.DELETE("/:sessionId/chat-messages/:messageId", func(c *gin.Context) {
		sessionId, ok := controllers.GetParamAsID(c, "sessionId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid session ID", nil)
			return
		}
		messageId, ok := controllers.GetParamAsID(c, "messageId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid chat message ID", nil)
			return
		}

		err := DeleteChatMessagesFrom(sessionId, messageId)
		controllers.RespondEmpty(c, err)
	})

	sessionRouter.POST("/:sessionId/chat-messages/:messageId/fork", func(c *gin.Context) {
		sessionId, ok := controllers.GetParamAsID(c, "sessionId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid session ID", nil)
			return
		}
		messageId, ok := controllers.GetParamAsID(c, "messageId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid chat message ID", nil)
			return
		}

		session, err := ForkChatSession(sessionId, messageId)
		controllers.RespondSingle(c, session, err)
	})

	sessionRouter.GET("/:sessionId/participants", func(c *gin.Context) {
		sessionId, ok := controllers.GetParamAsID(c, "sessionId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid session ID", nil)
			return
		}

		participants, err := GetAllParticipants(sessionId)
		controllers.RespondList(c, participants, err)
	})

	sessionRouter.POST("/:sessionId/participants/:characterId", func(c *gin.Context) {
		sessionId, ok := controllers.GetParamAsID(c, "sessionId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid session ID", nil)
			return
		}
		characterId, ok := controllers.GetParamAsID(c, "characterId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid character ID", nil)
			return
		}

		muted := strings.ToLower(c.Query("muted")) == "true"

		err := AddParticipant(sessionId, characterId, muted)
		controllers.RespondEmpty(c, err)
	})

	sessionRouter.DELETE("/:sessionId/participants/:characterId", func(c *gin.Context) {
		sessionId, ok := controllers.GetParamAsID(c, "sessionId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid session ID", nil)
			return
		}
		characterId, ok := controllers.GetParamAsID(c, "characterId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid character ID", nil)
			return
		}

		err := RemoveParticipant(sessionId, characterId)
		controllers.RespondEmpty(c, err)
	})

	sessionRouter.POST("/:sessionId/participants/:characterId/trigger-response", func(c *gin.Context) {
		sessionId, ok := controllers.GetParamAsID(c, "sessionId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid session ID", nil)
			return
		}
		characterId, ok := controllers.GetParamAsID(c, "characterId")
		if !ok {
			controllers.RespondBadRequest(c, "Invalid character ID", nil)
			return
		}

		participant, err := GetParticipantAsCharacter(sessionId, characterId)
		if err != nil {
			controllers.RespondBadRequest(c, "Character is not a participant", err)
			return
		}

		ChatParticipantResponseRequestedSignal.EmitBG(participant)
		controllers.RespondEmpty(c, nil)
	})
}
