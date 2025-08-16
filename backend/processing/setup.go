package processing

import (
	chatsessions "juraji.nl/chat-quest/model/chat-sessions"
	"juraji.nl/chat-quest/processing/chat_greeting"
)

func SetupProcessing() {
	chatsessions.ChatSessionCreatedSignal.AddListener(chat_greeting.CreateChatSessionGreetings)
}
