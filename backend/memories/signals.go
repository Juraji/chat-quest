package memories

import "github.com/maniartech/signals"

var MemoryCreatedSignal = signals.New[*Memory]()
var MemoryUpdatedSignal = signals.New[*Memory]()
var MemoryDeletedSignal = signals.New[int64]()

var MemoryPreferencesUpdatedSignal = signals.New[*MemoryPreferences]()
