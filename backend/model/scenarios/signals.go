package scenarios

import (
	"juraji.nl/chat-quest/core/sse"
	"juraji.nl/chat-quest/core/util/signals"
)

var ScenarioCreatedSignal = signals.New[*Scenario]()
var ScenarioUpdatedSignal = signals.New[*Scenario]()
var ScenarioDeletedSignal = signals.New[int]()

func init() {
	sse.RegisterOnSSE("ScenarioCreated", ScenarioCreatedSignal)
	sse.RegisterOnSSE("ScenarioUpdated", ScenarioUpdatedSignal)
	sse.RegisterOnSSE("ScenarioDeleted", ScenarioDeletedSignal)
}
