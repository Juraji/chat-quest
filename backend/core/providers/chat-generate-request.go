package providers

import "strings"

type LlmParameters struct {
	MaxTokens        int
	Temperature      float32
	TopP             float32
	PresencePenalty  float32
	FrequencyPenalty float32
	Stream           bool
	StopSequences    *string

	// An optional response format (JSON Schema)
	ResponseFormat *string
}

func (params *LlmParameters) StopSequencesAsSlice() []string {
	if params.StopSequences == nil {
		return nil
	}

	seqTrimmed := strings.TrimSpace(*params.StopSequences)
	if seqTrimmed == "" {
		return nil
	}

	var sequences []string
	sequences = strings.Split(seqTrimmed, ",")
	for i := range sequences {
		sequences[i] = strings.TrimSpace(sequences[i])
	}

	return sequences
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
	Content     string
	Error       error
	TotalTokens int
}
