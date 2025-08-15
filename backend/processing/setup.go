package processing

import chatsessions "juraji.nl/chat-quest/model/chat-sessions"

func SetupProcessing() {
	chatsessions.ChatSessionCreatedSignal.AddListener(onNewChatSessionHandleCharacterGreetings)
}
