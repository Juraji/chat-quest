package providers

import (
	"github.com/maniartech/signals"
	"juraji.nl/chat-quest/core/sse"
)

var ConnectionProfileCreatedSignal = signals.New[*ConnectionProfile]()
var ConnectionProfileUpdatedSignal = signals.New[*ConnectionProfile]()
var ConnectionProfileDeletedSignal = signals.New[int]()

var LlmModelCreatedSignal = signals.New[*LlmModel]()
var LlmModelUpdatedSignal = signals.New[*LlmModel]()
var LlmModelDeletedSignal = signals.New[int]()

func init() {
	sse.RegisterOnSSE("ConnectionProfileCreated", ConnectionProfileCreatedSignal)
	sse.RegisterOnSSE("ConnectionProfileUpdated", ConnectionProfileUpdatedSignal)
	sse.RegisterOnSSE("ConnectionProfileDeleted", ConnectionProfileDeletedSignal)
	sse.RegisterOnSSE("LlmModelCreated", LlmModelCreatedSignal)
	sse.RegisterOnSSE("LlmModelUpdated", LlmModelUpdatedSignal)
	sse.RegisterOnSSE("LlmModelDeleted", LlmModelDeletedSignal)
}
