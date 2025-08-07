package worlds

import "github.com/maniartech/signals"

var WorldCreatedSignal = signals.New[*World]()
var WorldUpdatedSignal = signals.New[*World]()
var WorldDeletedSignal = signals.New[int64]()

var ChatPreferencesUpdatedSignal = signals.New[*ChatPreferences]()
