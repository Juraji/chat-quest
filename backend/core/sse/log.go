package sse

import "juraji.nl/chat-quest/core/log"

func init() {
	RegisterOnSSE("LogMessages", log.LogMessagesSignal)
}
