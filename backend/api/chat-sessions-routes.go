package api

import (
	"strings"

	"github.com/gin-gonic/gin"
	cs "juraji.nl/chat-quest/model/chat-sessions"
	"juraji.nl/chat-quest/processing"
)

func ChatSessionsRoutes(router *gin.RouterGroup) {
	sessionRouter := router.Group("/worlds/:worldId/chat-sessions")

	sessionRouter.GET("", func(c *gin.Context) {
		worldId, ok := getParamAsID(c, "worldId")
		if !ok {
			respondBadRequest(c, "Invalid world ID", nil)
			return
		}

		sessions, err := cs.GetAllByWorldId(worldId)
		respondList(c, sessions, err)
	})

	sessionRouter.GET("/:sessionId", func(c *gin.Context) {
		worldId, ok := getParamAsID(c, "worldId")
		if !ok {
			respondBadRequest(c, "Invalid world ID", nil)
			return
		}
		sessionId, ok := getParamAsID(c, "sessionId")
		if !ok {
			respondBadRequest(c, "Invalid session ID", nil)
			return
		}

		session, err := cs.GetByWorldIdAndId(worldId, sessionId)
		respondSingle(c, session, err)
	})

	sessionRouter.POST("", func(c *gin.Context) {
		worldId, ok := getParamAsID(c, "worldId")
		if !ok {
			respondBadRequest(c, "Invalid world ID", nil)
			return
		}

		characterIds, _ := getQueryParamsAsInts(c, "characterId")

		var session cs.ChatSession
		if err := c.ShouldBindJSON(&session); err != nil || !session.CurrentTimeOfDay.IsValid() {
			respondBadRequest(c, "Invalid session data", nil)
			return
		}

		err := cs.Create(worldId, &session, characterIds)
		respondSingle(c, &session, err)
	})

	sessionRouter.PUT("/:sessionId", func(c *gin.Context) {
		worldId, ok := getParamAsID(c, "worldId")
		if !ok {
			respondBadRequest(c, "Invalid world ID", nil)
			return
		}
		sessionId, ok := getParamAsID(c, "sessionId")
		if !ok {
			respondBadRequest(c, "Invalid session ID", nil)
			return
		}

		var session cs.ChatSession
		if err := c.ShouldBindJSON(&session); err != nil {
			respondBadRequest(c, "Invalid session data", nil)
			return
		}

		err := cs.Update(worldId, sessionId, &session)
		respondSingle(c, &session, err)
	})

	sessionRouter.DELETE("/:sessionId", func(c *gin.Context) {
		worldId, ok := getParamAsID(c, "worldId")
		if !ok {
			respondBadRequest(c, "Invalid world ID", nil)
			return
		}
		sessionId, ok := getParamAsID(c, "sessionId")
		if !ok {
			respondBadRequest(c, "Invalid session ID", nil)
			return
		}

		err := cs.Delete(worldId, sessionId)
		respondEmpty(c, err)
	})

	sessionRouter.GET("/:sessionId/chat-messages", func(c *gin.Context) {
		sessionId, ok := getParamAsID(c, "sessionId")
		if !ok {
			respondBadRequest(c, "Invalid session ID", nil)
			return
		}

		messages, err := cs.GetAllChatMessages(sessionId)
		respondList(c, messages, err)
	})

	sessionRouter.POST("/:sessionId/chat-messages", func(c *gin.Context) {
		sessionId, ok := getParamAsID(c, "sessionId")
		if !ok {
			respondBadRequest(c, "Invalid session ID", nil)
			return
		}

		var message cs.ChatMessage
		if err := c.ShouldBindJSON(&message); err != nil {
			respondBadRequest(c, "Invalid chat message data", nil)
			return
		}

		err := cs.CreateChatMessage(sessionId, &message)
		respondSingle(c, &message, err)
	})

	sessionRouter.PUT("/:sessionId/chat-messages/:messageId", func(c *gin.Context) {
		sessionId, ok := getParamAsID(c, "sessionId")
		if !ok {
			respondBadRequest(c, "Invalid session ID", nil)
			return
		}

		messageId, ok := getParamAsID(c, "messageId")
		if !ok {
			respondBadRequest(c, "Invalid chat message ID", nil)
			return
		}

		var message cs.ChatMessage
		if err := c.ShouldBindJSON(&message); err != nil {
			respondBadRequest(c, "Invalid chat message data", nil)
			return
		}

		err := cs.UpdateChatMessage(sessionId, messageId, &message)
		respondSingle(c, &message, err)
	})

	sessionRouter.DELETE("/:sessionId/chat-messages/:messageId", func(c *gin.Context) {
		sessionId, ok := getParamAsID(c, "sessionId")
		if !ok {
			respondBadRequest(c, "Invalid session ID", nil)
			return
		}
		messageId, ok := getParamAsID(c, "messageId")
		if !ok {
			respondBadRequest(c, "Invalid chat message ID", nil)
			return
		}

		err := cs.DeleteChatMessagesFrom(sessionId, messageId)
		respondEmpty(c, err)
	})

	sessionRouter.POST("/:sessionId/chat-messages/:messageId/fork", func(c *gin.Context) {
		sessionId, ok := getParamAsID(c, "sessionId")
		if !ok {
			respondBadRequest(c, "Invalid session ID", nil)
			return
		}
		messageId, ok := getParamAsID(c, "messageId")
		if !ok {
			respondBadRequest(c, "Invalid chat message ID", nil)
			return
		}

		session, err := cs.ForkChatSession(sessionId, messageId)
		respondSingle(c, session, err)
	})

	sessionRouter.GET("/:sessionId/participants", func(c *gin.Context) {
		sessionId, ok := getParamAsID(c, "sessionId")
		if !ok {
			respondBadRequest(c, "Invalid session ID", nil)
			return
		}

		participants, err := cs.GetAllParticipants(sessionId)
		respondList(c, participants, err)
	})

	sessionRouter.POST("/:sessionId/participants/:characterId", func(c *gin.Context) {
		sessionId, ok := getParamAsID(c, "sessionId")
		if !ok {
			respondBadRequest(c, "Invalid session ID", nil)
			return
		}
		characterId, ok := getParamAsID(c, "characterId")
		if !ok {
			respondBadRequest(c, "Invalid character ID", nil)
			return
		}

		muted := strings.ToLower(c.Query("muted")) == "true"

		err := cs.AddParticipant(sessionId, characterId, muted)
		respondEmpty(c, err)
	})

	sessionRouter.DELETE("/:sessionId/participants/:characterId", func(c *gin.Context) {
		sessionId, ok := getParamAsID(c, "sessionId")
		if !ok {
			respondBadRequest(c, "Invalid session ID", nil)
			return
		}
		characterId, ok := getParamAsID(c, "characterId")
		if !ok {
			respondBadRequest(c, "Invalid character ID", nil)
			return
		}

		err := cs.RemoveParticipant(sessionId, characterId)
		respondEmpty(c, err)
	})

	sessionRouter.POST("/:sessionId/participants/:characterId/trigger-response", func(c *gin.Context) {
		sessionId, ok := getParamAsID(c, "sessionId")
		if !ok {
			respondBadRequest(c, "Invalid session ID", nil)
			return
		}
		characterId, ok := getParamAsID(c, "characterId")
		if !ok {
			respondBadRequest(c, "Invalid character ID", nil)
			return
		}

		participant, err := cs.GetParticipantAsCharacter(sessionId, characterId)
		if err != nil {
			respondBadRequest(c, "Character is not a participant", err)
			return
		}

		// TODO: Handle errors
		processing.GenerateResponseByParticipantTrigger(c, participant)
		respondEmpty(c, nil)
	})

	sessionRouter.POST("/:sessionId/generate-title", func(c *gin.Context) {
		sessionId, ok := getParamAsID(c, "sessionId")
		if !ok {
			respondBadRequest(c, "Invalid session ID", nil)
			return
		}

		// TODO: Handle errors
		processing.GenerateTitle(c, sessionId)
		respondEmpty(c, nil)
	})
}
