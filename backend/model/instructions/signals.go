package instructions

import (
	"github.com/maniartech/signals"
	"juraji.nl/chat-quest/core/sse"
)

var InstructionCreatedSignal = signals.New[*InstructionTemplate]()
var InstructionUpdatedSignal = signals.New[*InstructionTemplate]()
var InstructionDeletedSignal = signals.New[int]()

func init() {
	sse.RegisterOnSSE("InstructionCreated", InstructionCreatedSignal)
	sse.RegisterOnSSE("InstructionUpdated", InstructionUpdatedSignal)
	sse.RegisterOnSSE("InstructionDeleted", InstructionDeletedSignal)
}
