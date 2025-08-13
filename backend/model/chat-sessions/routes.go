package chat_sessions

import (
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/core"
	"juraji.nl/chat-quest/core/util"
)

func Routes(cq *core.ChatQuestContext, router *gin.RouterGroup) {
	sessionRouter := router.Group("/worlds/:worldId/chat-sessions")

	sessionRouter.GET("", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		worldId, err := util.GetIDParam(c, "worldId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid world ID")
			return
		}

		sessions, err := GetAllChatSessionsByWorldId(rcq, worldId)
		util.RespondList(rcq, c, sessions, err)
	})

	sessionRouter.GET("/:sessionId", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		worldId, err := util.GetIDParam(c, "worldId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid world ID")
			return
		}
		sessionId, err := util.GetIDParam(c, "sessionId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid session ID")
			return
		}

		session, err := GetChatSessionById(rcq, worldId, sessionId)
		util.RespondSingle(rcq, c, session, err)
	})

	sessionRouter.POST("", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		worldId, err := util.GetIDParam(c, "worldId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid world ID")
			return
		}

		characterIds, err := util.GetIDsFromQuery(c, "characterId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "invalid characterId in query")
			return
		}

		var session ChatSession
		if err := c.ShouldBindJSON(&session); err != nil {
			util.RespondBadRequest(rcq, c, "Invalid session data")
			return
		}

		err = CreateChatSession(rcq, worldId, &session, characterIds)
		util.RespondSingle(rcq, c, &session, err)
	})

	sessionRouter.PUT("/:sessionId", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		worldId, err := util.GetIDParam(c, "worldId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid world ID")
			return
		}
		sessionId, err := util.GetIDParam(c, "sessionId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid session ID")
			return
		}

		var session ChatSession
		if err := c.ShouldBindJSON(&session); err != nil {
			util.RespondBadRequest(rcq, c, "Invalid session data")
			return
		}

		err = UpdateChatSession(rcq, worldId, sessionId, &session)
		util.RespondSingle(rcq, c, &session, err)
	})

	sessionRouter.DELETE("/:sessionId", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		worldId, err := util.GetIDParam(c, "worldId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid world ID")
			return
		}
		sessionId, err := util.GetIDParam(c, "sessionId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid session ID")
			return
		}

		err = DeleteChatSessionById(rcq, worldId, sessionId)
		util.RespondDeleted(rcq, c, err)
	})

	sessionRouter.GET("/:sessionId/chat-messages", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		sessionId, err := util.GetIDParam(c, "sessionId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid session ID")
			return
		}

		messages, err := GetChatMessages(rcq, sessionId)
		util.RespondList(rcq, c, messages, err)
	})

	sessionRouter.POST("/:sessionId/chat-messages", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		sessionId, err := util.GetIDParam(c, "sessionId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid session ID")
			return
		}

		var message ChatMessage
		if err := c.ShouldBindJSON(&message); err != nil {
			util.RespondBadRequest(rcq, c, "Invalid chat message data")
			return
		}

		// TODO: Trigger response (embedding, fetching memories, prompt building, chat completion)
		//       go func...
		// TODO: Trigger chat truncation (creating memories)

		err = CreateChatMessage(rcq, sessionId, &message)
		util.RespondSingle(rcq, c, &message, err)
	})

	sessionRouter.PUT("/:sessionId/chat-messages/:messageId", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		sessionId, err := util.GetIDParam(c, "sessionId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid session ID")
			return
		}

		messageId, err := util.GetIDParam(c, "messageId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid chat message ID")
			return
		}

		var message ChatMessage
		if err := c.ShouldBindJSON(&message); err != nil {
			util.RespondBadRequest(rcq, c, "Invalid chat message data")
			return
		}

		// TODO: Trigger updating of memories

		err = UpdateChatMessage(rcq, sessionId, messageId, &message)
		util.RespondSingle(rcq, c, &message, err)
	})

	sessionRouter.DELETE("/:sessionId/chat-messages/:messageId", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		sessionId, err := util.GetIDParam(c, "sessionId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid session ID")
			return
		}
		messageId, err := util.GetIDParam(c, "messageId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid chat message ID")
			return
		}

		err = DeleteChatMessagesFrom(rcq, sessionId, messageId)
		util.RespondDeleted(rcq, c, err)
	})

	sessionRouter.GET("/:sessionId/participants", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		sessionId, err := util.GetIDParam(c, "sessionId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid session ID")
			return
		}

		participants, err := GetChatSessionParticipants(rcq, sessionId)
		util.RespondList(rcq, c, participants, err)
	})

	sessionRouter.POST("/:sessionId/participants/:characterId", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		sessionId, err := util.GetIDParam(c, "sessionId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid session ID")
			return
		}
		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid character ID")
			return
		}

		err = AddChatSessionParticipant(rcq, sessionId, characterId)
		util.RespondEmpty(rcq, c, err)
	})

	sessionRouter.DELETE("/:sessionId/participants/:characterId", func(c *gin.Context) {
		rcq := cq.WithContext(c.Request.Context())

		sessionId, err := util.GetIDParam(c, "sessionId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid session ID")
			return
		}
		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(rcq, c, "Invalid character ID")
			return
		}

		err = RemoveChatSessionParticipant(rcq, sessionId, characterId)
		util.RespondEmpty(rcq, c, err)
	})
}
