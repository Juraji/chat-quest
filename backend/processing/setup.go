package processing

import (
	chatsessions "juraji.nl/chat-quest/model/chat-sessions"
	"juraji.nl/chat-quest/processing/chat_init"
	"juraji.nl/chat-quest/processing/chat_response"
)

func SetupProcessing() {
	chatsessions.ChatSessionCreatedSignal.AddListener(
		"CreateChatSessionGreetings", chat_init.CreateChatSessionGreetings)
	chatsessions.ChatMessageCreatedSignal.AddListener(
		"GenerateChatSessionCharacterResponse", chat_response.GenerateChatSessionCharacterResponse)
}
