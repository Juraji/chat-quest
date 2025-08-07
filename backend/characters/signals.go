package characters

import "github.com/maniartech/signals"

var CharacterCreatedSignal = signals.New[*Character]()
var CharacterUpdatedSignal = signals.New[*Character]()
var CharacterDeletedSignal = signals.New[int64]()
