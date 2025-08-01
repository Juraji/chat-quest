package ai

import "juraji.nl/chat-quest/model"

type Provider interface {
	getAvailableModels() ([]*model.LlmModel, error)
	generateChatCompletions(llmModel model.LlmModel, messages []Message) (string, error)
}

type Message struct {
	Role    string
	Content string
}
