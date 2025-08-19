package chat_sessions

import (
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/core/util"
)

func Routes(router *gin.RouterGroup) {
	sessionRouter := router.Group("/worlds/:worldId/chat-sessions")

	sessionRouter.GET("", func(c *gin.Context) {
		worldId, err := util.GetIDParam(c, "worldId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid world ID")
			return
		}

		sessions, err := GetAllByWorldId(worldId)
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

		session, err := GetByWorldIdAndId(worldId, sessionId)
		util.RespondSingle(c, session, err)
	})

	sessionRouter.POST("", func(c *gin.Context) {
		worldId, err := util.GetIDParam(c, "worldId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid world ID")
			return
		}

		characterIds, err := util.GetIDsFromQuery(c, "characterId")
		if err != nil {
			util.RespondBadRequest(c, "invalid characterId in query")
			return
		}

		var session ChatSession
		if err := c.ShouldBindJSON(&session); err != nil {
			util.RespondBadRequest(c, "Invalid session data")
			return
		}

		err = Create(worldId, &session, characterIds)
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

		var session ChatSession
		if err := c.ShouldBindJSON(&session); err != nil {
			util.RespondBadRequest(c, "Invalid session data")
			return
		}

		err = Update(worldId, sessionId, &session)
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

		err = Delete(worldId, sessionId)
		util.RespondDeleted(c, err)
	})

	sessionRouter.GET("/:sessionId/chat-messages", func(c *gin.Context) {
		sessionId, err := util.GetIDParam(c, "sessionId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid session ID")
			return
		}

		messages, err := GetChatMessages(sessionId)
		util.RespondList(c, messages, err)
	})

	sessionRouter.POST("/:sessionId/chat-messages", func(c *gin.Context) {
		sessionId, err := util.GetIDParam(c, "sessionId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid session ID")
			return
		}

		var message ChatMessage
		if err := c.ShouldBindJSON(&message); err != nil {
			util.RespondBadRequest(c, "Invalid chat message data")
			return
		}

		message.IsSystem = false

		err = CreateChatMessage(sessionId, &message)
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

		var message ChatMessage
		if err := c.ShouldBindJSON(&message); err != nil {
			util.RespondBadRequest(c, "Invalid chat message data")
			return
		}

		err = UpdateChatMessage(sessionId, messageId, &message)
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

		err = DeleteChatMessagesFrom(sessionId, messageId)
		util.RespondDeleted(c, err)
	})

	sessionRouter.GET("/:sessionId/participants", func(c *gin.Context) {
		sessionId, err := util.GetIDParam(c, "sessionId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid session ID")
			return
		}

		participants, err := GetParticipants(sessionId)
		util.RespondList(c, participants, err)
	})

	sessionRouter.POST("/:sessionId/participants/:characterId", func(c *gin.Context) {
		sessionId, err := util.GetIDParam(c, "sessionId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid session ID")
			return
		}
		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid character ID")
			return
		}

		err = AddParticipant(sessionId, characterId)
		util.RespondEmpty(c, err)
	})

	sessionRouter.DELETE("/:sessionId/participants/:characterId", func(c *gin.Context) {
		sessionId, err := util.GetIDParam(c, "sessionId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid session ID")
			return
		}
		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid character ID")
			return
		}

		err = RemoveParticipant(sessionId, characterId)
		util.RespondEmpty(c, err)
	})

	sessionRouter.POST("/:sessionId/participants/:characterId/trigger-response", func(c *gin.Context) {
		sessionId, err := util.GetIDParam(c, "sessionId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid session ID")
			return
		}
		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(c, "Invalid character ID")
			return
		}

		participant, err := GetParticipant(sessionId, characterId)
		if err != nil {
			util.RespondBadRequest(c, "Character is not a participant")
			return
		}

		ChatParticipantResponseRequestedSignal.EmitBG(participant)
		util.RespondEmpty(c, nil)
	})
}
