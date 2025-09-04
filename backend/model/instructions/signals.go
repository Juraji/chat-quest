package instructions

import (
	"juraji.nl/chat-quest/core/sse"
	"juraji.nl/chat-quest/core/util/signals"
)

var InstructionCreatedSignal = signals.New[*Instruction]()
var InstructionUpdatedSignal = signals.New[*Instruction]()
var InstructionDeletedSignal = signals.New[int]()

func init() {
	sse.RegisterOnSSE("InstructionCreated", InstructionCreatedSignal)
	sse.RegisterOnSSE("InstructionUpdated", InstructionUpdatedSignal)
	sse.RegisterOnSSE("InstructionDeleted", InstructionDeletedSignal)
}
