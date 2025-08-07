package sse

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"juraji.nl/chat-quest/characters"
	"log"
	"net/http"
	"strconv"
	"time"
)

var sseSourceSignals = []sourceSignal{
	newSourceSignal("CHARACTER_CREATED", characters.CharacterCreatedSignal),
	newSourceSignal("CHARACTER_UPDATED", characters.CharacterUpdatedSignal),
	newSourceSignal("CHARACTER_DELETED", characters.CharacterDeletedSignal),
}

func sseHandler(c *gin.Context) {
	clientIp := c.ClientIP()

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	ctx := c.Request.Context()
	clientChan := make(chan messageBody, 10)
	pingTicker := time.NewTicker(30 * time.Second)
	defer pingTicker.Stop()

	// Create a map to store listener keys for cleanup
	listenerKeys := make(map[string]string)

	for _, source := range sseSourceSignals {
		sourceName := source.sourceName
		signal := source.signal
		key := "SSE_" + clientIp + "_" + sourceName

		// Store the key for later removal
		listenerKeys[sourceName] = key

		// Create a listener function that handles any payload type
		signal.AddListener(func(ctx context.Context, payload any) {
			clientChan <- messageBody{sourceName, payload}
		}, key)
	}

	for {
		select {
		case <-ctx.Done():
			// Remove all listeners before closing
			for _, source := range sseSourceSignals {
				source.signal.RemoveListener(listenerKeys[source.sourceName])
			}
			close(clientChan)
			return

		case msg := <-clientChan:
			j, err := json.Marshal(msg)
			if err != nil {
				log.Printf("json marshal error on SSE for %v: %v", clientIp, err)
				continue
			}

			if err = writeAndFlushEvent(c, "message", string(j)); err != nil {
				log.Printf("error writing data on SSE for %v: %v", clientIp, err)
			}

		case <-pingTicker.C:
			// Ping for keep-alive on client side
			timestamp := strconv.FormatInt(time.Now().Unix(), 10)
			if err := writeAndFlushEvent(c, "ping", timestamp); err != nil {
				log.Printf("error writing data on SSE for %v: %v", clientIp, err)
			}
		}
	}
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
