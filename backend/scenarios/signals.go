package scenarios

import "github.com/maniartech/signals"

var ScenarioCreatedSignal = signals.New[*Scenario]()
var ScenarioUpdatedSignal = signals.New[*Scenario]()
var ScenarioDeletedSignal = signals.New[int64]()
