package providers

import (
	"github.com/maniartech/signals"
	"juraji.nl/chat-quest/sse"
)

var ConnectionProfileCreatedSignal = signals.New[*ConnectionProfile]()
var ConnectionProfileUpdatedSignal = signals.New[*ConnectionProfile]()
var ConnectionProfileDeletedSignal = signals.New[int64]()

var LlmModelCreatedSignal = signals.New[*LlmModel]()
var LlmModelUpdatedSignal = signals.New[*LlmModel]()
var LlmModelDeletedSignal = signals.New[int64]()

func init() {
	sse.RegisterSseSourceSignal("ConnectionProfileCreated", ConnectionProfileCreatedSignal)
	sse.RegisterSseSourceSignal("ConnectionProfileUpdated", ConnectionProfileUpdatedSignal)
	sse.RegisterSseSourceSignal("ConnectionProfileDeleted", ConnectionProfileDeletedSignal)
	sse.RegisterSseSourceSignal("LlmModelCreated", LlmModelCreatedSignal)
	sse.RegisterSseSourceSignal("LlmModelUpdated", LlmModelUpdatedSignal)
	sse.RegisterSseSourceSignal("LlmModelDeleted", LlmModelDeletedSignal)
}
