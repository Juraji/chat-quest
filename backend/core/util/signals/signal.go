package signals

import (
	"context"
	"fmt"
	"juraji.nl/chat-quest/core/util"
	"sync"
)

// SignalListener defines a callback function type for handling signal events.
// Use this to process payloads when a signal is emitted.
type SignalListener[T any] func(ctx context.Context, payload T)

// subscriber represents an internal structure used to manage subscribers.
// You typically won't interact with this directly.
type subscriber[T any] = struct {
	key      string
	mut      sync.Mutex
	listener SignalListener[T]
}

// Signal implements a thread-safe publish-subscribe pattern for event handling.
// Use this to create signal channels that allow components to communicate asynchronously.
type Signal[T any] struct {
	subscribers    []*subscriber[T]
	subscribersMap *util.Set[string]
	subscriberMut  sync.Mutex
}

// New returns a new instance of a Signal.
func New[T any]() *Signal[T] {
	return &Signal[T]{
		subscribersMap: util.NewSet[string](0),
		subscriberMut:  sync.Mutex{},
	}
}

// AddListener registers your callback function to be notified when the signal is emitted.
// Provide a unique key for each listener to avoid duplicates and ensure proper removal later.
func (s *Signal[T]) AddListener(key string, listener SignalListener[T]) {
	s.subscriberMut.Lock()
	defer s.subscriberMut.Unlock()

	if exists := s.subscribersMap.Contains(key); exists {
		panic(fmt.Sprintf("key %s already exists", key))
	}

	s.subscribersMap.Add(key)
	s.subscribers = append(s.subscribers, &subscriber[T]{
		key:      key,
		mut:      sync.Mutex{},
		listener: listener,
	})
}

// RemoveListener unregisters a previously added signal listener.
// Use this when you no longer need to be notified about this signal.
func (s *Signal[T]) RemoveListener(key string) {
	s.subscriberMut.Lock()
	defer s.subscriberMut.Unlock()

	if removed := s.subscribersMap.Remove(key); !removed {
		return
	}

	for i, sub := range s.subscribers {
		if sub.key == key {
			s.subscribers = append(s.subscribers[:i], s.subscribers[i+1:]...)
			return
		}
	}
}

// emitOnGroup is an internal method used to execute all registered listeners.
// You typically won't call this directly, as it is not thread-safe - use Emit or EmitBG instead.
func (s *Signal[T]) emitOnGroup(ctx context.Context, payload T, wg *sync.WaitGroup) {
	for _, sub := range s.subscribers {
		wg.Add(1)
		go func(s *subscriber[T], p T, wg *sync.WaitGroup) {
			s.mut.Lock()
			defer s.mut.Unlock()
			defer wg.Done()

			s.listener(ctx, p)
		}(sub, payload, wg)
	}
}

// Emit sends a signal to all registered listeners and returns a WaitGroup.
// Use this when you need to wait for all listener functions to complete execution.
func (s *Signal[T]) Emit(ctx context.Context, payload T) *sync.WaitGroup {
	s.subscriberMut.Lock()
	defer s.subscriberMut.Unlock()
	var wg sync.WaitGroup

	s.emitOnGroup(ctx, payload, &wg)
	return &wg
}

// EmitBG is a convenience method for emitting signals with a background context.
// Use this when you don't need request-scoped data in your listeners.
func (s *Signal[T]) EmitBG(payload T) *sync.WaitGroup {
	return s.Emit(context.Background(), payload)
}

// EmitAll sends multiple related signals to all registered listeners in sequence.
// Use this when you need to emitOnGroup several related signals and want to wait for all of them to complete.
func (s *Signal[T]) EmitAll(ctx context.Context, payloads []T) *sync.WaitGroup {
	s.subscriberMut.Lock()
	defer s.subscriberMut.Unlock()
	var wg sync.WaitGroup

	for _, payload := range payloads {
		s.emitOnGroup(ctx, payload, &wg)
	}

	return &wg
}

// EmitAllBG sends multiple related signals using a background context.
// Use this when you need to emitOnGroup several related signals without request-scoped data.
func (s *Signal[T]) EmitAllBG(payloads []T) *sync.WaitGroup {
	return s.EmitAll(context.Background(), payloads)
}
