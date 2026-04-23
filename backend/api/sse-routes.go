package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"juraji.nl/chat-quest/core/log"
	"juraji.nl/chat-quest/core/sse"
)

func SseRoutes(router *gin.RouterGroup) {
	sseRouter := router.Group("/sse")
	logger := log.Get()

	sseRouter.GET("", func(c *gin.Context) {
		clientIP := c.ClientIP()
		connectionId := fmt.Sprintf("SSE::%s::%s", clientIP, uuid.New())

		logger.Info("New SSE subscriber",
			zap.String("clientIP", clientIP),
			zap.String("connectionId", connectionId))

		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")

		ctx := c.Request.Context()
		clientChan := make(chan sse.Message)
		pingTicker := time.NewTicker(30 * time.Second)
		defer pingTicker.Stop()

		sse.SseCombinedSignal.AddListener(connectionId, func(_ context.Context, m sse.Message) error {
			clientChan <- m
			return nil
		})

		// Write initial message to confirm connection with connection ID
		if err := writeAndFlushEvent(c, "connection", fmt.Sprintf("SSE connected! Connection ID: %s", connectionId)); err != nil {
			logger.Error("failed to send 'SSE connected' event to client", zap.Error(err))
		}

		for {
			select {
			case <-ctx.Done():
				sse.SseCombinedSignal.RemoveListener(connectionId)
				close(clientChan)
				logger.Info("SSE subscriber left, connection closed", zap.String("connectionId", connectionId))
				return

			case msg := <-clientChan:
				j, err := json.Marshal(msg)
				if err != nil {
					logger.Error("failed to marshal event",
						zap.Any("msg", msg),
						zap.Error(err))
					continue
				}

				if err = writeAndFlushEvent(c, "message", string(j)); err != nil {
					logger.Error("failed to write message to client",
						zap.Any("msg", msg),
						zap.Error(err))
				}

			case <-pingTicker.C:
				// Ping for keep-alive on client side with connection ID
				timestamp := strconv.FormatInt(time.Now().Unix(), 10)
				if err := writeAndFlushEvent(c, "ping", timestamp); err != nil {
					logger.Error("failed to write ping to client", zap.Error(err))
				}
			}
		}
	})
}

func writeAndFlushEvent(c *gin.Context, event string, data string) error {
	c.Header("Retry-After", "5000")

	if _, err := fmt.Fprintf(c.Writer, "event: %s\ndata: %s\n\n", event, data); err != nil {
		return err
	}

	if f, ok := c.Writer.(http.Flusher); ok {
		f.Flush()
	}

	return nil
}
