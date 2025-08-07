package providers

import "github.com/maniartech/signals"

var ConnectionProfileCreatedSignal = signals.New[*ConnectionProfile]()
var ConnectionProfileUpdatedSignal = signals.New[*ConnectionProfile]()
var ConnectionProfileDeletedSignal = signals.New[int64]()

var LlmModelCreatedSignal = signals.New[*LlmModel]()
var LlmModelUpdatedSignal = signals.New[*LlmModel]()
var LlmModelDeletedSignal = signals.New[int64]()
