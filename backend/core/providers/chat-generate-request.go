package providers

type ChatGenerateRequest struct {
	Messages      []ChatRequestMessage
	ModelId       string
	MaxTokens     int
	Temperature   float32
	TopP          float32
	Stream        bool
	StopSequences []string
}

type ChatMessageRole string

const (
	RoleSystem    ChatMessageRole = "SYSTEM"
	RoleUser      ChatMessageRole = "USER"
	RoleAssistant ChatMessageRole = "ASSISTANT"
)

type ChatRequestMessage struct {
	Role    ChatMessageRole
	Content string
}

type ChatGenerateResponse struct {
	Content string
	Error   error
}
