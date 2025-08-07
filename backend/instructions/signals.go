package instructions

import "github.com/maniartech/signals"

var InstructionTemplateCreatedSignal = signals.New[*InstructionTemplate]()
var InstructionTemplateUpdatedSignal = signals.New[*InstructionTemplate]()
var InstructionTemplateDeletedSignal = signals.New[int64]()
