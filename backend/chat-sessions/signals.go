package chat_sessions

import "github.com/maniartech/signals"

var ChatSessionCreatedSignal = signals.New[*ChatSession]()
var ChatSessionUpdatedSignal = signals.New[*ChatSession]()
var ChatSessionDeletedSignal = signals.New[int64]()

var ChatMessageCreatedSignal = signals.New[*ChatMessage]()
var ChatMessageUpdatedSignal = signals.New[*ChatMessage]()
var ChatMessageDeletedSignal = signals.New[int64]()
