package sse

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"juraji.nl/chat-quest/cq"
	"net/http"
	"strconv"
	"time"
)

func Routes(cq *cq.ChatQuestContext, router *gin.RouterGroup) {
	sseRouter := router.Group("/sse")

	sseRouter.GET("", func(c *gin.Context) {
		clientIp := c.ClientIP()

		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")

		ctx := c.Request.Context()
		clientChan := make(chan messageBody, 10)
		pingTicker := time.NewTicker(30 * time.Second)
		defer pingTicker.Stop()

		listenerKeys := make(map[string]string)

		for _, source := range sseSourceSignals {
			sourceName := source.sourceName
			signal := source.signal
			key := "SSE_" + clientIp + "_" + sourceName

			listenerKeys[sourceName] = key
			signal.AddListener(func(ctx context.Context, payload any) {
				clientChan <- messageBody{sourceName, payload}
			}, key)
		}

		// Write initial message to confirm connection
		if err := writeAndFlushEvent(c, "connection", "SSE connected!"); err != nil {
			cq.Logger().Error("failed to send 'SSE connected' event to client", zap.Error(err), zap.String("clientIp", clientIp))
		}

		for {
			select {
			case <-ctx.Done():
				for _, source := range sseSourceSignals {
					source.signal.RemoveListener(listenerKeys[source.sourceName])
				}
				close(clientChan)
				return

			case msg := <-clientChan:
				j, err := json.Marshal(msg)
				if err != nil {
					cq.Logger().Error("failed to marshal event", zap.Error(err), zap.String("clientIp", clientIp))
					continue
				}

				if err = writeAndFlushEvent(c, "message", string(j)); err != nil {
					cq.Logger().Error("failed to write message to client", zap.Error(err), zap.String("clientIp", clientIp))
				}

			case <-pingTicker.C:
				// Ping for keep-alive on client side
				timestamp := strconv.FormatInt(time.Now().Unix(), 10)
				if err := writeAndFlushEvent(c, "ping", timestamp); err != nil {
					cq.Logger().Error("failed to write ping to client", zap.Error(err), zap.String("clientIp", clientIp))
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
