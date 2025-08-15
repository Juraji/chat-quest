package sse

import (
  "context"
  "encoding/json"
  "fmt"
  "github.com/gin-gonic/gin"
  "github.com/google/uuid"
  "go.uber.org/zap"
  "juraji.nl/chat-quest/core/log"
  "net/http"
  "strconv"
  "time"
)

func Routes(router *gin.RouterGroup) {
  sseRouter := router.Group("/sse")

  sseRouter.GET("", func(c *gin.Context) {
    connectionId := fmt.Sprintf("SSE::%s::%s", c.ClientIP(), uuid.New())

    c.Header("Content-Type", "text/event-stream")
    c.Header("Cache-Control", "no-cache")
    c.Header("Connection", "keep-alive")

    log.Get().Info("New SSE subscriber",
      zap.String("connectionId", connectionId))

    ctx := c.Request.Context()
    clientChan := make(chan message)
    pingTicker := time.NewTicker(30 * time.Second)
    defer pingTicker.Stop()

    sseCombinedSignal.AddListener(func(_ context.Context, m message) {
      clientChan <- m
    }, connectionId)

    // Write initial message to confirm connection with connection ID
    if err := writeAndFlushEvent(c, "connection", fmt.Sprintf("SSE connected! Connection ID: %s", connectionId)); err != nil {
      log.Get().Error("failed to send 'SSE connected' event to client",
        zap.Error(err),
        zap.String("connectionId", connectionId))
    }

    for {
      select {
      case <-ctx.Done():
        sseCombinedSignal.RemoveListener(connectionId)
        close(clientChan)
        log.Get().Info("SSE subscriber left, connection closed", zap.String("connectionId", connectionId))
        return

      case msg := <-clientChan:
        j, err := json.Marshal(msg)
        if err != nil {
          log.Get().Error("failed to marshal event",
            zap.Error(err),
            zap.String("connectionId", connectionId),
            zap.Any("msg", msg))
          continue
        }

        if err = writeAndFlushEvent(c, "message", string(j)); err != nil {
          log.Get().Error("failed to write message to client",
            zap.Error(err),
            zap.String("connectionId", connectionId),
            zap.Any("msg", msg))
        }

      case <-pingTicker.C:
        // Ping for keep-alive on client side with connection ID
        timestamp := strconv.FormatInt(time.Now().Unix(), 10)
        if err := writeAndFlushEvent(c, "ping", timestamp); err != nil {
          log.Get().Error("failed to write ping to client",
            zap.Error(err),
            zap.String("connectionId", connectionId))
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
