package processing

import (
	chatsessions "juraji.nl/chat-quest/model/chat-sessions"
	"juraji.nl/chat-quest/processing/chat_init"
)

func SetupProcessing() {
	chatsessions.ChatSessionCreatedSignal.AddListener(chat_init.CreateChatSessionGreetings)
}
