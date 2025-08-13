package scenarios

import (
	"github.com/maniartech/signals"
	"juraji.nl/chat-quest/core/sse"
)

var ScenarioCreatedSignal = signals.New[*Scenario]()
var ScenarioUpdatedSignal = signals.New[*Scenario]()
var ScenarioDeletedSignal = signals.New[int]()

func init() {
	sse.RegisterSseSourceSignal("ScenarioCreated", ScenarioCreatedSignal)
	sse.RegisterSseSourceSignal("ScenarioUpdated", ScenarioUpdatedSignal)
	sse.RegisterSseSourceSignal("ScenarioDeleted", ScenarioDeletedSignal)
}
