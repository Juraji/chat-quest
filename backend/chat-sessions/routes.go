package chat_sessions

import (
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/cq"
	"juraji.nl/chat-quest/util"
)

func Routes(cq *cq.ChatQuestContext, router *gin.RouterGroup) {
	sessionRouter := router.Group("/worlds/:worldId/chat-sessions")

	sessionRouter.GET("", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		worldId, err := util.GetIDParam(c, "worldId")
		if err != nil {
			util.RespondBadRequest(cq, c, "Invalid world ID")
			return
		}

		sessions, err := GetAllChatSessionsByWorldId(cq, worldId)
		util.RespondList(cq, c, sessions, err)
	})

	sessionRouter.GET("/:sessionId", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		worldId, err := util.GetIDParam(c, "worldId")
		if err != nil {
			util.RespondBadRequest(cq, c, "Invalid world ID")
			return
		}
		sessionId, err := util.GetIDParam(c, "sessionId")
		if err != nil {
			util.RespondBadRequest(cq, c, "Invalid session ID")
			return
		}

		session, err := GetChatSessionById(cq, worldId, sessionId)
		util.RespondSingle(cq, c, session, err)
	})

	sessionRouter.POST("", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		worldId, err := util.GetIDParam(c, "worldId")
		if err != nil {
			util.RespondBadRequest(cq, c, "Invalid world ID")
			return
		}

		characterIds, err := util.GetIDsFromQuery(c, "characterId")
		if err != nil {
			util.RespondBadRequest(cq, c, "invalid characterId in query")
			return
		}

		var session ChatSession
		if err := c.ShouldBindJSON(&session); err != nil {
			util.RespondBadRequest(cq, c, "Invalid session data")
			return
		}

		err = CreateChatSession(cq, worldId, &session, characterIds)
		util.RespondSingle(cq, c, &session, err)
	})

	sessionRouter.PUT("/:sessionId", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		worldId, err := util.GetIDParam(c, "worldId")
		if err != nil {
			util.RespondBadRequest(cq, c, "Invalid world ID")
			return
		}
		sessionId, err := util.GetIDParam(c, "sessionId")
		if err != nil {
			util.RespondBadRequest(cq, c, "Invalid session ID")
			return
		}

		var session ChatSession
		if err := c.ShouldBindJSON(&session); err != nil {
			util.RespondBadRequest(cq, c, "Invalid session data")
			return
		}

		err = UpdateChatSession(cq, worldId, sessionId, &session)
		util.RespondSingle(cq, c, &session, err)
	})

	sessionRouter.DELETE("/:sessionId", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		worldId, err := util.GetIDParam(c, "worldId")
		if err != nil {
			util.RespondBadRequest(cq, c, "Invalid world ID")
			return
		}
		sessionId, err := util.GetIDParam(c, "sessionId")
		if err != nil {
			util.RespondBadRequest(cq, c, "Invalid session ID")
			return
		}

		err = DeleteChatSessionById(cq, worldId, sessionId)
		util.RespondDeleted(cq, c, err)
	})

	sessionRouter.GET("/:sessionId/chat-messages", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		sessionId, err := util.GetIDParam(c, "sessionId")
		if err != nil {
			util.RespondBadRequest(cq, c, "Invalid session ID")
			return
		}

		messages, err := GetChatMessages(cq, sessionId)
		util.RespondList(cq, c, messages, err)
	})

	sessionRouter.POST("/:sessionId/chat-messages", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		sessionId, err := util.GetIDParam(c, "sessionId")
		if err != nil {
			util.RespondBadRequest(cq, c, "Invalid session ID")
			return
		}

		var message ChatMessage
		if err := c.ShouldBindJSON(&message); err != nil {
			util.RespondBadRequest(cq, c, "Invalid chat message data")
			return
		}

		// TODO: Trigger response (embedding, fetching memories, prompt building, chat completion)
		//       go func...
		// TODO: Trigger chat truncation (creating memories)

		err = CreateChatMessage(cq, sessionId, &message)
		util.RespondSingle(cq, c, &message, err)
	})

	sessionRouter.PUT("/:sessionId/chat-messages/:messageId", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		sessionId, err := util.GetIDParam(c, "sessionId")
		if err != nil {
			util.RespondBadRequest(cq, c, "Invalid session ID")
			return
		}

		messageId, err := util.GetIDParam(c, "messageId")
		if err != nil {
			util.RespondBadRequest(cq, c, "Invalid chat message ID")
			return
		}

		var message ChatMessage
		if err := c.ShouldBindJSON(&message); err != nil {
			util.RespondBadRequest(cq, c, "Invalid chat message data")
			return
		}

		// TODO: Trigger updating of memories

		err = UpdateChatMessage(cq, sessionId, messageId, &message)
		util.RespondSingle(cq, c, &message, err)
	})

	sessionRouter.DELETE("/:sessionId/chat-messages/:messageId", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		sessionId, err := util.GetIDParam(c, "sessionId")
		if err != nil {
			util.RespondBadRequest(cq, c, "Invalid session ID")
			return
		}
		messageId, err := util.GetIDParam(c, "messageId")
		if err != nil {
			util.RespondBadRequest(cq, c, "Invalid chat message ID")
			return
		}

		err = DeleteChatMessagesFrom(cq, sessionId, messageId)
		util.RespondDeleted(cq, c, err)
	})

	sessionRouter.GET("/:sessionId/participants", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		sessionId, err := util.GetIDParam(c, "sessionId")
		if err != nil {
			util.RespondBadRequest(cq, c, "Invalid session ID")
			return
		}

		participants, err := GetChatSessionParticipants(cq, sessionId)
		util.RespondList(cq, c, participants, err)
	})

	sessionRouter.POST("/:sessionId/participants/:characterId", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		sessionId, err := util.GetIDParam(c, "sessionId")
		if err != nil {
			util.RespondBadRequest(cq, c, "Invalid session ID")
			return
		}
		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(cq, c, "Invalid character ID")
			return
		}

		err = AddChatSessionParticipant(cq, sessionId, characterId)
		util.RespondEmpty(cq, c, err)
	})

	sessionRouter.DELETE("/:sessionId/participants/:characterId", func(c *gin.Context) {
		cq = cq.WithContext(c.Request.Context())

		sessionId, err := util.GetIDParam(c, "sessionId")
		if err != nil {
			util.RespondBadRequest(cq, c, "Invalid session ID")
			return
		}
		characterId, err := util.GetIDParam(c, "characterId")
		if err != nil {
			util.RespondBadRequest(cq, c, "Invalid character ID")
			return
		}

		err = RemoveChatSessionParticipant(cq, sessionId, characterId)
		util.RespondEmpty(cq, c, err)
	})
}
